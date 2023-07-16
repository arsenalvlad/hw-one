package main

import (
	"fmt"
	"os"
)

func main() {
	//dir, err := ReadDir(os.Args[1])
	dir, err := ReadDir("testdata/env")
	if err != nil {
		fmt.Println(fmt.Errorf("could not ge env from dir: %w", err))
		return
	}

	code := RunCmd(os.Args[2:], dir)
	os.Exit(code)
}
