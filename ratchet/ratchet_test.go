package ratchet_test

import (
	"bytes"
	"testing"

	"github.com/cueniform/internal-tools/ratchet"
	"github.com/google/go-cmp/cmp"
)

func TestNewErrorsGivenNonExistentFileAsInput(t *testing.T) {
	t.Parallel()
	_, err := ratchet.New(t.TempDir()+"/bogus.json", "bogus")
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestNewSetExpectedDataGivenExistentAndValidFileAsInput(t *testing.T) {
	t.Parallel()
	want := []byte(`{
  "provider_schemas": {
    "provider.registry/provider_name": {
      "resource_schemas": {
        "resource_id": {
          "attributes": {
            "attribute_id": {
              "type": ["list", "string"],
              "required": true
            }
          }
        }
      }
    }
  }
}`)
	ratchet, err := ratchet.New("testdata/input.json", "bogus")
	if err != nil {
		t.Fatal(err)
	}
	got := ratchet.ProviderSchemaData
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestNewRemovesLastNewLineFromInputFile(t *testing.T) {
	t.Parallel()
	want := []byte("{}")
	ratchet, err := ratchet.New("testdata/newline.json", "bogus")
	if err != nil {
		t.Fatal(err)
	}
	got := ratchet.ProviderSchemaData
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestNewSetExpectedProviderAddressGivenStringContainingDot(t *testing.T) {
	t.Parallel()
	want := "provider.registry/provider_name"
	ratchet, err := ratchet.New("testdata/emptyfile.txt", "provider.registry/provider_name")
	if err != nil {
		t.Fatal(err)
	}
	got := ratchet.ProviderAddress
	if want != got {
		t.Fatalf("want provider address %q, got %q", want, got)
	}
}

func TestProviderData_ReturnsExpectedDataGivenProviderSchemaPathAndMatchingProviderAddress(t *testing.T) {
	t.Parallel()
	want := []byte(`{
"resource_schemas":{
"resource_id":{
"attributes":{
"attribute_id":{
"type":["list","string"],
"required":true
}
}
}
}
}`)
	ratchet, err := ratchet.New("testdata/input.json", "provider.registry/provider_name")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ratchet.ProviderData()
	if err != nil {
		t.Fatal(err)
	}
	got := bytes.ReplaceAll(data, []byte(" "), []byte(""))
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
