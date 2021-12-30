package regolint

test_unification {
    not _test_unification
}

_test_unification {
    tests := [
        {
            "title": "fail if unification is used",
            "input": data.testdata.unification.used,
            "exp": {"Do not use unification `=` in policies/policy.rego"},
        },
        {
            "title": "pass if unification is not used",
            "input": data.testdata.unification.not_used,
            "exp": set(),
        },
    ]

    t := tests[_]
	resp := fail with input as t.input
	resp != t.exp
	print(sprintf("failed '%s'. '%v' is expected, but got '%v'", [t.title, t.exp, resp]))
}
