package server

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/curtisnewbie/miso/miso"
)

func BashRun(rail miso.Rail, script []byte) (string, error) {
	c := exec.Command("bash")
	c.Stdin = bytes.NewReader(script)

	var outstr string
	var cmdout []byte
	var err error
	if cmdout, err = c.CombinedOutput(); err != nil {
		if cmdout != nil {
			outstr = string(cmdout)
		}
		return outstr, fmt.Errorf("failed to execute bash script, %s, %w", outstr, err)
	}
	return outstr, nil
}
