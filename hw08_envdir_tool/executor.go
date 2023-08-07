package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitError *exec.ExitError

	for name, item := range env {
		_, ok := os.LookupEnv(name)
		if ok {
			err := os.Unsetenv(name)
			if err != nil {
				fmt.Println(fmt.Errorf("could not unsetenv: %w", err))
				return exitError.ExitCode()
			}
		}

		if item.NeedRemove {
			continue
		}

		err := os.Setenv(name, item.Value)
		if err != nil {
			fmt.Println(fmt.Errorf("could not setenv: %w", err))
			return exitError.ExitCode()
		}
	}

	res := exec.Command(cmd[0], cmd[1:]...) //nolint: gosec
	res.Stdout = os.Stdout
	res.Stderr = os.Stderr
	res.Env = os.Environ()

	if err := res.Run(); err != nil {
		fmt.Println(fmt.Errorf("could not cmd run: %w", err))
		return exitError.ExitCode()
	}

	return 0
}
