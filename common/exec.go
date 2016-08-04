package common

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func ExecCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}

	cmd.Wait()
	return true
}
