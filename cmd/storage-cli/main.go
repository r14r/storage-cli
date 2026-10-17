package main

import (
 "fmt"
 "os"
 "github.com/r14r/storage-cli/internal/storage"
)
func main() {
 if err := storage.Run(os.Args[1:], os.Stdout); err != nil {
  fmt.Fprintln(os.Stderr,"error:",err)
  os.Exit(1)
 }
}
