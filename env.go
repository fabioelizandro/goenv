package goenv

import (
	"fmt"
	"os"
)

// Env combines both, .env and OS variables in a single API
// OS variables always takes precedence of any other variable declaration
type Env struct {
	dotfile Dotfile
}

func NewEnv(dotfile Dotfile) Env {
	return Env{dotfile: dotfile}
}

func (e Env) ReadOrDefault(envVar string, defaultValue string) string {
	value, ok := e.read(envVar)
	if !ok {
		return defaultValue
	}

	return value
}

func (e Env) MustRead(envVar string) string {
	value, ok := e.read(envVar)
	if !ok {
		panic(fmt.Errorf("variable %s is not declared", envVar))
	}

	return value
}

func (e Env) read(envVar string) (string, bool) {
	osValue, ok := os.LookupEnv(envVar)
	if ok {
		return osValue, true
	}

	dotenvValue, ok := e.dotfile.vars[envVar]
	if ok {
		return dotenvValue, true
	}

	return "", false
}
