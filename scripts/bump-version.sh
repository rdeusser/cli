#!/usr/bin/env bash

repo_root=$(git rev-parse --show-toplevel)
version=$(grep -oE "[0-9]+[.][0-9]+[.][0-9]+" "${repo_root}/version/version.go")

major="$(echo "${version}" | tr -d "v" | awk -F '.' '{print $1}')"
minor="$(echo "${version}" | tr -d "v" | awk -F '.' '{print $2}')"
patch="$(echo "${version}" | tr -d "v" | awk -F '.' '{print $3}')"

case $1 in
    "major")
	major=$((major + 1))
	patch=0
	;;
    "minor")
	minor=$((minor + 1))
	patch=0
	;;
    "patch")
	patch=$((patch + 1))
	;;
esac

new_version="${major}.${minor}.${patch}"

if [[ "$new_version" == "$version" ]]; then
    echo "Refusing to bump version. Must pass 'major', 'minor', or 'patch' as an option."
    exit 1
fi

echo "Bumping version from ${version} to ${new_version}"
sed -i "s/${version}/${new_version}/g" "${repo_root}/version/version.go"
git add "${repo_root}/version/version.go"
git commit -vsam "chore: bump version to ${new_version}"
