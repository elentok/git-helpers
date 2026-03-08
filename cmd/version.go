package cmd

import (
	"fmt"
	"io"
	"runtime/debug"
)

// version is set at build time via -ldflags "-X gx/cmd.version=vX.Y.Z"
var version = ""

func getVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" && info.Main.Version != "" {
		return info.Main.Version
	}
	return "dev"
}

func runVersion(w io.Writer) error {
	fmt.Fprintf(w, "gx %s\n", getVersion())
	return nil
}
