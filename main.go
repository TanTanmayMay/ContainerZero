package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		panic("Wrong Command")
	}

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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// syscall.Sethostname([]byte("container")); -> Setting prompt name here will not work because
	// Namespace is not set until we hit the cmd.Run() command
	checkErr(cmd.Run())
}

func child() {
	fmt.Printf("Running (In Child) %v\n", os.Args[2:])

	cg()

	// Unshare mount namespace
	checkErr(syscall.Unshare(syscall.CLONE_NEWNS))

	// Remount / as private to make sure changes are not propagated to the host
	checkErr(syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, ""))

	checkErr(syscall.Sethostname([]byte("container_zero")))
	checkErr(syscall.Chroot("./manjaro_fs"))
	// Change directory after chroot
	checkErr(os.Chdir("/"))

	fmt.Println("Mounting in the new namespace...")
	// Mount proc and tmpfs in the new mount namespace
	checkErr(syscall.Mount("proc", "/proc", "proc", 0, ""))
	checkErr(syscall.Mount("something2", "/mytemp", "tmpfs", 0, ""))

	// cmd is typically a variable representing an instance of exec.Cmd, which is used to execute external commands.
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())

	fmt.Println("Unmounting in the new namespace...")

	checkErr(syscall.Unmount("/proc", 0))
	checkErr(syscall.Unmount("/mytemp", 0))
}

func cg() {
	cgroups := "/sys/fs/cgroup"
	mem := filepath.Join(cgroups, "memory")
	os.Mkdir(filepath.Join(mem, "tantanmaymay"), 0755)
	checkErr(os.WriteFile(filepath.Join(mem, "tantanmaymay/memory.limit_in_bytes"), []byte("999424"), 0700))
	checkErr(os.WriteFile(filepath.Join(mem, "tantanmaymay/memory.memsw.limit_in_bytes"), []byte("999424"), 0700))
	checkErr(os.WriteFile(filepath.Join(mem, "tantanmaymay/notify_on_release"), []byte("1"), 0700))

	pid := strconv.Itoa(os.Getpid())
	checkErr(os.WriteFile(filepath.Join(mem, "tantanmaymay/cgroup.proc"), []byte(pid), 0700))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
