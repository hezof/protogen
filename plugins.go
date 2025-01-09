package protogen

var Plugins = []Plugin{
	{
		Mode:    Protoc,
		Name:    "protoc",
		Module:  "com/google/protobuf/protoc",
		Version: "v3.21.9",
	},
	{
		Mode:    GoSrc,
		Name:    "include",
		Module:  "github.com/hezof/protogen/cmd/include",
		Version: "v0.5.0",
	},
	{
		Mode:    GoBin,
		Name:    "protoc-gen-go-protoapi",
		Module:  "github.com/hezof/protogen/cmd/protoc-gen-go-protoapi",
		Version: "v0.5.0",
	},
	{
		Mode:    GoBin,
		Name:    "protoc-gen-go-openapi",
		Module:  "github.com/hezof/protogen/cmd/protoc-gen-go-openapi",
		Version: "v0.5.0",
	},
	{
		Mode:    GoBin,
		Name:    "protoc-gen-go",
		Module:  "google.golang.org/protobuf/cmd/protoc-gen-go",
		Version: "v1.36.2",
	},
	{
		Mode:    GoBin,
		Name:    "protoc-gen-go-grpc",
		Module:  "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		Version: "v1.5.1",
	},
}

type Mode uint8

const (
	GoBin  Mode = 0
	GoSrc  Mode = 1
	Protoc Mode = 2
)

type Plugin struct {
	Mode    Mode
	Name    string
	Module  string
	Version string
}
