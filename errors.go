package main

import "github.com/m-mizutani/goerr"

var (
	errInvalidConfig = goerr.New("invalid configuration")

	errEvalFailed = goerr.New("got evaluation failure")
)
