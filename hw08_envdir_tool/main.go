package main

import (
	"fmt"
	"os"
)

func main() {
	mapEnv, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Errorf("could not ge env from dir: %w", err))
		return
	}

	code := RunCmd(os.Args[2:], mapEnv)
	os.Exit(code)
}
