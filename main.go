package main

import (
	"os"

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
	// example.com is expected.
	println(conf.host)
}
