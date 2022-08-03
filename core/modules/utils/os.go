package utils

import (
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"os"
)

//  GetCpuModel get cpu model
func GetCpuModel() string {
	cpuInfo, err := cpu.Info()

	if err != nil {
		return "unknow CPU"
	}

	if len(cpuInfo) > 0 {
		return cpuInfo[0].ModelName
	} else {
		return "unknow CPU"
	}
}

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
