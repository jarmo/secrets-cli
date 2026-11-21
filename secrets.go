package main

import (
  "os"
  "github.com/jarmo/secrets-cli/v6/cli"
)

const VERSION = "6.1.2"

var version = "v" + VERSION

func main() {
  cli.Command("secrets " + version, os.Args[1:]).Execute()
}

