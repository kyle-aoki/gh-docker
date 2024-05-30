package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	version               = "v1.0.0"
	petraConfigFile       = "petra-config.json"
	applicationConfigFile = "config.json"
)

var (
	shouldPrintVersion = flag.Bool("version", false, "print version")
	shouldDeployTag    = flag.String("tag", "", "deploy a specific tag")
)

type PetraConfig struct {
	DockerUsername   string
	DockerPassword   string
	DockerRepository string
}

var CFG *PetraConfig = &PetraConfig{}

func main() {
	flag.Parse()
	if *shouldPrintVersion {
		fmt.Println(version)
		os.Exit(0)
	}
	log.Println("petra", version)

	bytes := must(os.ReadFile(homefile(petraConfigFile)))
	check(json.Unmarshal(bytes, &CFG))

	exec(*shouldDeployTag != "", func() { dockerDeploy(*shouldDeployTag) })

	flag.PrintDefaults()
	os.Exit(1)
}

func exec(cond bool, fn func()) {
	if cond {
		fn()
		os.Exit(0)
	}
}
