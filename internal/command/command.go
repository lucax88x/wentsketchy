package command

import (
	"fmt"
	"os/exec"
)

func Run(name string, arg ...string) (string, error) {
	fmt.Println(fmt.Sprintf("%s %d %v", name, len(arg), arg))
	cmd := exec.Command(name, arg...)

	out, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("could not run command %w", err)
	}

	return string(out), nil
}
