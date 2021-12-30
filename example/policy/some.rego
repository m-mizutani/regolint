package regolint

# Check if variables are declared
fail[msg] {
    file := input.files[_]
    bodies := file.rego.rules[_].body

    tv := walk(bodies[_].terms)
    variables := { tv[_v].value | is_variable(tv[_v]) }

    v := variables[_]

    # Check if assigned at local
    not is_assigned_local(v, bodies)
    not is_assigned_symbol(v, bodies)
    not is_assigned_head(v, file.rego.rules)

    msg := sprintf("'%v' is not declared. Assign or use some keyword explicitly", [v])
}

is_variable(term) {
    term.type == "var"
    term.value != "input"
    term.value != "equal"
    term.value != "assign"
    term.value != "eq"
    is_string(term.value)
}

is_assigned_local(v, bodies) {
    local_terms := bodies[_].terms
    local_terms[0].type == "ref"
    local_terms[0].value[0].type == "var"
    local_terms[0].value[0].value == "assign"
    local_terms[1].type == "var"
    v == local_terms[1].value
}

is_assigned_symbol(v, bodies) {
    symbol_terms := bodies[_].terms
    symbol_terms.symbols[_s].type == "var"
    v == symbol_terms.symbols[_s].value
}

is_assigned_head(v, rules) {
    rules[_r].head.assign
    v == rules[_r].head.name
}
