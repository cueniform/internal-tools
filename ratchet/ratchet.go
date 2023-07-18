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

func formatPrimitiveTypes(key, value string, typeAttributes gjson.Result) string {
	var output string
	switch {
	case typeAttributes.Get("required").Bool():
		output = fmt.Sprintf("    %s!: %s", key, value)
	case typeAttributes.Get("optional").Bool():
		output = fmt.Sprintf("    %s?: %s", key, value)
	default:
		log.Fatalf("Attribute %q is neither required or optional", key)
	}
	return output
}

func formatSetOrListOfComplexObject(key string, objFields gjson.Result) string {
	output := []string{fmt.Sprintf("    %s: [..._#%s]", key, key)}
	output = append(output, fmt.Sprintf("    _#%s: {", key))
	objFields.ForEach(func(key, value gjson.Result) bool {
		output = append(output, fmt.Sprintf("        %s!: %s", key, value.String()))
		return true
	})
	output = append(output, "    }")
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
	// json schema missing required field
	case gjson.Null:
		log.Fatalf("Attribute field not found in %q", terraformAttribute.String())
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

func EmitAttributes(entityID string, entityType string, terraformAttributes gjson.Result) string {
	output := []string{}
	required := map[string]gjson.Result{}
	optional := map[string]gjson.Result{}
	terraformAttributes.ForEach(func(attrID, terraformAttribute gjson.Result) bool {
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
		log.Fatalf("Attribute %q is neither required or optional", attrID)
		return true
	})
	keys := make([]string, 0, len(required))
	for k := range required {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		output = append(output, EmitAttribute(key, entityID, entityType, required[key]))
	}
	keys = make([]string, 0, len(optional))
	for k := range optional {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, key := range keys {
		output = append(output, EmitAttribute(key, entityID, entityType, optional[key]))
	}
	return strings.Join(output, "\n")
}

func EmitEntities(providerID string, JSONData []byte) string {
	output := []string{}
	gjson.GetBytes(JSONData, "provider_schemas").ForEach(func(providerAddress, providerValue gjson.Result) bool {
		if providerAddress.String() == providerID {
			providerValue.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "data_source_schemas" {
					value.ForEach(func(datasourceID, datasourceValue gjson.Result) bool {
						output = append(output, fmt.Sprintf("%s: %s: {", datasourceID.String(), "#DataSource"))
						output = append(output, EmitAttributes(datasourceID.String(), "#DataSource", datasourceValue.Get("block").Get("attributes")))
						output = append(output, "}")
						return true
					})
				}
				if key.String() == "resource_schemas" {
					value.ForEach(func(resourceID, value gjson.Result) bool {
						output = append(output, fmt.Sprintf("%s: %s: {", resourceID.String(), "#Resource"))
						output = append(output, EmitAttributes(resourceID.String(), "#Resource", value.Get("block").Get("attributes")))
						output = append(output, "}")
						return true
					})
				}
				return true
			})
		}
		return true
	})
	return strings.Join(output, "\n")
}

func Main() int {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [terraform-provider-schema.json] [provider_address]\n", os.Args[0])
		return 1
	}
	JSONData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	CUESchema := EmitEntities(os.Args[2], JSONData)
	ctx := cuecontext.New()
	fmt.Printf("%#v\n", ctx.CompileString(CUESchema))
	return 0
}
