package filemenager

import (
	"os/exec"
	"read_files/config"
)

func MoveFile(destination string) error {
	cmd := exec.Command(constants.Mv, destination, constants.TargetDir)
	return cmd.Run()
}
