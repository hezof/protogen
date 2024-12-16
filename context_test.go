package protogen

import "testing"

func TestPrintUsage(t *testing.T) {
	ctx := NewContext([]string{"protogen", "-help"})
	defer ctx.Close()

	ctx.PrintVersion()
}
