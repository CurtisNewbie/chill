package build

import (
	"embed"
	"path/filepath"
)

//go:embed scripts/*
var buildScriptFs embed.FS

const (
	baseFolder = "scripts"
)

// Lookup build script under scripts/ folder.
//
// E.g., for `scripts/echo.sh`, path should be 'echo.sh'
func LookupBuildScript(path string) ([]byte, error) {
	return buildScriptFs.ReadFile(filepath.Join(baseFolder, path))
}
