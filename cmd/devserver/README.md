# DevServer (Local Update Emulator)

This **devserver** is a tiny HTTP server for local development/testing of **go-autoupdater**.

It emulates your real update hosting by serving:

- `manifest.json` (at `GET /manifest.json`)
- update artifacts (binaries/exe) (at `GET /files/<filename>`)

So you can test update flows **locally** without uploading to a real server.

> This is a **development tool**. Do not use it as a production update server.

---

## What you need

- Go **1.22+**
- The repo checked out locally
- One or more built artifacts (example: `agent.exe`, `agent_linux_arm64`, etc.)

---

## Folder layout (recommended)

Put your artifacts and (optionally) a static manifest under a single folder.

Recommended structure:

```text
<repo-root>/
  cmd/
    devserver/
      main.go
      README.md
  testdata/
    updates/
      manifest.json               (optional if you use --gen)
      agent_1.0.12_windows_amd64.exe
      agent_1.0.12_linux_amd64
      agent_1.0.12_linux_arm64
      agent_1.0.12_darwin_arm64
```

Why `testdata/`?
- It’s a common Go convention for sample/test resources.
- It keeps dev-only artifacts separate from library code.

---

## Two ways to run devserver

### Mode A — Static manifest (you write `manifest.json` yourself)

Use this when you want full control over the manifest content.

1) Create `testdata/updates/manifest.json`
2) Put your artifacts in the same folder
3) Run:

```bash
go run ./cmd/devserver --root ./testdata/updates
```

Then open:
- Manifest: `http://127.0.0.1:8089/manifest.json`
- Files: `http://127.0.0.1:8089/files/<filename>`


### Mode B — Auto-generate manifest (fastest)

Use this when you want devserver to generate a manifest automatically.

**Requirement:** artifact filenames must contain `_<os>_<arch>` at the end.

Examples that will be detected:

- `agent_1.0.12_windows_amd64.exe`
- `agent_1.0.12_linux_amd64`
- `agent_1.0.12_linux_arm64`
- `agent_1.0.12_darwin_arm64`

Run:

```bash
go run ./cmd/devserver --root ./testdata/updates --gen --version 1.0.12
```

devserver will:
- scan files in `--root`
- compute SHA256 for each detected artifact
- output a generated manifest at `GET /manifest.json`

---

## Flags

| Flag | Default | Description |
|---|---:|---|
| `--addr` | `127.0.0.1:8089` | Address to listen on |
| `--root` | `./testdata/updates` | Directory containing artifacts (+ optional manifest) |
| `--manifest` | `manifest.json` | Manifest filename under `--root` (Mode A) |
| `--gen` | `false` | Enable auto-generate manifest (Mode B) |
| `--product` | `agent` | Product name used in generated manifest |
| `--channel` | `stable` | Channel used in generated manifest |
| `--version` | `0.0.0-dev` | Version used in generated manifest |
| `--notes` | `local dev build` | Notes used in generated manifest |


### Example: change port

```bash
go run ./cmd/devserver --addr 127.0.0.1:9090 --root ./testdata/updates --gen --version 1.0.12
```

---

## Build devserver binary

You can also build it as a standalone binary.

### macOS / Linux

```bash
go build -o devserver ./cmd/devserver
./devserver --root ./testdata/updates --gen --version 1.0.12
```

### Windows (PowerShell)

```powershell
go build -o devserver.exe .\cmd\devserver
.\devserver.exe --root .\testdata\updates --gen --version 1.0.12
```

---

## Using devserver with `updaterctl`

Once devserver is running, point your updater to the local manifest URL.

Example (standalone mode):

```bash
./updaterctl \
  --manifest "http://127.0.0.1:8089/manifest.json" \
  --dir "./" \
  --exe "agent" \
  --current "1.0.11"
```

### Windows service update test (NSSM/SC)

- Place `updater-helper.exe` next to your `agent.exe` (same install dir)
- Run `updaterctl.exe` with:

```powershell
updaterctl.exe `
  --manifest "http://127.0.0.1:8089/manifest.json" `
  --dir "C:\\fck\\agent" `
  --exe "agent.exe" `
  --current "1.0.11" `
  --service "FCKAgent" `
  --nssm "C:\\tools\\nssm.exe"
```

> If you want to force Windows built-in `sc.exe` controller, use: `--nssm SC`

---

## What you may need to change

### 1) Your artifact filenames (for `--gen` mode)

If you want auto-generation, make sure your file name ends with:

- `_<os>_<arch>` (plus `.exe` for Windows)

Valid OS values (Go runtime):
- `windows`, `linux`, `darwin`

Common arch values:
- `amd64`, `arm64`, `arm`, `386`

Examples:
- `agent_1.0.12_windows_amd64.exe`
- `agent_1.0.12_linux_arm64`
- `agent_1.0.12_darwin_arm64`

If you don’t want to rename files, use **Mode A** (static manifest) and write URLs manually.

### 2) The `--root` directory

If your artifacts live somewhere else, just point devserver to it:

```bash
go run ./cmd/devserver --root /path/to/your/artifacts --gen --version 1.0.12
```

### 3) Using `localhost` vs `127.0.0.1`

If your updater runs in a VM/container/another device, `127.0.0.1` will refer to that device itself.

In that case:
- bind devserver on `0.0.0.0:<port>`
- use your host IP address

Example:

```bash
go run ./cmd/devserver --addr 0.0.0.0:8089 --root ./testdata/updates --gen --version 1.0.12
```

Then access from another machine:
- `http://<YOUR_HOST_IP>:8089/manifest.json`

---

## Troubleshooting

### “No artifacts detected” (when using `--gen`)

This means devserver can’t infer `os/arch` from filenames.

Fix:
- rename artifacts to include `_<os>_<arch>` near the end, e.g. `agent_1.0.12_linux_amd64`
- or use Mode A (static `manifest.json`)

### 404 for a file

- Ensure the artifact file exists inside the `--root` directory
- Ensure the URL in the manifest points to `/files/<filename>`

### Updater says “sha256 mismatch”

- You changed the artifact file after manifest was created
- Regenerate manifest (`--gen`) or update the SHA256 in your static `manifest.json`

---

## Notes

- devserver sets `Cache-Control: no-store` to avoid caching during development.
- This server is intentionally simple and should not be used for production.

---

## License

This devserver is part of the main repository and follows the repository license.
