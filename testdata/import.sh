#!/bin/sh

if [ "$#" -ne "2" ]
then
    echo "usage: import <sqlite_target_db> <records_source_csv>"
    exit 2
fi

# ./import.sh ~/.ttt.sqlite testdata/someRecordsWithStartNew.csv
sqlite3 -csv $1 ".import $2 records --skip 1"