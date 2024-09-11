package command

import (
	"fmt"
	"os/exec"
)

func Run(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	fmt.Println(fmt.Sprintf("%s %v", cmd.Path, cmd.Args))

	out, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("could not run command %w", err)
	}

	return string(out), nil
}
