package main

import (
	"github.com/vncommunityleague/kazusa/cmd/daemon"
	"github.com/vncommunityleague/kazusa/registry"
)

func main() {
	r := registry.NewRegistryDefault()

	daemon.ServePublic(r)
}
