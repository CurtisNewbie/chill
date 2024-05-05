package server

import "embed"

//go:embed static
var staticFs embed.FS

const (
	taticFsPre = "/static"
)
