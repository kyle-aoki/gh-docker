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
	versionFlag = flag.Bool("version", false, "print version")
	tag         = flag.String("tag", "", "deploy a specific tag")
	job         = flag.String("job", "", "run a docker image as a job")
)

type PetraConfig struct {
	DockerUsername   string
	DockerPassword   string
	DockerRepository string
}

var CFG *PetraConfig = &PetraConfig{}

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	log.Println("petra", version)

	bytes := must(os.ReadFile(homefile(petraConfigFile)))
	check(json.Unmarshal(bytes, &CFG))

	command(*tag != "", func() { dockerDeploy(*tag) })
	command(*job != "", func() { runJob(*job) })

	flag.PrintDefaults()
	os.Exit(1)
}

func command(cond bool, fn func()) {
	if cond {
		fn()
		os.Exit(0)
	}
}
