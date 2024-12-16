package protogen

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
)

func TestPrintUsage(t *testing.T) {
	cwd, _ := filepath.Abs("./../")
	fmt.Println(cwd)
	filepath.Walk(cwd, func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
}
