#!/bin/sh
out_file=$1

append_coverage() {
    local profile="$1"
    if [ -f $profile ]; then
        cat $profile | grep -v "mode: count" >> "$out_file"
        rm $profile
    fi
}

echo "mode: count" > "$out_file"

for pkg in $(go list ./...); do
    go test -covermode=count -coverprofile=profile.out "$pkg"
    append_coverage profile.out
done
