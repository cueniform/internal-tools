default:
	false

schemata/%.json.zst:
	mkdir -p schemata
	curl -SsfL https://github.com/cueniform/raw-terraform-registry-providers-schemata/raw/main/$@ -o $@

build/%.json: schemata/%.json.zst
	zstdcat "$^" >"$@"

.PHONY: test
test: unitest
test: integration

.PHONY: ci
ci:
	act --rm -C ../

.PHONY: unitest
unitest:
	go test ./... --tags=fixtures
	go test -race ./...

.PHONY: integration
integration:
	test -z "$(which ratchet)" || rm "$(which ratchet)"
	rm -f cmd/ratchet/ratchet
	go build -o cmd/ratchet/ratchet cmd/ratchet/main.go
	PATH="$(PATH):$(PWD)/cmd/ratchet" testscript testdata/script/integration/*
	rm cmd/ratchet/ratchet

.PHONY: deps
deps:
	go install github.com/rogpeppe/go-internal/cmd/testscript@latest
