package main

import (
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
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

	// Copy & paste from https://qiita.com/sky_jokerxx/items/fd79c71143a72cb4efcd
	key, err := ioutil.ReadFile(conf.keyPath)
	ce(err, "private key")

	signer, err := ssh.ParsePrivateKey(key)
	ce(err, "signer")

	auth := []ssh.AuthMethod{ssh.PublicKeys(signer)}
	sshConfig := &ssh.ClientConfig{
		User:            conf.user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// SSH connect.
	client, err := ssh.Dial("tcp", conf.host+":"+strconv.Itoa(conf.port), sshConfig)
	ce(err, "dial")

	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, 25, 80, modes)
	ce(err, "request pty")

	err = session.Shell()
	ce(err, "start shell")

	err = session.Wait()
	ce(err, "return")
}
