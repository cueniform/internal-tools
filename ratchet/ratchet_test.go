package ratchet_test

import (
	"strings"
	"testing"

	"github.com/cueniform/internal-tools/ratchet"
	"github.com/google/go-cmp/cmp"
)

func TestProviderData_ReturnsExpectedStringGivenProviderSchemaPathAndMatchingProviderAddress(t *testing.T) {
	t.Parallel()
	want := `{
"entity_schemas":{
"entity_id":{
"attributes":{
"attribute_id":{
"type":["list","string"],
"required":true
}
}
}
}
}`
	rt, err := ratchet.New("testdata/example.json", "provider.registry/provider_name")
	if err != nil {
		t.Fatal(err)
	}
	got := strings.ReplaceAll(rt.ProviderSchema, " ", "")
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
