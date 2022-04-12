#!/bin/bash


MIGRATIONS=`ls db/migrations/*.sql`
for m in $MIGRATIONS; do
  echo $m
  sqlite3 db/kees.sqlite < $m
done
