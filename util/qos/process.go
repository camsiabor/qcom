package qos

import (
	"bytes"
	"os/exec"
	"time"
)


func ExecCmd(timeoutsec int, cmd string, args ... string) (stdoutstr string,stderrstr string, dotimeout bool, err error) {

	var stdout bytes.Buffer;
	var stderr bytes.Buffer;
	var execcmd = exec.Command(cmd, args...);
	execcmd.Stdout = &stdout;
	execcmd.Stderr = &stderr;

	err = execcmd.Start()
	if (err != nil) {
		return "", "", false, err;
	}

	go func() {
		var chtimeout = time.After(time.Duration(timeoutsec) * time.Second);
		select {
		case <-chtimeout:
			execcmd.Process.Kill()
		}
	}();

	err = execcmd.Wait();
	if (err == nil) {
		dotimeout = false;
	} else {
		err = nil;
		dotimeout = true;
	}

	stdoutstr = stdout.String();
	stderrstr = stderr.String();
	execcmd.Process.Release();
	return stdoutstr, stderrstr, dotimeout, err;
}
