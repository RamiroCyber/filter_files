package filemenager

import (
	"fmt"
	"os"
	"read_files/util/constants"
)

func CreateDirIfNotExist() error {
	if _, err := os.Stat(constants.TempDirPath); os.IsNotExist(err) {
		return os.Mkdir(constants.TempDirPath, 0755)
	}
	return nil
}

func GenerateTempFilePath(filePath string) string {
	return fmt.Sprintf("./%s/ %s", constants.TempDirPath, filePath)
}

func RemoveFiles(tempFilePath string) error {
	err := os.Remove(tempFilePath)
	if err != nil {
		fmt.Println("Erro ao excluir o arquivo temporário:", err)
	}
	return err
}
