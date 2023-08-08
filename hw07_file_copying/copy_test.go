package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fromPath   = "testdata/input.txt"
	toPath     = "testdata/test1.txt"
	checkPath  = "testdata/out_offset2_limit4.txt"
	offsetTest = 2
	limitTest  = 4
)

func TestCopy(t *testing.T) {
	err := Copy(fromPath, toPath, offsetTest, limitTest)
	if err != nil {
		return
	}

	checkFile, err := os.Open(checkPath)
	if err != nil {
		return
	}

	defer checkFile.Close()

	statCheck, err := checkFile.Stat()
	if err != nil {
		return
	}

	toFile, err := os.Open(toPath)
	if err != nil {
		return
	}

	defer toFile.Close()

	statTo, err := toFile.Stat()
	if err != nil {
		return
	}

	require.Equal(t, statCheck.Size(), statTo.Size())

	bufCheck := make([]byte, 1024)

	_, err = checkFile.Read(bufCheck)
	if err != nil {
		return
	}

	bufTo := make([]byte, 1024)

	_, err = toFile.Read(bufTo)
	if err != nil {
		return
	}

	require.Equal(t, bufCheck, bufTo)

	cmd := exec.Command("rm", "-f", toPath)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
