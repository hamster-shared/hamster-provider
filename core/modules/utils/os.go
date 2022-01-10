package utils

import "github.com/shirou/gopsutil/cpu"

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
