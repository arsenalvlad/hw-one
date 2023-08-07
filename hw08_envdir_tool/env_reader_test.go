package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("done empty and unset case env", func(t *testing.T) {
		mapEnv, err := ReadDir("testdata/env")
		if err != nil {
			fmt.Println(fmt.Errorf("could not ge env from dir: %w", err))
			return
		}

		require.Equal(t, mapEnv["EMPTY"].Value, "")
		require.Equal(t, mapEnv["UNSET"].NeedRemove, true)
	})
}
