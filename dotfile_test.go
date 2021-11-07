package goenv_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/fabioelizandro/goenv"
	"github.com/stretchr/testify/assert"
)

func TestMustParseDotenvFile(t *testing.T) {
	t.Run("empty variables when file is empty", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{}, dotenv.Vars())
	})

	t.Run("returns single variable", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=BAR")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR"},
			dotenv.Vars(),
		)
	})

	t.Run("returns many variables", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=BAR\nBAR=FOO")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("ignores blank lines", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("\nFOO=BAR\n\nBAR=FOO\n")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("ignores lines with spaces", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("  \nFOO=BAR\n   \nBAR=FOO\n  ")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("ignores surround whitespaces", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("   FOO=BAR  ")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR"},
			dotenv.Vars(),
		)
	})

	t.Run("ignores value whitespaces", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=   BAR  ")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR"},
			dotenv.Vars(),
		)
	})

	t.Run("supports empty variables", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": ""},
			dotenv.Vars(),
		)
	})

	t.Run("preserves spaces between double quotes", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile(`FOO=" BAR "`)
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": " BAR "},
			dotenv.Vars(),
		)
	})

	t.Run("preserves spaces between single quotes", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile(`FOO=' BAR '`)
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": " BAR "},
			dotenv.Vars(),
		)
	})

	t.Run("ignores commented lines", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=BAR\n#BAR=FOO")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR"},
			dotenv.Vars(),
		)
	})

	t.Run("ignores commented lines", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfile("FOO=BAR\n#BAR=FOO")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR"},
			dotenv.Vars(),
		)
	})

	t.Run("returns error when file format is invalid", func(t *testing.T) {
		_, err := goenv.ParseDotfile("FOO=BAR\nINVALID_LINE")
		assert.EqualError(t, err, "dotenv file contains invalid value at line 2")

		_, err = goenv.ParseDotfile("FOO=BAR\nBAR=FOO\nINVALID_LINE")
		assert.EqualError(t, err, "dotenv file contains invalid value at line 3")
	})
}

func TestParseDotfileFromIOReader(t *testing.T) {
	t.Run("parses dotenv file from IO reader", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfileFromIOReader(bytes.NewBufferString("FOO=BAR\nBAR=FOO"))
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("returns error when fails to read", func(t *testing.T) {
		_, err := goenv.ParseDotfileFromIOReader(newBrokenIOReader(errors.New("broken reader")))
		assert.EqualError(t, err, "broken reader")
	})
}

func TestParseDotfileFromFilepath(t *testing.T) {
	t.Run("parses dotenv file from filepath", func(t *testing.T) {
		dotenv, err := goenv.ParseDotfileFromFilepath(".env-sample")
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("returns error when fails to read file", func(t *testing.T) {
		_, err := goenv.ParseDotfileFromFilepath(".env-not-found")
		assert.EqualError(t, err, "open .env-not-found: no such file or directory")
	})
}

func TestMustParseDotfileFromFilepath(t *testing.T) {
	t.Run("parses dotenv file from filepath", func(t *testing.T) {
		dotenv := goenv.MustParseDotfileFromFilepath(".env-sample")
		assert.Equal(
			t,
			map[string]string{"FOO": "BAR", "BAR": "FOO"},
			dotenv.Vars(),
		)
	})

	t.Run("panics when fails to read file", func(t *testing.T) {
		assert.PanicsWithError(t, "open .env-not-found: no such file or directory", func() {
			_ = goenv.MustParseDotfileFromFilepath(".env-not-found")
		})
	})
}

// --- Test Doubles

type brokenIOReader struct {
	err error
}

func newBrokenIOReader(err error) *brokenIOReader {
	return &brokenIOReader{err: err}
}

func (b *brokenIOReader) Read(p []byte) (n int, err error) {
	return 0, b.err
}
