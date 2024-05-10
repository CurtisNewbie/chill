package server

import "github.com/curtisnewbie/miso/miso"

const (
	PropUsername = "auth.basic.username"
	PropPassword = "auth.basic.password"
)

func EnableBasicAuth() {
	if !miso.IsProdMode() {
		return
	}
	miso.EnableBasicAuth(func(username, password, url, method string) bool {
		return username == miso.GetPropStr(PropUsername) && password == miso.GetPropStr(PropPassword)
	})
}
