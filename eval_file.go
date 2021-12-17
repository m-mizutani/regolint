package main

import (
	"context"
	"fmt"

	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/rego"
)

type input struct {
	Files []*RegoFile `json:"files"`
}

func evalWithFile(policyPath string, targets []*RegoFile) error {
	query := []func(*rego.Rego){
		rego.Query(`data`),
		rego.Input(input{Files: targets}),
	}

	policies, err := loadFiles(policyPath)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		query = append(query, rego.ParsedModule(policy.Rego))
	}

	q := rego.New(query...)
	rs, err := q.Eval(context.Background())
	if err != nil {
		return goerr.Wrap(err)
	}

	if len(rs) != 1 {
		panic("result set of Rego must have only one result set")
	}
	if len(rs[0].Expressions) != 1 {
		panic("expression must have only one result set")
	}

	values, ok := rs[0].Expressions[0].Value.(map[string]interface{})
	if !ok {
		logger.With("value", rs[0].Expressions[0].Value).Debug("Value is not map[string]interface{}")
		return nil
	}
	regoLint, ok := values["regolint"].(map[string]interface{})
	if !ok {
		logger.With("regolint", values["regolint"]).Debug("package `regolint` is not found")
		return nil
	}

	fail, ok := regoLint["fail"].([]interface{})
	if !ok {
		logger.With("regolint", regoLint).Debug("'fail' is not found or not array")
		return nil
	}

	if len(fail) > 0 {
		fmt.Println("Failed")
		for _, f := range fail {
			fmt.Println(f)
		}
		return errEvalFailed
	}

	return nil
}
