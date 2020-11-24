package main

import (
  "github.com/sidilabs/kishell/pkg/options"
)

func main() {
  option := options.Parse()
  option.Run()
}
