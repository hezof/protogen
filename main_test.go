package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMainCase(t *testing.T) {
	pid := os.Getpid()
	prc, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	fmt.Println(prc.Pid)
}
