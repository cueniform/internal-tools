package ratchet_test

import (
	"fmt"
	"testing"

	"github.com/cueniform/internal-tools/ratchet"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
)

func TestValidTFSChemaVersion_ReturnsTrueGivenResultStringWithKnownVersion(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenResultStringWithUnknownVersion(t *testing.T) {
	t.Parallel()
	want := false
	input := gjson.Result{
		Type: gjson.String,
		Raw:  "666",
		Str:  "666",
	}
	got := ratchet.ValidTFSchemaVersion(input)
	if want != got {
		t.Fatalf("want valid %t, got %t", want, got)
	}
}

func TestValidTFSChemaVersion_ReturnsFalseGivenResultNumber(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenResultFalse(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenResultTrue(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenResultNull(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenResultJSON(t *testing.T) {
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

func TestValidTFSChemaVersion_ReturnsFalseGivenUnkownResultType(t *testing.T) {
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

func TestEmitEntities_SkipsItemsWithComputedTrueGivenDataSource(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {

}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
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

func TestEmitEntities_SkipsItemsWithComputedTrueGivenResource(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {

}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
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
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
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

func TestEmitEntities_ReturnsExpectedStringGivenResourceWithRequiredAttribute(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {
    bogus!: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
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
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
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

func TestEmitEntities_ReturnsExpectedStringGivenResourceWithOptionalAttribute(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {
    bogus?: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
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

func TestEmitEntities_ReturnsExpectedStringGivenDataSourceWithPrimitiveTypes(t *testing.T) {
	testCases := []struct {
		desc     string
		rawInput string
		want     string
	}{
		{
			desc:     "String",
			rawInput: `"string"`,
			want:     "string",
		},
		{
			desc:     "Number",
			rawInput: `"number"`,
			want:     "number",
		},
		{
			desc:     "Boolean",
			rawInput: `"bool"`,
			want:     "bool",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			wantTmpl := "bogus: #DataSource: {\n    bogus!: %s\n}"
			want := fmt.Sprintf(wantTmpl, tC.want)
			inputTmpl := `{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
                "type": %s,
                "required": true
              }
            }
          }
        }
      }
    }
  }
}`
			inputJSON := []byte(fmt.Sprintf(inputTmpl, tC.rawInput))
			got := ratchet.EmitEntities("bogus", inputJSON)
			if !cmp.Equal(want, got) {
				t.Fatal(cmp.Diff(want, got))
			}
		})
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenResourceWithPrimitiveTypes(t *testing.T) {
	testCases := []struct {
		desc     string
		rawInput string
		want     string
	}{
		{
			desc:     "String",
			rawInput: `"string"`,
			want:     "string",
		},
		{
			desc:     "Number",
			rawInput: `"number"`,
			want:     "number",
		},
		{
			desc:     "Boolean",
			rawInput: `"bool"`,
			want:     "bool",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			wantTmpl := "bogus: #Resource: {\n    bogus!: %s\n}"
			want := fmt.Sprintf(wantTmpl, tC.want)
			inputTmpl := `{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
                "type": %s,
                "required": true
              }
            }
          }
        }
      }
    }
  }
}`
			inputJSON := []byte(fmt.Sprintf(inputTmpl, tC.rawInput))
			got := ratchet.EmitEntities("bogus", inputJSON)
			if !cmp.Equal(want, got) {
				t.Fatal(cmp.Diff(want, got))
			}
		})
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenDataSourceWithSetListAndMapOfPrimitiveTypes(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc     string
		rawInput string
		want     string
	}{
		{
			desc:     "List of string",
			rawInput: `["list", "string"]`,
			want:     "[...string]",
		},
		{
			desc:     "List of number",
			rawInput: `["list", "number"]`,
			want:     "[...number]",
		},
		{
			desc:     "List of bool",
			rawInput: `["list", "bool"]`,
			want:     "[...bool]",
		},
		{
			desc:     "Set of string",
			rawInput: `["set", "string"]`,
			want:     "[...string]",
		},
		{
			desc:     "Set of number",
			rawInput: `["set", "number"]`,
			want:     "[...number]",
		},
		{
			desc:     "Set of bool",
			rawInput: `["set", "bool"]`,
			want:     "[...bool]",
		},
		{
			desc:     "Map of string",
			rawInput: `["map", "string"]`,
			want:     "[string]: string",
		},
		{
			desc:     "Map of number",
			rawInput: `["map", "number"]`,
			want:     "[string]: number",
		},
		{
			desc:     "Map of bool",
			rawInput: `["map", "bool"]`,
			want:     "[string]: bool",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			wantTmpl := "bogus: #DataSource: {\n    bogus!: %s\n}"
			want := fmt.Sprintf(wantTmpl, tC.want)
			inputTmpl := `{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
                "type": %s,
                "required": true
              }
            }
          }
        }
      }
    }
  }
}`
			inputJSON := []byte(fmt.Sprintf(inputTmpl, tC.rawInput))
			got := ratchet.EmitEntities("bogus", inputJSON)
			if !cmp.Equal(want, got) {
				t.Fatal(cmp.Diff(want, got))
			}
		})
	}
}

func TestEmitEntities_ReturnsExpectedStringGivenResourceWithSetListAndMapOfPrimitiveTypes(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc     string
		rawInput string
		want     string
	}{
		{
			desc:     "List of string",
			rawInput: `["list", "string"]`,
			want:     "[...string]",
		},
		{
			desc:     "List of number",
			rawInput: `["list", "number"]`,
			want:     "[...number]",
		},
		{
			desc:     "List of bool",
			rawInput: `["list", "bool"]`,
			want:     "[...bool]",
		},
		{
			desc:     "Set of string",
			rawInput: `["set", "string"]`,
			want:     "[...string]",
		},
		{
			desc:     "Set of number",
			rawInput: `["set", "number"]`,
			want:     "[...number]",
		},
		{
			desc:     "Set of bool",
			rawInput: `["set", "bool"]`,
			want:     "[...bool]",
		},
		{
			desc:     "Map of string",
			rawInput: `["map", "string"]`,
			want:     "[string]: string",
		},
		{
			desc:     "Map of number",
			rawInput: `["map", "number"]`,
			want:     "[string]: number",
		},
		{
			desc:     "Map of bool",
			rawInput: `["map", "bool"]`,
			want:     "[string]: bool",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			wantTmpl := "bogus: #Resource: {\n    bogus!: %s\n}"
			want := fmt.Sprintf(wantTmpl, tC.want)
			inputTmpl := `{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "bogus": {
                "type": %s,
                "required": true
              }
            }
          }
        }
      }
    }
  }
}`
			inputJSON := []byte(fmt.Sprintf(inputTmpl, tC.rawInput))
			got := ratchet.EmitEntities("bogus", inputJSON)
			if !cmp.Equal(want, got) {
				t.Fatal(cmp.Diff(want, got))
			}
		})
	}
}

func TestEmitEntities_ReturnsRequiredAttributesFirstGivenDataSourceWithTwoAttributes(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    bar!: number
    foo?: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "foo": {
                "type": "string",
                "optional": true
              },
              "bar": {
                "type": "number",
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

func TestEmitEntities_ReturnsRequiredAttributesFirstGivenResourceWithTwoAttributes(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {
    bar!: number
    foo?: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "foo": {
                "type": "string",
                "optional": true
              },
              "bar": {
                "type": "number",
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

func TestEmitEntities_ReturnsExpectedEntitiesGivenTwoDataSources(t *testing.T) {
	t.Parallel()
	want := `data_source1: #DataSource: {
    bogus!: string
}
data_source2: #DataSource: {
    bogus!: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "data_source1": {
          "block": {
            "attributes": {
              "bogus": {
                "type": "string",
                "required" true
              }
            }
          }
        },
        "data_source2": {
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
                "required" true
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

func TestEmitEntities_ReturnsExpectedEntitiesGivenTwoResources(t *testing.T) {
	t.Parallel()
	want := `resource1: #Resource: {
    bogus!: string
}
resource2: #Resource: {
    bogus!: string
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "resource1": {
          "block": {
            "attributes": {
              "bogus": {
                "type": "string",
                "required" true
              }
            }
          }
        },
        "resource2": {
          "block": {
            "attributes": {
              "bogus": {
                "type": "string"
                "required" true
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

func TestEmitEntities_ReturnsExpectedEntityGivenDataSourceWithComplexObject(t *testing.T) {
	t.Parallel()
	want := `bogus: #DataSource: {
    myComplexObj: [..._#myComplexObj]
    _#myComplexObj: {
        field1!: string
        field2!: string
    }
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "data_source_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "myComplexObj": {
                "type": ["set",[
                    "object",
                    {
                      "field1": "string",
                      "field2": "string"
                    }
                  ]
				],
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

func TestEmitEntities_ReturnsExpectedEntityGivenResourceWithComplexObject(t *testing.T) {
	t.Parallel()
	want := `bogus: #Resource: {
    myComplexObj: [..._#myComplexObj]
    _#myComplexObj: {
        field1!: string
        field2!: string
    }
}`
	inputJSON := []byte(`{
  "provider_schemas": {
    "bogus": {
      "resource_schemas": {
        "bogus": {
          "block": {
            "attributes": {
              "myComplexObj": {
                "type": ["set",[
                    "object",
                    {
                      "field1": "string",
                      "field2": "string"
                    }
                  ]
				],
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
