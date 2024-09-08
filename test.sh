#!/usr/bin/env bash

for file in testdata/**/*.json; do
  echo "$file"
  go run . "$file"
done
