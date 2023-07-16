package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error open from file: %w", err)
	}

	stat, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("could not get stat from file: %w", err)
	}

	countLimit := stat.Size()

	if countLimit < offset {
		return fmt.Errorf("could not offet size: %w", ErrOffsetExceedsFileSize)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error open to file: %w", err)
	}

	defer func(fromFile *os.File) {
		err := fromFile.Close()
		if err != nil {
			fmt.Println(fmt.Errorf("error close from file: %w", err))
		}
	}(fromFile)

	defer func(toFile *os.File) {
		err := toFile.Close()
		if err != nil {
			fmt.Println(fmt.Errorf("error close to file: %w", err))
		}
	}(toFile)

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("could not seek offset file: %w", err)
	}

	if limit != 0 {
		countLimit = limit
		if stat.Size()-offset-limit < 0 {
			countLimit = stat.Size() - offset
		}
	}

	bar := pb.Full.Start64(countLimit)

	bar.SetRefreshRate(100 * time.Nanosecond)

	// create proxy reader
	barReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barReader, countLimit)
	if err != nil {
		return fmt.Errorf("error copyN: %w", err)
	}

	bar.Finish()

	// Place your code here.
	return nil
}
