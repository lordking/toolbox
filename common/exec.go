package common

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

//ExecCommand 执行系统命令
func ExecCommand(commandName string, arg ...string) (bool, error) {
	cmd := exec.Command(commandName, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, NewError(ErrCodeInternal, err.Error())
	}
	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
	}

	cmd.Wait()
	return true, nil
}
