# regolint

Lint Rego file with policy written by Rego

## Usage

### Command Line

```bash
$ regolint -p ./lint ./policy
./policy/auth.rego: package path and directory path is not matched
```

### GitHub Actions

An image on GitHub Container Registry is available: `ghcr.io/m-mizutani/regolint:latest`.

You can use the container image with such following GitHub Actions workflow.

```yaml
name: Lint

on: [push]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout upstream repo
        uses: actions/checkout@v2
        with:
          ref: ${{ github.head_ref }}
      - uses: docker://ghcr.io/m-mizutani/regolint:latest
        with:
          args: "-p ./lint ./policy"
```

### Options

- `-p, --policy`: lint policy file/dir. If no policy file, output only parsed rego files
- `-o, --output`: specify output file. `-` means stdout

## Rule guide

### Parameters

These requirements are for local enforce policy file specified by `-p`.

- package name: must be `regolint`
- input: `files`
    - `path` (array of string): File path splitted by delimiter
    - `rego` (ast): Rego rule structured as `ast.Module` in `github.com/open-policy-agent/opa`
- output: `fail[msg]` should be set failure message(s).

### Example

```rego
package regolint

# Check file path
fail[msg] {
    file := input.files[_]
    count(file.path) <= 1
    msg := sprintf("%s: .rego file at top level is not allowed", [concat("/", file.path)])
}

# Check matching with directory path and package path
fail[msg] {
    file := input.files[_]

    count(file.path) > 1

    dirpath := array.slice(file.path, 0, count(file.path) - 2)
    pkgpath := array.slice(file.rego.package.path, 1, count(file.rego.package.path) - 1)

    some i
    count(dirpath[i] != pkgpath[i]) > 0
    msg := sprintf("%s: package path and directory path is not matched", [concat("/", file.path)])
}
```

## License

MIT License
