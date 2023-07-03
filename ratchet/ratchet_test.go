package ratchet_test

import (
	"testing"

	"github.com/cueniform/internal-tools/ratchet"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
)

func TestValidTFSChemaVersion_ReturnsTrueGivenKnownVersion(t *testing.T) {
	t.Parallel()
	want := true
	input := gjson.Result{
		Type: gjson.String,
		Raw:  "1.0",
		Str:  "1.0",
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}

}

func TestValidTFSChemaVersion_ReturnsFalseGivenNumber(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.Number,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenFalse(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.False,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenTrue(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.True,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenNull(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.Null,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenJSON(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.JSON,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenUnkownType(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: 666,
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenString(t *testing.T) {
	t.Parallel()
	want := "string"
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.String,
		Raw:  `"string"`,
		Str:  "string",
	})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenNumber(t *testing.T) {
	t.Parallel()
	want := "number"
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.String,
		Raw:  `"number"`,
		Str:  "number",
	})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenListString(t *testing.T) {
	t.Parallel()
	want := "[...string]"
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.JSON,
		Raw:  `["list", "string"]`,
		Str:  `["list", "string"]`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenSetString(t *testing.T) {
	t.Parallel()
	want := "[...string]"
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.JSON,
		Raw:  `["set", "string"]`,
		Str:  `["set", "string"]`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenListBool(t *testing.T) {
	t.Parallel()
	want := "[...bool]"
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.JSON,
		Raw:  `["list", "bool"]`,
		Str:  `["list", "bool"]`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConvertTerraformType_ReturnsExpectedStringGivenMapString(t *testing.T) {
	t.Parallel()
	want := `[...close({
        role_arn: string
        role_type: string
    })]`
	got, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.JSON,
		Raw: `["set", [
                    "object",
                    {
                      "role_arn": "string",
                      "role_type": "string"
                    }
                  ]]`,
		Str: `["set", [
                    "object",
                    {
                      "role_arn": "string",
                      "role_type": "string"
                    }
                  ]]`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestConvertTerraformType_ErrorsGivenListWithThreeItems(t *testing.T) {
	t.Parallel()
	_, err := ratchet.ConvertTerraformType(gjson.Result{
		Type: gjson.JSON,
		Raw:  `["a", "b", "c"]`,
		Str:  `["a", "b", "c"]`,
	})
	if err == nil {
		t.Fatal("want error but got nil")
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenDataSourceWithBasicAttribute(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    bogus: string
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_SkipsItemsWithComputedTrue(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
				"computed": true
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenDataSourceWithRequiredAttribute(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    bogus!: string
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string",
				"required": true
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenDataSourceWithOptionalAttribute(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    bogus?: string
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string",
				"optional": true
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenResource(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {
    bogus: string
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_ReturnsExpectedEntityGivenDataSourceWithTwoAttributes(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    foo: string
    bar: number
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "version": 0,
          "block": {
            "attributes": {
              "foo": {
                "type": "string"
              },
			  "bar": {
                "type": "number"
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}

func TestEmitEntities_ReturnsExpectedEntitiesGivenTwoDataSources(t *testing.T) {
	t.Parallel()
	want := `data_source1: #DataSource: {
    bogus: string
}
data_source2: #DataSource: {
    bogus: string
}`
	inputJSON := []byte(`{
  "format_version": "1.0",
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "data_source1": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
              }
            }
          }
        },
		"data_source2": {
          "version": 0,
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
              }
            }
          }
        }
      }
    }
  }
}`)
	got := ratchet.EmitEntities("bogus", inputJSON)
	if !cmp.Equal(want, got) {
		t.Fatal(cmp.Diff(want, got))
	}
}
