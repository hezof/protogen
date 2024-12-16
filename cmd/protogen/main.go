package main

import (
	"github.com/hezof/protogen"
	"os"
)

func main() {
	protogen.Main(os.Args[1:])
}
