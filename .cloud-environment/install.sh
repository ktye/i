#!/bin/bash
set -euo pipefail
setup=/workspace/.cloud-setup
mkdir -p "$setup/bin" "$setup/gocache" "$setup/gomodcache"
if ! "$setup/go/bin/go" version 2>/dev/null | grep -q 'go1.27.2 linux/amd64'; then
  archive=$(mktemp /tmp/cloud-go.XXXXXX.tar.gz)
  curl -fsSL https://go.dev/dl/go1.27.2.linux-amd64.tar.gz -o "$archive"
  echo "ecbadb99091a3f46e31f5f934b068b1864eafa7995211b39eaddf76996045fe5  $archive" | sha256sum -c -
  tar -xzf "$archive" -C "$setup"
  rm "$archive"
fi
export PATH="$setup/go/bin:$PATH"
export GOCACHE="$setup/gocache" GOMODCACHE="$setup/gomodcache" GOTOOLCHAIN=local
cd /workspace/i/_/i5
go build -mod=readonly -o "$setup/bin/k" ./k
# TestK regenerates tracked readme.md: run unchanged tests on a disposable copy.
testdir=$(mktemp -d /tmp/i5-tests.XXXXXX)
trap 'rm -rf "$testdir"' EXIT
cp -a . "$testdir/"
cd "$testdir"
go test -mod=readonly -count=1 -v ./... > "$setup/i5-tests.log" 2>&1
cat "$setup/i5-tests.log"
python3 - <<'PY'
import subprocess
p = subprocess.run(['/workspace/.cloud-setup/bin/k'], input='1+2\n+/!5\n', text=True, capture_output=True, check=True)
assert [s.strip() for s in p.stdout.splitlines() if s.strip()] == ['ktye/k', '3', '10'], p.stdout
print('CLI smoke: 1+2 => 3; +/!5 => 10')
PY
