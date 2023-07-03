default:
	false

FORCE: ;

build/%.json:
	zstdcat schemata/$*.json.zst > $@

output/%.json: build/%.json
	go run cmd/ratchet/main.go build/$*.json registry.terraform.io/integrations/github

clean: FORCE
	rm -f output/*.json build/*.json

.PHONY: test
test: clean 
test: unitest
test: output/registry.terraform.io_integrations_github_5.9.2.json

.PHONY: unittest
unitest:
	gotestdox
