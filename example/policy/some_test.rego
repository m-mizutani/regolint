package regolint

test_some_declared {
    not _test_some_declared
}

_test_some_declared {
    tests := [
        {
            "title": "fail if variable is not declared",
            "input": data.testdata["some"].undeclared,
            "exp": {"'c' is not declared. Assign or use some keyword explicitly"},
        },
        {
            "title": "pass if unification is declared explicitly",
            "input": data.testdata["some"].declared,
            "exp": set(),
        },
    ]

    t := tests[_]
	resp := fail with input as t.input
	resp != t.exp
	print(sprintf("failed '%s'. '%v' is expected, but got '%v'", [t.title, t.exp, resp]))
}

