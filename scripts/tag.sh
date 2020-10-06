#!/usr/bin/env bash

repo_root=$(git rev-parse --show-toplevel)
version=$(grep -oE "[0-9]+[.][0-9]+[.][0-9]+" "${repo_root}/version/version.go")
remote=$(git remote -v | awk '{print $1}' | head -n 1)

git tag -a "v${version}" -m "v${version}"
git push "$remote" "v${version}"`
