default:
	false

schemata/%.json.zst:
	mkdir -p schemata
	curl -SsfL https://github.com/cueniform/raw-terraform-registry-providers-schemata/raw/main/$@ -o $@

build/%.json: schemata/%.json.zst
	zstdcat "$^" >"$@"

.PHONY: test
test:
	go test ./... --tags=fixtures
	go test -race ./...
