package hls

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"syscall"
)

func execute(cmdPath string, args []string) (data []byte, err error) {
	cmd := exec.Command(cmdPath, args...)
	stdout, err1 := cmd.StdoutPipe()
	defer stdout.Close()
	if err1 != nil {
		err = fmt.Errorf("Error opening stdout of command: %v", err)
		return
	}

	log.Debugf("Executing: %v %v", cmdPath, args)
	err2 := cmd.Start()
	if err2 != nil {
		err = fmt.Errorf("Error starting command: %v", err)
		return
	}
	var buffer bytes.Buffer
	_, err3 := io.Copy(&buffer, stdout)
	if err3 != nil {
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		err = fmt.Errorf("Error copying stdout to buffer: %v", err)
		return
	}
	err4 := cmd.Wait()
	if err4 != nil {
		err = fmt.Errorf("Command failed %v", err4)
		return
	}
	data = buffer.Bytes()
	return
}
