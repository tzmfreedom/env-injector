package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/okzk/env-injector/envinjector"
)

func main() {
	envinjector.InjectEnviron()

	args := os.Args
	if len(args) <= 1 {
		log.Fatal("missing command")
	}

	path, err := exec.LookPath(args[1])
	if err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(path, args[1:], os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}
