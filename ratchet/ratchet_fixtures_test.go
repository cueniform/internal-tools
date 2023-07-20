//go:build fixtures

package ratchet_test

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestRatchet(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script/fixtures",
		// UpdateScripts: true,
	})
}
