package goenv

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var variableExpression = regexp.MustCompile(`(.*)= *[",']?([^",']*)[",']? *`)

// Dotfile is the interface of the .env file in the user's file system
type Dotfile struct {
	vars map[string]string
}

func MustParseDotfileFromFilepath(filepath string) Dotfile {
	dotfile, err := ParseDotfileFromFilepath(filepath)
	if err != nil {
		panic(err)
	}

	return dotfile
}

func ParseDotfileFromFilepath(filepath string) (Dotfile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Dotfile{}, err
	}

	return ParseDotfileFromIOReader(file)
}

func ParseDotfileFromIOReader(reader io.Reader) (Dotfile, error) {
	fileContent, err := io.ReadAll(reader)
	if err != nil {
		return Dotfile{}, err
	}

	return ParseDotfile(string(fileContent))
}

func ParseDotfile(fileContent string) (Dotfile, error) {
	vars := map[string]string{}

	for lineIndex, variable := range strings.Split(fileContent, "\n") {
		variable = strings.TrimSpace(variable)

		if variable == "" || strings.HasPrefix(variable, "#") {
			continue
		}

		result := variableExpression.FindStringSubmatch(variable)
		if len(result) != 3 {
			return Dotfile{}, fmt.Errorf("dotenv file contains invalid value at line %d", lineIndex+1)
		}

		vars[result[1]] = result[2]
	}

	return Dotfile{vars: vars}, nil
}

func (f Dotfile) Vars() map[string]string {
	return f.vars
}
