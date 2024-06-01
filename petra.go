package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	version               = "v1.0.0"
	petraConfigFile       = "/petra/petra-config.json"
	applicationConfigFile = "/petra/app-config.json"
	petraLoopDelay        = time.Second * 5
)

var (
	versionFlag = flag.Bool("version", false, "print version")
	job         = flag.String("job", "", "run a docker image as a job")
	background  = flag.Bool("background", false, "run petra background process")
	set         = flag.String("set", "", "set target docker tag")
)

type PetraConfig struct {
	DockerUsername   string
	DockerPassword   string
	DockerRepository string
	TargetDockerTag  string
	CurrentDockerTag string
}

var CFG *PetraConfig = &PetraConfig{}

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	log.Println("petra", version)

	bytes := must(os.ReadFile(petraConfigFile))
	check(json.Unmarshal(bytes, &CFG))

	command(*job != "", func() { runJob(*job) })
	command(*set != "", func() { setTargetDockerTag(*set) })
	command(*background, func() { petraWorker() })

	flag.PrintDefaults()
	os.Exit(1)
}

func command(cond bool, fn func()) {
	if cond {
		fn()
		os.Exit(0)
	}
}

func setTargetDockerTag(tag string) {
	log.Println("setting target docker tag:", tag)
	cfg := readJson[PetraConfig](petraConfigFile)
	cfg.TargetDockerTag = tag
	writeJson(petraConfigFile, cfg)
}

func petraWorker() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go infiniteLoop()
	wg.Wait()
}

func infiniteLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("CRASH:", r)
			go infiniteLoop()
		}
	}()
	loop()
}

func loop() {
	log.Println("started petra background loop")
	for {
		petraCFG := readJson[PetraConfig](petraConfigFile)
		if petraCFG.CurrentDockerTag != petraCFG.TargetDockerTag {
			log.Println("petra config tags do not match:", fmt.Sprintf("%s != %s", petraCFG.CurrentDockerTag, petraCFG.TargetDockerTag))
			dockerDeploy(petraCFG.TargetDockerTag)
			petraCFG.CurrentDockerTag = petraCFG.TargetDockerTag
			writeJson(petraConfigFile, petraCFG)
		} else {
			log.Println("petra config tags match:", fmt.Sprintf("%s == %s", petraCFG.CurrentDockerTag, petraCFG.TargetDockerTag))
		}
		time.Sleep(petraLoopDelay)
	}
}
