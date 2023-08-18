package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("done check env for empty, unset", func(t *testing.T) {
		mapEnv := make(map[string]EnvValue)

		mapEnv["EMPTY"] = EnvValue{
			Value:      "",
			NeedRemove: false,
		}

		mapEnv["UNSET"] = EnvValue{
			Value:      "",
			NeedRemove: true,
		}

		code := RunCmd([]string{"echo", "test work"}, mapEnv)

		_, okEmpty := os.LookupEnv("EMPTY")
		_, okUnset := os.LookupEnv("UNSET")

		require.Equal(t, code, 0)
		require.Equal(t, okEmpty, true)
		require.Equal(t, okUnset, false)
	})
}
