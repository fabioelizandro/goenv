# Goenv

Golang-library that facilitates the use of `.env` files. 

## Installation

```
go get github.com/fabioelizandro/goenv
```

## Usage

Place a .env file in the root of your directory.
```
FOO=BAR
BAR=FOO
```

Now use the goenv.Env object to read it.

```go
package main

import "github.com/fabioelizandro/goenv"

func main() {
	env := goenv.NewEnv(goenv.MustParseDotfileFromFilepath(".env"))
	
	println(env.ReadOrDefault("FOO", "my default value")) // returns default value if variable is not present in .env and OS
	println(env.MustRead("BAR")) // panics if variable is not present in .env and OS 
}
```

## Contributing

First run:
```
make setup
```

Run `make help` for help.
