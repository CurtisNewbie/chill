package server

import (
	"os"

	"github.com/curtisnewbie/chill/chillback/internal/schema"
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
	// declare http endpoints, jobs/tasks, and other components here
	return nil
}

func PostServerBootstrap(rail miso.Rail) error {
	// do stuff right after server being fully bootstrapped
	return nil
}
