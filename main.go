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
	case "child":
		child()
	default:
		panic("Wrong Command")
	}
}

func run() {
	fmt.Printf("Running (In Parent) %v\n", os.Args[2:])

	// cmd is typically a variable representing an instance of exec.Cmd, which is used to execute external commands.
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// syscall.Sethostname([]byte("container")); -> Setting prompt name here will not work because
	// Namespace is not set until we hit the cmd.Run() command
	checkErr(cmd.Run())
}

func child() {
	fmt.Printf("Running (In Child) %v\n", os.Args[2:])

	// cmd is typically a variable representing an instance of exec.Cmd, which is used to execute external commands.
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	checkErr(syscall.Sethostname([]byte("container_zero")))
	checkErr(syscall.Chroot("./manjaro_fs"))
	// Change directory after chroot
	checkErr(os.Chdir("./manjaro_fs"))
	checkErr(cmd.Run())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
