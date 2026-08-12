#!/usr/bin/env sh
# Install the latest RiceBench release into ~/.local/bin.
#
#   curl -sSfL https://raw.githubusercontent.com/gabrielassisxyz/ricebench/main/install.sh | sh
#
# Set RICEBENCH_VERSION to pin a tag, and RICEBENCH_BIN_DIR to install elsewhere.
set -eu

repo=gabrielassisxyz/ricebench
bin_dir=${RICEBENCH_BIN_DIR:-$HOME/.local/bin}

for tool in curl tar uname; do
    command -v "$tool" >/dev/null 2>&1 || {
        echo "install.sh: missing required tool: $tool" >&2
        exit 1
    }
done

os=$(uname -s)
arch=$(uname -m)
if [ "$os" != Linux ] || [ "$arch" != x86_64 ]; then
    echo "install.sh: releases currently cover linux x86_64 only, this is $os $arch." >&2
    echo "Build from source instead: bin/build-web, then go build ./cmd/ricebench" >&2
    exit 1
fi

version=${RICEBENCH_VERSION:-}
if [ -z "$version" ]; then
    # Resolving through the redirect avoids parsing the releases API, which is rate limited for
    # unauthenticated callers and would fail exactly when several people try this at once.
    version=$(curl -sSfL -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest" | sed 's|.*/tag/||')
fi
[ -n "$version" ] || {
    echo "install.sh: could not resolve the latest version." >&2
    exit 1
}

number=${version#v}
archive="ricebench_${number}_linux_amd64.tar.gz"
url="https://github.com/$repo/releases/download/$version/$archive"

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

echo "Downloading $archive"
curl -sSfL "$url" -o "$work/$archive"

echo "Verifying the checksum"
curl -sSfL "https://github.com/$repo/releases/download/$version/checksums.txt" -o "$work/checksums.txt"
(cd "$work" && grep " $archive\$" checksums.txt | sha256sum -c -) || {
    echo "install.sh: checksum mismatch, refusing to install." >&2
    exit 1
}

tar -xzf "$work/$archive" -C "$work" ricebench
mkdir -p "$bin_dir"
install -m 0755 "$work/ricebench" "$bin_dir/ricebench"

echo "Installed ricebench $version to $bin_dir/ricebench"
case ":$PATH:" in
    *":$bin_dir:"*) ;;
    *) echo "Note: $bin_dir is not on PATH." ;;
esac
