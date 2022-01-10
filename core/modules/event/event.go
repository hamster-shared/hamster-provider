package event

import "strconv"

type VmRequest struct {
	Tag         OperationTag
	Cpu         uint64
	Mem         uint64
	Disk        uint64
	AccessPort  uint64
	Type        string
	Image       string
	System      string
	PublicKey   string
	OrderNo     uint64
	AgreementNo uint64
}

func (req *VmRequest) getName() string {
	return "order_" + strconv.Itoa(int(req.OrderNo))
}

type OperationTag int

const OPCreatedVm OperationTag = 1
const OPDestroyVm OperationTag = 2
const OPRenewVM OperationTag = 3
const OPRecoverVM OperationTag = 4
