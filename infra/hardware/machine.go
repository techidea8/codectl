package hardware

import (
	"fmt"

	"github.com/super-l/machine-code/machine"
	"github.com/super-l/machine-code/machine/types"
	"github.com/techidea8/codectl/infra/cryptor"
)

// 硬件相关信息
func machineData() types.Information {
	return machine.GetMachineData()
}

type MachineInfo struct {
	types.Information
	MachineId string
}

// 硬件ID
func Machine() MachineInfo {
	data := machineData()
	id := cryptor.Md5String(fmt.Sprintf("bordsn=%s,cpusn=%s", data.BoardSerialNumber, data.CpuSerialNumber))
	machineInfo := MachineInfo{
		data, id,
	}
	return machineInfo
}
