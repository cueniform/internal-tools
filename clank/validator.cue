package clank

#TerraformProvidersSchema: {
	format_version!: =~"^1\\."
	provider_schemas?: [providerAddress=string]: {
		provider?: _#Schema
		resource_schemas?: [resourceName=string]:      _#Schema
		data_source_schemas?: [dataSourceName=string]: _#Schema
	}
}
_#Schema: {
	version?: int
	block?:   _#Block
}
_#Block: {
	description?:      string
	description_kind?: "plain" | "markdown"
	deprecated?:       bool
	attributes?: [attributeName=string]: {
		type?: "string" | "number" | "bool" | "dynamic" |
			["map", "string"] | ["map", "bool"] | ["map", "number"] |
			[ "list", "number"] | ["list", "string"] | ["list", ["map", "string"]] | ["list", ["list", "number"]] |
			["set", "number"] | ["set", "string"] | ["set", [ "map", "string"]] |
					["object", _] | ["list", ["object", _]] | ["set", ["object", _]]
		description?:      string
		description_kind?: "plain" | "markdown"
		required?:         bool
		optional?:         bool
		computed?:         bool
		sensitive?:        bool
		deprecated?:       bool
		nested_type?:      "foo" // TODO: investigate hashicorp/awscc's use of this
	}
	block_types?: [blockName=string]: {
		nesting_mode?: "single" | "list" | "set" | "map"
		block?:        _#Block
		min_items?:    int
		max_items?:    int
	}
}

#In: #TerraformProvidersSchema
