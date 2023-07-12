default:
	false

schemata/%.json.zst:
	mkdir -p schemata
	curl -SsfL https://github.com/cueniform/raw-terraform-registry-providers-schemata/raw/main/$@ -o $@

build/%.json: schemata/%.json.zst
	zstdcat "$^" >"$@"

testdata/script/unit/datasources/%.txtar.tmp: SOURCE=$(shell echo "$@" | sed -e "s/datasources/resources/" -e "s/.tmp//")
testdata/script/unit/datasources/%.txtar.tmp:
	cp $(SOURCE) "$@"

testdata/script/unit/datasources/%.txtar: SOURCE=$(shell echo "$@" | sed "s/datasources/resources/")
testdata/script/unit/datasources/%.txtar: testdata/script/unit/datasources/%.txtar.tmp
	sed \
	-e "s/resource_id/datasource_id/" \
	-e "s/#Resource/#DataSource/" \
	-e "s/resource_schemas/data_source_schemas/" \
	"$^" >"$@"
	rm "$^" # Pending to understand why the file is not being deleted automatically
