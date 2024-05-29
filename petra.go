package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PetraConfig struct {
	DockerUsername   string
	DockerPassword   string
	DockerRepository string
	Key              string
}

var CFG *PetraConfig = &PetraConfig{}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		panic("petra <path-to-config-file.json>")
	}
	bytes := must(os.ReadFile(args[0]))
	check(json.Unmarshal(bytes, &CFG))

	http.HandleFunc("/deploy", httpDeploy)
	http.HandleFunc("/list", httpList)

	am := &AuthMiddleware{handler: http.DefaultServeMux}

	log.Println("listening on port 10000")
	check(http.ListenAndServe("0.0.0.0:10000", am))
}

func handleHttpPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("%s", r)))
	}
}

func httpDeploy(w http.ResponseWriter, r *http.Request) {
	type DeployInput struct {
		Tag string
	}
	defer handleHttpPanic(w)
	bytes := must(io.ReadAll(r.Body))
	input := &DeployInput{}
	check(json.Unmarshal(bytes, &input))
	containers := deploy(input.Tag)
	w.Write(must(json.Marshal(containers)))
}

func httpList(w http.ResponseWriter, r *http.Request) {
	defer handleHttpPanic(w)
	w.Write(must(json.Marshal(listContainers())))
}
