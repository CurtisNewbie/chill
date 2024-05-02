package server

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
)

func TestLookupBuildScript(t *testing.T) {
	miso.SetProp(PropScriptsBaseFolder, "../../testdata")
	f, err := LookupBuildScript("echo.sh")
	if err != nil {
		t.Fatal(err)
	}
	fstr := string(f)
	if fstr != "echo \"Lets go!\"" {
		t.Fatal(fstr)
	}
	t.Log(fstr)
}

func TestLoadBuildConf(t *testing.T) {
	err := miso.LoadConfigFromFile("../../conf.yml", miso.EmptyRail())
	if err != nil {
		t.Fatal(err)
	}
	builds := LoadBuilds()
	t.Logf("%#v", builds)
}
