package qos

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func ExecCmd(timeoutsec int, cmd string, args ...string) (stdoutstr string, stderrstr string, dotimeout bool, err error) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var execcmd = exec.Command(cmd, args...)
	execcmd.Stdout = &stdout
	execcmd.Stderr = &stderr

	err = execcmd.Start()
	if err != nil {
		return "", "", false, err
	}

	go func() {
		var chtimeout = time.After(time.Duration(timeoutsec) * time.Second)
		select {
		case <-chtimeout:
			execcmd.Process.Kill()
		}
	}()

	err = execcmd.Wait()
	if err == nil {
		dotimeout = false
	} else {
		err = nil
		dotimeout = true
	}
	stdoutstr = stdout.String()
	stderrstr = stderr.String()
	execcmd.Process.Release()
	return stdoutstr, stderrstr, dotimeout, err
}

// Fork and exec this same image without dropping the net.Listener.
func Fork() (*os.Process, error) {
	execpath, err := lookPath()
	if nil != err {
		return nil, err
	}
	wd, err := os.Getwd()
	if nil != err {
		return nil, err
	}
	os.Environ()
	files := make([]*os.File, 3)
	files[syscall.Stdin] = os.Stdin
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr
	return os.StartProcess(execpath, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
}

func Kill(pid int) error {
	var process, err = os.FindProcess(pid)
	if err == nil {
		err = process.Kill()
	}
	return err
}

func lookPath() (argv0 string, err error) {
	argv0, err = exec.LookPath(os.Args[0])
	if nil != err {
		return
	}
	if _, err = os.Stat(argv0); nil != err {
		return
	}
	return
}
