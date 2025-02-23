package wrapper

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

type JavaExec interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
	Stop() error
	Kill() error
}

type defaultJavaExec struct {
	cmd  *exec.Cmd
	done chan error
}

func (j *defaultJavaExec) Stdout() io.ReadCloser {
	r, _ := j.cmd.StdoutPipe()
	return r
}

func (j *defaultJavaExec) Stdin() io.WriteCloser {
	w, _ := j.cmd.StdinPipe()
	return w
}

func (j *defaultJavaExec) Start() error {
	if err := j.cmd.Start(); err != nil {
		return err
	}

	go func() {
		j.done <- j.cmd.Wait()
	}()

	return nil
}

func (j *defaultJavaExec) Stop() error {
	if j.cmd.Process == nil {
		return fmt.Errorf("server not running")
	}

	select {
	case <-j.done:
		return nil
	case <-time.After(10 * time.Second):
		if err := j.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
		return j.Kill()
	}
}

func (j *defaultJavaExec) Kill() error {
	return j.cmd.Process.Kill()
}

func javaExecCmd(serverPath string, initialHeapSize, maxHeapSize int) *defaultJavaExec {
	initialHeapFlag := fmt.Sprintf("-Xms%dM", initialHeapSize)
	maxHeapFlag := fmt.Sprintf("-Xmx%dM", maxHeapSize)
	cmd := exec.Command("java", initialHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")
	return &defaultJavaExec{cmd: cmd, done: make(chan error, 1)}
}
