package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"gopkg.in/ini.v1"
)

type ConnectionConfig struct {
	user    string
	host    string
	port    int
	keyPath string
}

func ce(err error, msg string) {
	if err != nil {
		println(msg)
		println(err.Error())
		os.Exit(1)
	}
}

type FirewallClient struct{}

func NewFirewallClient() *FirewallClient {
	if os.Getuid() == 0 {
		fmt.Fprint(os.Stdout, "OK")
	} else {
		_, filename, _, _ := runtime.Caller(1)
		out, err := exec.Command("sudo", "-p", "Password:[local sudo] ", "go", "run", filename).Output()
		if err != nil {
			log.Fatalf(err.Error())
		}
		println(string(out))
	}
	return nil
}

func run(conf ConnectionConfig) {
	NewFirewallClient()
}

func main() {
	cfg, err := ini.Load("acrossh.conf")
	ce(err, "config")
	conf := ConnectionConfig{
		// Users are "Krouton" if users don't provide thier username to config.
		user:    cfg.Section("connection").Key("user").MustString("Krouton"),
		host:    cfg.Section("connection").Key("host").MustString("localhost"),
		port:    cfg.Section("connection").Key("port").MustInt(22),
		keyPath: cfg.Section("connection").Key("key_path").MustString("~/.ssh/id_rsa"),
	}
	run(conf)
}
