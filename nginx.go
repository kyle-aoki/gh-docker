package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	nginxconf = "/etc/nginx/nginx.conf"
	proxypass = "proxy_pass http://localhost:{{ PORT }};"
)

func oppositePort(port string) string {
	if port == "8080" {
		return "8081"
	}
	if port == "8081" {
		return "8080"
	}
	panic("unknown port " + port)
}

func switchProxyPort(port string) {
	log.Println("switching nginx proxy port from", oppositePort(port), "to", port)
	currNginxConf := string(must(os.ReadFile(nginxconf)))
	currentProxyPass := strings.ReplaceAll(proxypass, "{{ PORT }}", oppositePort(port))
	if !strings.Contains(currNginxConf, currentProxyPass) {
		panic("unknown ngnix.conf state:\n" + nginxconf)
	}
	newProxyPass := strings.ReplaceAll(proxypass, "{{ PORT }}", port)
	newNginxConf := strings.ReplaceAll(currNginxConf, currentProxyPass, newProxyPass)
	check(os.WriteFile(nginxconf, []byte(newNginxConf), 0600))
	log.Println("updated nginx.conf with new directive:", newProxyPass)
	cmd := exec.Command("nginx", "-s", "reload")
	check(cmd.Run())
	log.Println("reloaded nginx config")
}