package server

import (
	"os"
	"testing"

	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

func TestBashRunNormal(t *testing.T) {
	_, err := BashRun(miso.EmptyRail(), []byte(`touch "abc.txt" && \
echo "yes" > "abc.txt"`))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("abc.txt")
	s, err := util.ReadFileAll("abc.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(s))
}

func TestBashRunFailed(t *testing.T) {
	out, err := BashRun(miso.EmptyRail(), []byte(`ech`))
	if err == nil {
		t.Fatal("command should failed")
	}
	t.Log(out)
}
