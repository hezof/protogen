package protogen

import "testing"

func TestMainCase(t *testing.T) {
	args := []string{
		"--grpc",
		"--root_path",
		"D:\\Workspace\\hezof\\github.com\\hezof\\protogen\\test",
		"protoapi.proto",
	}
	Main(args)
}
