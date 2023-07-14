#!/bin/sh -e

src_dir="testdata/script/unit/entities"
for type_camel in Resource DataSource;
do
    if [[ "${type_camel}" == "Resource" ]]; then
        type_lower="resource"
    elif [[ "${type_camel}" == "DataSource" ]]; then
        type_lower="data_source"
    fi
    dest_dir="testdata/script/unit/${type_lower}s"
    rm -vf "${dest_dir}"/*.common.txtar
    cp -v "${src_dir}"/*.common.txtar "${dest_dir}"/
    sed -i.bak \
        -e "s/ENTITY_LOWERCASE/${type_lower}/" \
        -e "s/ENTITY_CAMELCASE/${type_camel}/" \
        ${dest_dir}/*.common.txtar
    rm -vf "${dest_dir}"/*.common.txtar.bak
done
