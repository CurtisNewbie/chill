package build

import (
	"embed"
	"path/filepath"

	"github.com/curtisnewbie/miso/miso"
)

//go:embed scripts/*
var buildScriptFs embed.FS

const (
	PropScriptsBaseFolder = "scripts.base-folder"
)

func init() {
	miso.SetDefProp(PropScriptsBaseFolder, "./")
}

// Lookup build script under the specified folder.
func LookupBuildScript(path string) ([]byte, error) {
	return buildScriptFs.ReadFile(filepath.Join(miso.GetPropStr(PropScriptsBaseFolder), path))
}
