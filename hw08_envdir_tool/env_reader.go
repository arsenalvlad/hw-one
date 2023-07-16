package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)

	readDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, item := range readDir {
		open, err := os.Open(dir + "/" + item.Name())
		if err != nil {
			return nil, fmt.Errorf("could not oepn file: %w", err)
		}

		defer func(open *os.File) {
			err := open.Close()
			if err != nil {
				fmt.Println(fmt.Errorf("could not close file: %w", err))
			}
		}(open)

		info, err := open.Stat()
		if err != nil {
			return nil, fmt.Errorf("could not get stat: %w", err)
		}

		if info.Size() == 0 {
			result[item.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		res, err := prepareEnv(open)
		if err != nil {
			return nil, fmt.Errorf("could not prepare env: %w", err)
		}

		result[item.Name()] = *res
	}
	// Place your code here
	return result, nil
}

func prepareEnv(reader io.Reader) (*EnvValue, error) {
	rd := bufio.NewReader(reader)

	line, err := rd.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("could not read string: %w", err)
	}

	un, err := strconv.Unquote(line)
	if err != nil {
		un = line
	}

	res := strings.ReplaceAll(un, "0x00", string('\n'))

	value := strings.TrimRight(res, " ")

	matchString, err := regexp.MatchString("^[a-zA-Z]+$", value)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("could not match string: %w", err)
	}

	if !matchString {
		return &EnvValue{
			Value:      "",
			NeedRemove: true,
		}, nil
	}

	return &EnvValue{
		Value:      value,
		NeedRemove: false,
	}, nil
}
