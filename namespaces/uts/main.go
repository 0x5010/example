package main

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		log.Errorln("missing commands")
		return
	}
	switch os.Args[1] {
	case "run":
		run()
	default:
		log.Errorln("wrong command")
		return
	}
}

func run() {
	log.Infof("Running %v", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	check(cmd.Run())
}

func check(err error) {
	if err != nil {
		log.Errorln(err)
	}
}
