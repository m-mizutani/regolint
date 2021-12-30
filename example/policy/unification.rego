package regolint

# Do not allow use '=' (unification)
fail[msg] {
    file := input.files[_]
    body := file.rego.rules[_].body[_]
    term := body.terms[_]
    term.type == "ref"
    term.value[_].value == "eq"

    msg := sprintf("Do not use unification `=` in %s", [concat("/", file.path)])
}
