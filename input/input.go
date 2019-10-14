package input

import (
  "os"
  "fmt"
  "strings"
  "bufio"
  "github.com/jarmo/secrets/input"
)

func Ask(message string) string {
  fmt.Print(message)
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  return scanner.Text()
}

func AskMultiline(message string) string {
  fmt.Print(message)
  var value []string
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
      value = append(value, scanner.Text())
  }

  return replaceUnprintableCharacters(strings.Join(value, "\n"))
}

func AskVaultPassword() []byte {
  return input.AskPassword("Enter vault password: ")
}

func replaceUnprintableCharacters(s string) string {
  ctrlD, ctrlX, ctrlZ := "\x04", "\x18", "\x1A"

  return strings.Replace(
    strings.Replace(
      strings.Replace(s,
      ctrlD, "", -1),
      ctrlX, "", -1),
      ctrlZ, "", -1)
}
