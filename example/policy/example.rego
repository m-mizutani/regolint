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

    dir_path := array.slice(file.path, 0, count(file.path) - 2)
    pkg_path := array.slice(file.rego["package"].path, 1, count(file.rego["package"].path) - 1)

    some i
    count({ i | dir_path[i] != pkg_path[i] }) > 0
    msg := sprintf("%s: package path and directory path is not matched", [concat("/", file.path)])
}
