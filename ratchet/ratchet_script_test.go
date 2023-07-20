package ratchet_test

import (
	"os"
	"testing"

	"github.com/cueniform/internal-tools/ratchet"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestRatchetEmitResourcesEntities(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script/unit/resources",
		// UpdateScripts: highly discouraged here,
	})
}

func TestRatchetEmitDataSourcesEntities(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script/unit/data_sources",
		// UpdateScripts: highly discouraged here,
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"ratchet": ratchet.Main,
	}))
}
