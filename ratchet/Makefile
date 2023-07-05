default:
	false

schemata/%.json.zst:
	mkdir -p schemata
	curl -SsfL https://github.com/cueniform/raw-terraform-registry-providers-schemata/raw/main/$@ -o $@

build/%.json: schemata/%.json.zst
	zstdcat "$^" >"$@"

output/%.cue: ADDRESS=$(shell echo "$*" | tr '_' '/' | cut -d/ -f1-3)
output/%.cue: build/%.json
	go run cmd/ratchet/main.go "$^" $(ADDRESS) >"$@"

.PHONY: clean
clean:
	rm -v output/*.cue build/*.json schemata/*.json.zst || true

.PHONY: test
test: clean 
test: output/registry.terraform.io_integrations_github_5.9.2.cue
test: output/registry.terraform.io_hashicorp_aws_5.6.1.cue
