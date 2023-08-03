package ratchet

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cuelang.org/go/cue/cuecontext"
	"github.com/tidwall/gjson"
)

// Ratchet stores the CLI runtime information
type Ratchet struct {
	Debug           io.Writer
	Output          []string
	ProviderAddress string
	ProviderSchema  string
}

// New creates a new Ratchet instance with the provided provider schema data and address.
func New(providerSchemaPath, providerAddress string) (*Ratchet, error) {
	rt := &Ratchet{
		ProviderAddress: providerAddress,
		Debug:           io.Discard,
	}
	err := rt.ProviderData(providerSchemaPath)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

// ProviderData stores the data of the provider associated with the Ratchet instance as a slice of bytes.
func (rt *Ratchet) ProviderData(providerSchemaPath string) error {
	providerSchemaData, err := os.ReadFile(providerSchemaPath)
	if err != nil {
		return err
	}
	pathEscaped := strings.ReplaceAll(rt.ProviderAddress, ".", "\\.")
	rt.ProviderSchema = gjson.GetBytes(providerSchemaData, "provider_schemas").Get(pathEscaped).String()
	return nil
}

// FormatCUEKey returns the translated string from gjson.Result to CUE value
// or returns an error if the key is neither required nor optional.
func FormatCUEKey(keyID string, key gjson.Result) (string, error) {
	switch {
	case key.Get("required").Bool():
		return fmt.Sprintf("%s!:", keyID), nil
	case key.Get("optional").Bool():
		return fmt.Sprintf("%s?:", keyID), nil
	default:
		return "", fmt.Errorf("%s is neither required nor optional. Received: %v", keyID, key)
	}
}

// EmitEntities generates and emits entities for data sources and resources based on
// the Ratchet's provider schema.
func (rt *Ratchet) EmitEntities() {
	rt.Output = append(rt.Output, `import "list"`)
	rt.Output = append(rt.Output, "\n")
	gjson.Get(rt.ProviderSchema, "data_source_schemas").ForEach(func(dataSourceID, dataSourceValue gjson.Result) bool {
		rt.Output = append(rt.Output, fmt.Sprintf("%s?: %s?: {\n", dataSourceID, "#DataSource"))
		rt.EmitDataSource(dataSourceValue)
		rt.Output = append(rt.Output, "}\n")
		return true
	})
	gjson.Get(rt.ProviderSchema, "resource_schemas").ForEach(func(resourceID, resourceValue gjson.Result) bool {
		rt.Output = append(rt.Output, fmt.Sprintf("%s?: %s?: {\n", resourceID, "#Resource"))
		rt.EmitResource(resourceValue)
		rt.Output = append(rt.Output, "}\n")
		return true
	})
}

func (rt *Ratchet) EmitBlockAttributes(attributes gjson.Result) {
	attributes.Get("block.block_types").ForEach(func(blockID, blockValue gjson.Result) bool {
		rt.Output = append(rt.Output, fmt.Sprintf("%s?: {\n", blockID))
		rt.EmitResource(blockValue)
		rt.Output = append(rt.Output, "}\n")
		return true
	})
}

func (rt *Ratchet) EmitSingleAttributes(attributes gjson.Result) {
	required := []map[string]gjson.Result{}
	optional := []map[string]gjson.Result{}
	attributes.Get("block.attributes").ForEach(func(attributeID, attributeValue gjson.Result) bool {
		if attributeValue.Get("nested_type").Exists() {
			return true
		}
		if attributeValue.Get("computed").Bool() {
			return true
		}
		if attributeValue.Get("required").Bool() {
			required = append(required, map[string]gjson.Result{attributeID.String(): attributeValue})
			return true
		}
		if attributeValue.Get("optional").Bool() {
			optional = append(optional, map[string]gjson.Result{attributeID.String(): attributeValue})
			return true
		}
		log.Fatalf("%s is neither required nor optional", attributeID.String())
		return false
	})
	for _, r := range required {
		for k, v := range r {
			rt.EmitAttribute(k, v)
		}
	}
	for _, o := range optional {
		for k, v := range o {
			rt.EmitAttribute(k, v)
		}
	}
}

func (rt *Ratchet) EmitNestedAttribute(attributeID string, attributeValue gjson.Result) {
	if !attributeValue.Get("optional").Bool() && attributeValue.Get("computed").Bool() {
		return
	}
	CUEKey, err := FormatCUEKey(attributeID, attributeValue)
	if err != nil {
		log.Fatal(err)
	}
	rt.Output = append(rt.Output, CUEKey)
	if attributeValue.Get("nested_type.attributes").Exists() {
		attributeValue.Get("nested_type.attributes").ForEach(func(key, value gjson.Result) bool {
			rt.EmitNestedAttribute(key.String(), value)
			return true
		})
	}
	rt.ConvertType(attributeValue.Get("type"))
}

func (rt *Ratchet) EmitNestedAttributes(attributes gjson.Result) {
	attributes.Get("block.attributes").ForEach(func(attributeID, attributeValue gjson.Result) bool {
		if attributeValue.Get("nested_type.attributes").Exists() {
			rt.Output = append(rt.Output, fmt.Sprintf("%s?: {\n", attributeID))
			attributeValue.Get("nested_type.attributes").ForEach(func(key, value gjson.Result) bool {
				rt.EmitNestedAttribute(key.String(), value)
				return true
			})
			rt.Output = append(rt.Output, "}\n")
		}

		return true
	})
}

func (rt *Ratchet) EmitResource(attributes gjson.Result) {
	rt.EmitSingleAttributes(attributes)
	rt.EmitBlockAttributes(attributes)
	rt.EmitNestedAttributes(attributes)
}

func (rt *Ratchet) EmitDataSource(attributes gjson.Result) {
	required := []map[string]gjson.Result{}
	optional := []map[string]gjson.Result{}
	attributes.Get("block.attributes").ForEach(func(attributeID, attributeValue gjson.Result) bool {
		if attributeValue.Get("nested_type").Exists() {
			return true
		}
		if !attributeValue.Get("optional").Bool() && attributeValue.Get("computed").Bool() {
			return true
		}
		if attributeValue.Get("required").Bool() {
			required = append(required, map[string]gjson.Result{attributeID.String(): attributeValue})
			return true
		}
		if attributeValue.Get("optional").Bool() {
			optional = append(optional, map[string]gjson.Result{attributeID.String(): attributeValue})
			return true
		}
		log.Fatalf("%s is neither required nor optional", attributeID.String())
		return false
	})
	for _, r := range required {
		for k, v := range r {
			rt.EmitAttribute(k, v)
		}
	}
	for _, o := range optional {
		for k, v := range o {
			rt.EmitAttribute(k, v)
		}
	}
}

func (rt *Ratchet) EmitAttribute(attributeID string, attributeValue gjson.Result) {
	CUEKey, err := FormatCUEKey(attributeID, attributeValue)
	if err != nil {
		log.Fatal(err)
	}
	rt.Output = append(rt.Output, CUEKey)
	rt.ConvertType(attributeValue.Get("type"))
}

// String returns the string representation of Ratchet instance.
func (rt *Ratchet) String() string {
	return strings.Join(rt.Output, " ")
}

func (rt *Ratchet) ConvertType(attributeType gjson.Result) {
	switch attributeType.Type {
	case gjson.String:
		s := attributeType.String()
		if s == "dynamic" {
			s = "_"
		}
		rt.Output = append(rt.Output, s)
		rt.Output = append(rt.Output, "\n")
	case gjson.JSON:
		if !attributeType.IsArray() {
			log.Fatalf("Type must be string or array. Got: %q", attributeType.String())
		}
		if attributeType.Array()[0].String() == "list" {
			rt.Output = append(rt.Output, "[...")
			rt.ConvertType(attributeType.Array()[1])
			rt.Output = append(rt.Output, "]\n")
		}
		if attributeType.Array()[0].String() == "set" {
			rt.Output = append(rt.Output, "[...")
			rt.ConvertType(attributeType.Array()[1])
			rt.Output = append(rt.Output, "] & list.UniqueItems()\n")
		}
		if attributeType.Array()[0].String() == "map" {
			rt.Output = append(rt.Output, "{[string]:")
			rt.ConvertType(attributeType.Array()[1])
			rt.Output = append(rt.Output, "}\n")
		}
		if attributeType.Array()[0].String() == "object" {
			rt.Output = append(rt.Output, "{\n")
			attributeType.Array()[1].ForEach(func(key, value gjson.Result) bool {
				rt.Output = append(rt.Output, fmt.Sprintf("%s!: {", key))
				rt.ConvertType(value)
				rt.Output = append(rt.Output, "}\n")
				return true
			})
			rt.Output = append(rt.Output, "}\n")
		}
	}
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
	rt.EmitEntities()
	ctx := cuecontext.New()
	v := ctx.CompileString(fmt.Sprintln(rt))
	if v.Err() != nil {
		fmt.Fprintln(os.Stderr, v.Err())
		fmt.Fprintln(rt.Debug, fmt.Sprint(rt))
		return 1
	}
	fmt.Printf("%#v\n", v)
	return 0
}
