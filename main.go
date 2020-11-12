package main

import (
	"os"
	"os/exec"

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

type FirewallClient struct {
	argv         []string
	childProcess *exec.Cmd
	// TODO: Define appropriate type instead of interface{}
	autoNets []interface{}
}

func NewFirewallClient() *FirewallClient {
	fc := FirewallClient{}
	ex, err := os.Executable()
	if err != nil {
		ce(err, "executable path")
	}

	elevPrefix := []string{"sudo", "-p", "[local sudo] Password: ", ex}
	path := which(elevPrefix[0])
	if path != "" {
		elevPrefix[0] = path
	}

	argvTries := make([][]string, 0)
	argvTries = append(argvTries, elevPrefix)
	if os.Getuid() == 0 {
		fc.argv = argvTries[len(argvTries)-1]
	}
	for _, argv := range argvTries {
		println(argv)
	}
	fc.argv = argvTries[len(argvTries)-1]
	return &fc
}

func which(cmd string) string {
	// TODO: implement, quit hard coding for *nix.
	return "/usr/bin/" + cmd
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
