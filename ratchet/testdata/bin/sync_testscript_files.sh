#!/bin/sh -e

for SRC_FILE in $(ls testdata/script/unit/resources/*.txtar); 
do 
    DEST_FILE=$(echo $SRC_FILE | sed "s/resources/datasources/")
    rm $DEST_FILE
    make $DEST_FILE
done
rm -v testdata/script/unit/datasources/*given*block*.txtar
