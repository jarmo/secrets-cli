package input

import (
  "os"
  "fmt"
  "strings"
  "bufio"
  "syscall"
  "golang.org/x/crypto/ssh/terminal"
)

func Ask(message string) string {
  fmt.Print(message)
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Buffer(make([]byte, 4096), 4096)
  if !scanner.Scan() {
    if err := scanner.Err(); err != nil {
      panic(fmt.Errorf("Unable to read input: %w", err))
    }
    return ""
  }
  return scanner.Text()
}

func AskMultiline(message string) string {
  fmt.Print(message)
  var value []string
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Buffer(make([]byte, 4096), 8<<20)
  for scanner.Scan() {
      value = append(value, scanner.Text())
  }

  if err := scanner.Err(); err != nil {
    panic(fmt.Errorf("Unable to read secret value: %w", err))
  }
  return replaceUnprintableCharacters(strings.Join(value, "\n"))
}

func AskPassword(message string) []byte {
  fmt.Print(message)
  password, err := terminal.ReadPassword(int(syscall.Stdin))
  if err != nil {
    panic(err)
  }
  fmt.Println()

  return password
}

func AskVaultPassword() []byte {
  return AskPassword("Enter vault password: ")
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
