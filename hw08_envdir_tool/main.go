package main

import (
	"fmt"
	"os"
)

func main() {
	//dir, err := ReadDir(os.Args[1])
	mapEnv, err := ReadDir("testdata/env")
	if err != nil {
		fmt.Println(fmt.Errorf("could not ge env from dir: %w", err))
		return
	}

	//code := RunCmd(os.Args[2:], mapEnv)
	code := RunCmd([]string{"bash", "-c", "$(pwd)/testdata/echo.sh", "arg1=1", "arg2=2"}, mapEnv)
	os.Exit(code)
}
