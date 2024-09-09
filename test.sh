#!/usr/bin/env bash

for file in testdata/**/*.json; do
  echo "$file"
  case "$file" in
  *"final"*)
    echo "final files aren't used"
    ;;
  *)
    go run . "$file"
    ;;
  esac
done
