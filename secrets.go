package main

import (
  "os"
  "github.com/jarmo/secrets-cli/v6/cli"
)

const VERSION = "6.1.1"

func main() {
  cli.Command(VERSION, os.Args[1:]).Execute()
}

