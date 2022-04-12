#!/bin/bash

date=$(date +%s_%Y-%m-%d_%H-%M-%S)
name=$1
filename=${date}_${name}.sql

echo "INSERT INTO migrations (name, migrated_at)" >> $filename
echo "  VALUES(\"$filename\", strftime('%s', \"now\"));" >> $filename
