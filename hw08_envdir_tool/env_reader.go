package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
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
		name := strings.ReplaceAll(item.Name(), "=", "")
		open, err := os.Open(dir + "/" + name)
		if err != nil {
			return nil, fmt.Errorf("could not oepn file: %w", err)
		}

		info, err := open.Stat()
		if err != nil {
			return nil, fmt.Errorf("could not get stat: %w", err)
		}

		if info.Size() == 0 {
			result[name] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		res, err := prepareEnv(open)
		if err != nil {
			return nil, fmt.Errorf("could not prepare env: %w", err)
		}

		result[name] = *res

		err = open.Close()
		if err != nil {
			fmt.Println(fmt.Errorf("could not close file: %w", err))
		}
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

	value1 := strings.TrimRight(line, " ")
	value2 := strings.ReplaceAll(value1, string([]byte{0x00}), string('\n'))
	value3 := strings.Trim(value2, string('"'))

	return &EnvValue{
		Value:      value3,
		NeedRemove: false,
	}, nil
}
