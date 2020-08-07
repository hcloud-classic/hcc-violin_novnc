package driver

import (
	"strconv"

	"hcc/violin-novnc/dao"
	"hcc/violin-novnc/lib/logger"
)

const (
	portMin int = 30000
	portMax int = 40000
)

var PM = &PortManager{portMax, []int{}}

type PortManager struct {
	lastPort          int
	availablePortList []int
}

func (pm *PortManager) CheckExistWSPort(token, srvUUID string) ([]int, error) {
	var wspList []int
	selectList, err := dao.SelectWSPortByParam(token, srvUUID)
	if err != nil {
		return nil, err
	}

	for _, wsStr := range selectList {
		wsp, _ := strconv.Atoi(wsStr)
		wspList = append(wspList, wsp)
	}

	return wspList, nil

}

func (pm *PortManager) GetAvailablePort() string {
	var port int

	logger.Logger.Println(pm)
	if len(pm.availablePortList) > 0 {
		port = pm.availablePortList[0]
		if len(pm.availablePortList) == 1 {
			pm.availablePortList = pm.availablePortList[:0]
		} else {
			pm.availablePortList = pm.availablePortList[1:]
		}
		return strconv.Itoa(port)
	}

	if pm.lastPort == portMin {
		return ""
	} else {
		port = pm.lastPort
		pm.lastPort -= 1
		return strconv.Itoa(port)
	}
}

func (pm *PortManager) ReturnPort(port ...int) {
	pm.availablePortList = append(pm.availablePortList, port...)
}
