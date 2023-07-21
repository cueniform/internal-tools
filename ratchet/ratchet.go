package ratchet

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"github.com/tidwall/gjson"
)

type Ratchet struct {
	OutputLines        []string
	ProviderAddress    string
	ProviderSchemaData []byte
}

// New creates a new Ratchet instance with the provided provider schema data and address.
// It reads the content of the file specified by the 'providerSchemaPath' parameter and
// initializes a new Ratchet struct with the given 'providerAddress' and data of 'providerSchemaPath'.
//
// Parameters:
//
//   - providerSchemaPath: string
//     The path to the file containing the provider schema data.
//
//   - providerAddress: string
//     The address of the provider to emit the entities.
//
// Returns:
//
//   - *Ratchet: A pointer to the newly created Ratchet instance.
//
//   - error: If there is an error reading the provider schema data file or any other
//     error encountered during the initialization, it will be returned.
func New(providerSchemaPath, providerAddress string) (*Ratchet, error) {
	rt := &Ratchet{
		ProviderAddress: providerAddress,
	}
	err := rt.ProviderData(providerSchemaPath)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

// ProviderData returns the data of the provider associated with the Ratchet instance as a slice of bytes.
// It is the public API to allow users to retrieve the provider data stored in the Ratchet instance.
// Returns:
//
//   - []byte: The provider data as a slice of bytes.
//
//   - error: If any error occurs during the retrieval of the provider data, it will be returned as an error.
func (rt *Ratchet) ProviderData(providerSchemaPath string) error {
	providerSchemaData, err := os.ReadFile(providerSchemaPath)
	if err != nil {
		return err
	}
	if len(providerSchemaData) > 0 {
		if providerSchemaData[len(providerSchemaData)-1] == '\n' {
			providerSchemaData = providerSchemaData[:len(providerSchemaData)-1]
		}
	}
	rt.providerData(providerSchemaData)
	return nil
}

// providerData is the internal method that hides the implementation details of how to get the provider data.
// Currently, it uses the Go module gjson which returns an empty string if the keys do not exist.
func (rt *Ratchet) providerData(providerSchemaData []byte) {
	rt.ProviderSchemaData = []byte(gjson.GetBytes(providerSchemaData, "provider_schemas").Get(strings.ReplaceAll(rt.ProviderAddress, ".", "\\.")).String())
}

func (rt *Ratchet) EmitDatasources(dataSourceID string, terraformAttributes gjson.Result) {
	required := map[string]gjson.Result{}
	optional := map[string]gjson.Result{}
	terraformAttributes.ForEach(func(attrID, terraformAttribute gjson.Result) bool {
		if terraformAttribute.Get("required").Bool() {
			required[attrID.String()] = terraformAttribute
			return true
		}
		if terraformAttribute.Get("optional").Bool() {
			optional[attrID.String()] = terraformAttribute
			return true
		}
		if terraformAttribute.Get("computed").Bool() {
			return true
		}
		log.Fatalf("(datasource) %s: Attribute %q is neither required nor optional: %v", dataSourceID, attrID, terraformAttribute)
		return false
	})
	keys := make([]string, 0, len(required))
	for k := range required {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		rt.OutputLines = append(rt.OutputLines, EmitAttribute(key, dataSourceID, "#DataSource", required[key]))
	}
	keys = make([]string, 0, len(optional))
	for k := range optional {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		rt.OutputLines = append(rt.OutputLines, EmitAttribute(key, dataSourceID, "#DataSource", optional[key]))
	}
}

func (rt *Ratchet) EmitResources(resourceID string, terraformBlock gjson.Result) {
	required := map[string]gjson.Result{}
	optional := map[string]gjson.Result{}
	if terraformBlock.Get("attributes").Exists() {
		terraformBlock.Get("attributes").ForEach(func(attrID, terraformAttribute gjson.Result) bool {
			if terraformAttribute.Get("computed").Bool() {
				return true
			}
			if terraformAttribute.Get("required").Bool() {
				required[attrID.String()] = terraformAttribute
				return true
			}
			if terraformAttribute.Get("optional").Bool() {
				optional[attrID.String()] = terraformAttribute
				return true
			}
			log.Fatalf("(resource) %s: Attribute %q is neither required nor optional: %v", resourceID, attrID, terraformAttribute)
			return false
		})
	}
	keys := make([]string, 0, len(required))
	for k := range required {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		rt.OutputLines = append(rt.OutputLines, EmitAttribute(key, resourceID, "#Resource", required[key]))
	}
	keys = make([]string, 0, len(optional))
	for k := range optional {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		rt.OutputLines = append(rt.OutputLines, EmitAttribute(key, resourceID, "#Resource", optional[key]))
	}
	if terraformBlock.Get("block_types").Exists() {
		rt.EmitBlocks(resourceID, terraformBlock.Get("block_types"))
	}
}

func (rt *Ratchet) EmitBlocks(resourceID string, blocks gjson.Result) {
	blocks.ForEach(func(blockID, value gjson.Result) bool {
		rt.OutputLines = append(rt.OutputLines, fmt.Sprintf("%s?: {", blockID.String()))
		if value.Get("block").Get("attributes").Exists() {
			value.Get("block").Get("attributes").ForEach(func(attrID, terraformAttribute gjson.Result) bool {
				if terraformAttribute.Get("required").Bool() {
					rt.OutputLines = append(rt.OutputLines, EmitAttribute(attrID.String(), resourceID, "#Resource", terraformAttribute))
					return true
				}
				if terraformAttribute.Get("optional").Bool() {
					rt.OutputLines = append(rt.OutputLines, EmitAttribute(attrID.String(), resourceID, "#Resource", terraformAttribute))
					return true
				}
				if terraformAttribute.Get("computed").Bool() {
					return true
				}
				log.Fatalf("(block) %s: Attribute %q is neither required nor optional: %v", resourceID, attrID, terraformAttribute)
				return false
			})
		}
		if value.Get("block").Get("block_types").Exists() {
			rt.EmitBlocks(resourceID, value.Get("block").Get("block_types"))
		}
		rt.OutputLines = append(rt.OutputLines, "}")
		return true
	})
}

func (rt *Ratchet) EmitEntities() string {
	gjson.GetBytes(rt.ProviderSchemaData, "data_source_schemas").ForEach(func(dataSourceID, dataSourceValue gjson.Result) bool {
		rt.OutputLines = append(rt.OutputLines, fmt.Sprintf("%s: %s: {", dataSourceID.String(), "#DataSource"))
		if dataSourceValue.Get("block").Get("attributes").Exists() {
			rt.EmitDatasources(dataSourceID.String(), dataSourceValue.Get("block").Get("attributes"))
		}
		rt.OutputLines = append(rt.OutputLines, "}")
		return true
	})
	gjson.GetBytes(rt.ProviderSchemaData, "resource_schemas").ForEach(func(resourceID, resourceValue gjson.Result) bool {
		rt.OutputLines = append(rt.OutputLines, fmt.Sprintf("%s: %s: {", resourceID.String(), "#Resource"))
		if resourceValue.Get("block").Exists() {
			rt.EmitResources(resourceID.String(), resourceValue.Get("block"))
		}
		rt.OutputLines = append(rt.OutputLines, "}")
		return true
	})
	return strings.Join(rt.OutputLines, "\n")
}

func Main() int {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [terraform-provider-schema.json] [provider_address]\n", os.Args[0])
		return 1
	}
	rt, err := New(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	ctx := cuecontext.New()
	fmt.Printf("%#v\n", ctx.CompileString(rt.EmitEntities()))
	return 0
}

func formatPrimitiveTypes(key, value string, typeAttributes gjson.Result) string {
	var output string
	switch {
	case typeAttributes.Get("required").Bool():
		output = fmt.Sprintf("%s!: %s", key, value)
	case typeAttributes.Get("optional").Bool():
		output = fmt.Sprintf("%s?: %s", key, value)
	default:
		log.Fatalf("Attribute %q is neither required or optional: %v", key, typeAttributes)
	}
	return output
}

func formatSetOrListOfComplexObject(key string, objFields gjson.Result) string {
	output := []string{fmt.Sprintf("%s: [..._#%s]", key, key)}
	output = append(output, fmt.Sprintf("_#%s: {", key))
	objFields.ForEach(func(key, value gjson.Result) bool {
		output = append(output, fmt.Sprintf("%s!: %s", key, value.String()))
		return true
	})
	output = append(output, "}")
	return strings.Join(output, "\n")
}

func validType(items []gjson.Result) bool {
	if len(items) != 2 {
		return false
	}
	if items[0].Type != gjson.String {
		return false
	}
	if items[0].String() != "list" && items[0].String() != "set" && items[0].String() != "map" {
		return false
	}
	return true
}

func EmitAttribute(attrID string, entityID string, entityType string, terraformAttribute gjson.Result) string {
	var CUEType string
	attrType := terraformAttribute.Get("type")
	switch attrType.Type {
	// it is a primitive type
	case gjson.String:
		CUEType = attrType.String()
		return formatPrimitiveTypes(attrID, CUEType, terraformAttribute)
	// json schema missing required field.
	case gjson.Null:
		// if it happens means that cue vet did not run
		// or there is a bug in the schema validation
		log.Fatalf("BUG (validator): Attribute field not found in %q.", terraformAttribute.String())
	// it is a complex type
	case gjson.JSON:
		if attrType.IsArray() {
			attrTypeItems := attrType.Array()
			if !validType(attrTypeItems) {
				log.Fatalf("Invalid input for terraform attribute type: %q\n", attrType.String())
			}
			// it is a set or list of a primitive type
			if (attrTypeItems[0].String() == "list" || attrTypeItems[0].String() == "set") && attrTypeItems[1].Type == gjson.String {
				CUEType = fmt.Sprintf("[...%s]", attrTypeItems[1].String())
				return formatPrimitiveTypes(attrID, CUEType, terraformAttribute)
			}
			// it is a map of a primitive type
			if attrTypeItems[0].String() == "map" && attrTypeItems[1].Type == gjson.String {
				CUEType = fmt.Sprintf("[string]: %s", attrTypeItems[1].String())
				return formatPrimitiveTypes(attrID, CUEType, terraformAttribute)
			}
			// it is a set or list of a complex type
			if (attrTypeItems[0].String() == "list" || attrTypeItems[0].String() == "set") && attrTypeItems[1].Type == gjson.JSON {
				return formatSetOrListOfComplexObject(attrID, attrTypeItems[1].Array()[1])
			}
			log.Fatalf("Unable to emit %q %q: cannot translate type for %q. Received -> %q\n", entityType, entityID, attrID, attrType.String())
		}
		log.Fatalf("Unable to emit %q %q: cannot translate type for %q. Received -> %q\n", entityType, entityID, attrID, attrType.String())
	default:
		log.Fatalf("Unkown type %q\n", attrType.Type)
	}
	return ""
}
