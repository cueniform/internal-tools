default:
	false

build/%.json: schemata/%.json.zst
	zstdcat "$^" >"$@"

output/%.cue: ADDRESS=$(shell echo "$*" | tr '_' '/' | cut -d/ -f1-3)
output/%.cue: build/%.json
	go run cmd/ratchet/main.go "$^" $(ADDRESS) >"$@"

.PHONY: clean
clean:
	rm -vf output/*.cue build/*.json

.PHONY: test
test: clean 
test: unitest
test: output/registry.terraform.io_integrations_github_5.9.2.cue

.PHONY: unitest
unitest:
	go test ./...
