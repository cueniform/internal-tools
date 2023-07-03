package ratchet

import (
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/gjson"
)

func ValidTFSchemaVersion(schemaVersion gjson.Result) bool {
	switch schemaVersion.Type {
	case gjson.String:
		return schemaVersion.String() == "1.0"
	default:
		return false
	}
}

func EmitEntities(providerID string, JSONData []byte) string {
	output := []string{}
	gjson.GetBytes(JSONData, "provider_schemas").ForEach(func(key, value gjson.Result) bool {
		if key.String() == providerID {
			// providerID := key.String()
			// fmt.Println("emit entities for " + providerID)
			value.ForEach(func(key, value gjson.Result) bool {
				if key.String() == "data_source_schemas" {
					// fmt.Println("datasources:")
					value.ForEach(func(key, value gjson.Result) bool {
						output = append(output, fmt.Sprintf("%s: #DataSource: {", key.String()))
						value.Get("block").Get("attributes").ForEach(func(key, value gjson.Result) bool {
							if value.Get("computed").Bool() {
								return true
							}
							tfType := value.Get("type")
							if tfType.IsArray() {
								if tfType.Array()[0].String() == "set" {
									if tfType.Array()[1].IsArray() {
										if tfType.Array()[1].Array()[0].String() == "object" {
											output = append(output, fmt.Sprintf("    %s: [..._#%s]", key.String(), key.String()))
										}
										output = append(output, fmt.Sprintf("    _#%s: {", key.String()))
										for k, v := range tfType.Array()[1].Array()[1].Map() {
											output = append(output, fmt.Sprintf("        %s: %s", k, v))
										}
										output = append(output, "    }")
									}
								}
								return true
							}
							//fmt.Println(tfType.Type, tfType.Raw, tfType.String())
							CUEType, err := ConvertTerraformType(tfType)
							if err != nil {
								log.Fatal(err)
							}
							switch {
							case value.Get("required").Bool():
								output = append(output, fmt.Sprintf("    %s!: %s", key.String(), CUEType))
							case value.Get("optional").Bool():
								output = append(output, fmt.Sprintf("    %s?: %s", key.String(), CUEType))
							default:
								output = append(output, fmt.Sprintf("    %s: %s", key.String(), CUEType))
							}
							return true
						})
						output = append(output, "}")
						return true
					})
				}
				if key.String() == "resource_schemas" {
					// fmt.Println("datasources:")
					value.ForEach(func(key, value gjson.Result) bool {
						output = append(output, fmt.Sprintf("%s: #Resource: {", key.String()))
						value.Get("block").Get("attributes").ForEach(func(key, value gjson.Result) bool {
							if value.Get("computed").Bool() {
								return true
							}
							tfType := value.Get("type")
							if tfType.IsArray() {
								if tfType.Array()[0].String() == "set" {
									if tfType.Array()[1].IsArray() {
										if tfType.Array()[1].Array()[0].String() == "object" {
											output = append(output, fmt.Sprintf("    %s: [..._#%s]", key.String(), key.String()))
										}
										output = append(output, fmt.Sprintf("    _#%s: {", key.String()))
										for k, v := range tfType.Array()[1].Array()[1].Map() {
											output = append(output, fmt.Sprintf("        %s: %s", k, v))
										}
										output = append(output, "    }")
									}
								}
								return true
							}
							CUEType, err := ConvertTerraformType(tfType)
							if err != nil {
								log.Fatal(err)
							}
							switch {
							case value.Get("required").Bool():
								output = append(output, fmt.Sprintf("    %s!: %s", key.String(), CUEType))
							case value.Get("optional").Bool():
								output = append(output, fmt.Sprintf("    %s?: %s", key.String(), CUEType))
							default:
								output = append(output, fmt.Sprintf("    %s: %s", key.String(), CUEType))
							}
							return true
						})
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

func ConvertTerraformType(TFType gjson.Result) (string, error) {
	switch TFType.Type {
	case gjson.String:
		return TFType.String(), nil
	case gjson.JSON:
		if TFType.IsArray() {
			TFTypeItems := TFType.Array()
			if len(TFTypeItems) != 2 {
				return "", fmt.Errorf("%d is invalid number of items. Expecting 2", len(TFTypeItems))
			}
			switch TFTypeItems[0].String() {
			case "list", "set":
				if !TFTypeItems[1].IsArray() {
					return fmt.Sprintf("[...%s]", TFTypeItems[1].String()), nil
				}
				return "", fmt.Errorf("i dont know what to return")
			case "map":
				return fmt.Sprintf("[string]: %s", TFTypeItems[1].String()), nil
			default:
				return "", fmt.Errorf("not sure what to do with %q", TFTypeItems[0].String())
			}
		}
		return "", fmt.Errorf("i dont know what to return")
	default:
		return "", fmt.Errorf("unkown type %q", TFType.Type)
	}
}
