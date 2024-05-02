package server

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/chill/chill/internal/schema"
	"github.com/curtisnewbie/miso/miso"
)

func BootstrapServer() {
	// automatic MySQL schema migration using svc
	schema.EnableSchemaMigrateOnProd()

	miso.PreServerBootstrap(PreServerBootstrap)
	miso.PostServerBootstrapped(PostServerBootstrap)
	miso.BootstrapServer(os.Args)
}

func PreServerBootstrap(rail miso.Rail) error {
	EnableBasicAuth()
	RegisterEndpoints(rail)
	bc := LoadBuildsConf()
	if err := CheckBuildsConf(bc); err != nil {
		return fmt.Errorf("build conf illegal, %v", err)
	}
	InitBuildStatusMap(bc)
	return nil
}

func PostServerBootstrap(rail miso.Rail) error {
	// do stuff right after server being fully bootstrapped
	bc := LoadBuildsConf()
	if err := InitBuildInfo(rail, bc, miso.GetMySQL()); err != nil {
		return fmt.Errorf("failed to init build_info records, %v", err)
	}
	return nil
}
