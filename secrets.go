package main

import (
  "os"
  "github.com/jarmo/secrets-cli/v7/cli"
)

const VERSION = "7.1.0"

var version = "v" + VERSION

func main() {
  cli.Command("secrets " + version, os.Args[1:]).Execute()
}

