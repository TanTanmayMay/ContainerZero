package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("Wrong Command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	// cmd is typically a variable representing an instance of exec.Cmd, which is used to execute external commands.
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	
	checkErr(cmd.Run())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}