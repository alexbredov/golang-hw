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
	target, err := os.Create(toPath)
	if err != nil {
		return ErrFileCreate
	}
	defer target.Close()
	pBar := pb.Full.Start64(limit)
	var buf []byte
	if limit == 0 {
		read, err := io.ReadAll(file)
		buf = read[offset:]
		if err != nil {
			return err
		}
	} else {
		var bufLength int64
		if offset+limit > fileInfo.Size() {
			bufLength = fileInfo.Size() - offset
		} else {
			bufLength = limit
		}
		buf = make([]byte, bufLength)
		file.ReadAt(buf, offset)
		pBar.Finish()
	}
	_, err = target.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
