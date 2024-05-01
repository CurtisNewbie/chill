package build

import "testing"

func TestLookupBuildScript(t *testing.T) {
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
