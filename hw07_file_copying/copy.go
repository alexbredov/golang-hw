package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
	ErrFileOpen              = errors.New("file open failed")
	ErrFileCreate            = errors.New("file creation failed")
	ErrNoOriginalFileSize    = errors.New("original file size undefined")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if filepath.Ext(fromPath) == "" || filepath.Ext(fromPath) != filepath.Ext(toPath) {
		return ErrUnsupportedFile
	}
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	file, err := os.Open(fromPath)
	if err != nil {
		return ErrFileOpen
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return ErrFileOpen
	}
	defer file.Close()
	if fileInfo.Size() <= limit {
		limit = fileInfo.Size()
	}
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if fileInfo.Size() == 0 {
		return ErrNoOriginalFileSize
	}
	var totalToCopy int64
	if limit == 0 || offset+limit > fileInfo.Size() {
		totalToCopy = fileInfo.Size() - offset
	} else {
		totalToCopy = limit
	}
	targetFile, err := os.Create(toPath)
	if err != nil {
		return ErrFileCreate
	}
	defer targetFile.Close()
	pBar := pb.Full.Start64(totalToCopy)
	defer pBar.Finish()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	bufSize := 1024
	buf := make([]byte, bufSize)
	var copied int64
	for copied < totalToCopy {
		bytesLeft := totalToCopy - copied
		if int64(bufSize) > bytesLeft {
			buf = make([]byte, bytesLeft)
		}
		n, err := file.Read(buf)
		if n > 0 {
			written, writeErr := targetFile.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			if written < n {
				return io.ErrShortWrite
			}
			copied += int64(n)
			pBar.Add(n)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
