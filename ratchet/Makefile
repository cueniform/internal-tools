default:
	false

FORCE: ;

build/%.json:
	zstdcat schemata/$*.json.zst > $@

output/%.json: build/%.json
output/%.json: ADDRESS=$(shell echo "$*" | tr '_' '/' | cut -d/ -f1-3)
output/%.json:
	go run cmd/ratchet/main.go build/$*.json $(ADDRESS)

clean: FORCE
	rm -f output/*.json build/*.json

.PHONY: test
test: clean 
test: unitest
test: output/registry.terraform.io_integrations_github_5.9.2.json

.PHONY: unitest
unitest:
	go test ./...
