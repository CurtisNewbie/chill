package server

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/curtisnewbie/miso/miso"
)

func BashRun(rail miso.Rail, script []byte) (string, error) {
	c := exec.Command("bash")
	c.Stdin = bytes.NewReader(script)

	var cmdout []byte
	var err error
	if cmdout, err = c.CombinedOutput(); err != nil {
		var outstr string
		if cmdout != nil {
			outstr = strings.TrimSpace(miso.UnsafeByt2Str(cmdout))
		}
		return outstr, fmt.Errorf("%s, %w", outstr, err)
	}
	return strings.TrimSpace(miso.UnsafeByt2Str(cmdout)), nil
}
