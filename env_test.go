package goenv_test

import (
	"os"
	"testing"

	"github.com/fabioelizandro/goenv"
	"github.com/stretchr/testify/assert"
)

func TestReadOrDefault(t *testing.T) {
	t.Run("returns default value when dotenv and OS variable are not present", func(t *testing.T) {
		assert.NoError(t, os.Unsetenv("MY_ENV_VARIABLE"))
		dotfile, err := goenv.ParseDotfile("")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.Equal(
			t,
			"my-default-value",
			env.ReadOrDefault("MY_ENV_VARIABLE", "my-default-value"),
		)
	})

	t.Run("returns dotenv value when OS value is not present", func(t *testing.T) {
		assert.NoError(t, os.Unsetenv("MY_ENV_VARIABLE"))
		dotfile, err := goenv.ParseDotfile("MY_ENV_VARIABLE=dotenv-value")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.Equal(
			t,
			"dotenv-value",
			env.ReadOrDefault("MY_ENV_VARIABLE", "my-default-value"),
		)
	})

	t.Run("returns OS variable when present", func(t *testing.T) {
		assert.NoError(t, os.Setenv("MY_ENV_VARIABLE", "os-value"))
		dotfile, err := goenv.ParseDotfile("MY_ENV_VARIABLE=dotenv-value")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.Equal(
			t,
			"os-value",
			env.ReadOrDefault("MY_ENV_VARIABLE", "my-default-value"),
		)
	})
}

func TestMustRead(t *testing.T) {
	t.Run("panics when variable is not present in dotenv file and OS", func(t *testing.T) {
		assert.NoError(t, os.Unsetenv("MY_ENV_VARIABLE"))
		dotfile, err := goenv.ParseDotfile("")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.PanicsWithError(t, "variable MY_ENV_VARIABLE is not declared", func() {
			env.MustRead("MY_ENV_VARIABLE")
		})
	})

	t.Run("returns dotenv value when os value is not present", func(t *testing.T) {
		assert.NoError(t, os.Unsetenv("MY_ENV_VARIABLE"))
		dotfile, err := goenv.ParseDotfile("MY_ENV_VARIABLE=dotenv-value")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.Equal(
			t,
			"dotenv-value",
			env.MustRead("MY_ENV_VARIABLE"),
		)
	})

	t.Run("returns OS variable when present", func(t *testing.T) {
		assert.NoError(t, os.Setenv("MY_ENV_VARIABLE", "os-value"))
		dotfile, err := goenv.ParseDotfile("MY_ENV_VARIABLE=dotenv-value")
		assert.NoError(t, err)
		env := goenv.NewEnv(dotfile)

		assert.Equal(
			t,
			"os-value",
			env.MustRead("MY_ENV_VARIABLE"),
		)
	})
}
