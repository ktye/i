# Codex cloud environment

This setup prepares the K library and CLI in `_/i5`. It installs Go 1.27.2 for Linux amd64, verifies the download with SHA-256, builds the CLI, runs the existing tests, and checks two expressions through the CLI.

The cloud checkout must be at `/workspace/i`. Setup requires Bash, curl, tar, sha256sum, Python 3, and HTTPS access to `go.dev` and its download destination `dl.google.com`. It requires no application credentials or background services.

## Install or refresh

From the existing checkout:

```sh
bash /workspace/i/.cloud-environment/install.sh
```

The compiler, binary, caches, and test logs are stored outside the checkout in `/workspace/.cloud-setup`. The script is safe to rerun. The repository's source and dependency files remain unchanged. Tests run on a disposable copy because `TestK` regenerates `readme.md`.

For the managed cloud environment, use the contents of `install.sh` as the environment's `install_script`, and the startup instructions below as its `start_skill`. Saving these files in GitHub does not publish a cloud environment snapshot; review, save, and publish the environment in its settings separately.

## Startup instructions

Use the existing isolated checkout at /workspace/i; do not create a Git worktree. The prepared representative development workflow is the dependency-free K library and CLI in /workspace/i/_/i5. Other historical/experimental implementations and WebAssembly are optional and were not validated.
Before Go commands, export PATH=/workspace/.cloud-setup/go/bin:$PATH GOCACHE=/workspace/.cloud-setup/gocache GOMODCACHE=/workspace/.cloud-setup/gomodcache GOTOOLCHAIN=local. The system /usr/bin/go is a different program, not the Go compiler. Verify go version reports go1.27.2 linux/amd64.
No background services or credentials are needed. For interactive use, cd /workspace/i/_/i5 and run /workspace/.cloud-setup/bin/k. EOF exits. For validation, rebuild with go build -mod=readonly -o /workspace/.cloud-setup/bin/k ./k, then send 1+2 and +/!5 on stdin; expect 3 and 10 after the ktye/k banner.
Run tests from a temporary copy: testdir=$(mktemp -d /tmp/i5-tests.XXXXXX); cp -a /workspace/i/_/i5/. "$testdir/"; (cd "$testdir" && go test -mod=readonly -count=1 -v ./...). Remove only that temporary copy afterward. TestK rewrites readme.md, so never execute it in the original checkout during setup. TestK must pass; CLI and mat packages currently have no test files. Preserve tracked source, tests, manifests, lockfiles and existing changes. Check git -C /workspace/i status --short before and after setup.
If setup needs refreshing, execute the saved install_script with bash; a tested local copy is /workspace/.cloud-setup/install.sh. Outputs and caches are retained under /workspace/.cloud-setup. No running process must survive restoration.

## Validation and scope

The setup was validated with Go 1.27.2 on Linux amd64. `TestK` passed; the CLI and matrix packages have no test files. The CLI returned `3` for `1+2` and `10` for `+/!5`. Other historical implementations and WebAssembly builds were not validated.
