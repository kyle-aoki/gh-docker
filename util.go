package main

import (
	"os"
	"path/filepath"
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

func homefile(filename string) string {
	home := must(os.UserHomeDir())
	return filepath.Join(home, filename)
}
