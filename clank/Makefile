default:
	false

test:
	ls schemata/registry.terraform.io_hashicorp_*.json.zst | shuf | while read I ; do \
		echo "$$(date): $$I"; \
		zstdcat $$I \
		| cue vet -c .:clank json: - -d '#TerraformProvidersSchema' >"errors.$$(basename $$I).txt" 2>&1 \
		&& rm "errors.$$(basename $$I).txt"; \
		echo "$$(date): done"; \
	done

errors.%.txt:
	zstdcat schemata/$* \
	| cue vet -c .:clank json: - -d '#TerraformProvidersSchema' >"$@" 2>&1 \
	&& rm "$@"; \
	true
