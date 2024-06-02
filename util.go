package main

import (
	"encoding/json"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](t T, err error) T {
	check(err)
	return t
}

func command(cond bool, fn func()) {
	if cond {
		fn()
		os.Exit(0)
	}
}

func readJson[T any](file string) T {
	var t T
	bytes := must(os.ReadFile(file))
	check(json.Unmarshal(bytes, &t))
	return t
}

func writeJson[T any](file string, t T) {
	bytes := must(json.Marshal(t))
	check(os.WriteFile(file, bytes, 0600))
}
