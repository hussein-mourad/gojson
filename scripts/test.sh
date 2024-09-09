#!/usr/bin/env bash

for file in testdata/**/*.json; do
  echo "$file"
  case "$file" in
  *"final"*)
    go run main.go "$file"
    # echo "final files aren't used"
    ;;
  *)
    # go run main.go "$file"
    ;;
  esac
done
