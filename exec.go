package wrapper

import (
	"fmt"
	"io"
	"os/exec"
)

type JavaExec interface {
	Stdout() io.ReadCloser
	Stdin() io.WriteCloser
	Start() error
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

func (j *defaultJavaExec) Kill() error {
	return j.cmd.Process.Kill()
}

func javaExecCmd(serverPath string, initialHeapSize string, maxHeapSize string) *defaultJavaExec {
	initialHeapFlag := fmt.Sprintf("-Xms%s", initialHeapSize)
	maxHeapFlag := fmt.Sprintf("-Xmx%s", maxHeapSize)
	cmd := exec.Command("java", initialHeapFlag, maxHeapFlag, "-jar", serverPath, "nogui")
	return &defaultJavaExec{cmd: cmd, done: make(chan error, 1)}
}
