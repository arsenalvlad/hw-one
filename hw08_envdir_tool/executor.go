package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	res := exec.Command(cmd[0], cmd[1:]...)
	res.Stdout = os.Stdout
	res.Stderr = os.Stderr

	for name, item := range env {
		_, ok := os.LookupEnv(name)
		if ok {
			err := os.Unsetenv(name)
			if err != nil {
				fmt.Println(fmt.Errorf("could not unsetenv: %w", err))
				return res.ProcessState.ExitCode()
			}
		}

		if item.NeedRemove {
			continue
		}

		err := os.Setenv(name, item.Value)
		if err != nil {
			fmt.Println(fmt.Errorf("could not setenv: %w", err))
			return res.ProcessState.ExitCode()
		}

	}

	res.Env = append(os.Environ())

	if err := res.Run(); err != nil {
		fmt.Println(fmt.Errorf("could not cmd run: %w", err))
		return res.ProcessState.ExitCode()
	}
	// Place your code here.
	return res.ProcessState.ExitCode()
}
