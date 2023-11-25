package filemenager

import (
	"os/exec"
	"read_files/util/constants"
)

func MoveFile(destination string) error {
	cmd := exec.Command(constants.Mv, destination, constants.TempDirPath)
	return cmd.Run()
}
