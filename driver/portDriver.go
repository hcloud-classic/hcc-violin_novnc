package driver

import (
	"strconv"

	"hcc/violin-novnc/lib/logger"
)

const (
	portMin int = 30000
	portMax int = 40000
)

var PD = &PortDriver{portMax, []int{}}

type PortDriver struct {
	lastPort          int
	availablePortList []int // TODO: change to rbtree
}

func (pd *PortDriver) GetAvailablePort() string {
	var port int

	if len(pd.availablePortList) > 0 {
		port = pd.availablePortList[0]
		if len(pd.availablePortList) == 1 {
			pd.availablePortList = pd.availablePortList[:0]
		} else {
			pd.availablePortList = pd.availablePortList[1:]
		}

		return strconv.Itoa(port)
	}

	if pd.lastPort < portMin {
		return ""
	} else {
		port = pd.lastPort
		pd.lastPort -= 1

		return strconv.Itoa(port)
	}
}

func (pd *PortDriver) SetLastPort(port string) {
	p, _ := strconv.Atoi(port)

	if p > pd.lastPort {
		return
	}

	for ; pd.lastPort >= p; pd.lastPort-- {
		pd.availablePortList = append(pd.availablePortList, p)
	}

}

func (pd *PortDriver) ReturnPort(port string) {
	p, err := strconv.Atoi(port)
	if err != nil {
		logger.Logger.Println("Wrong port string")
		return
	}

	if p < portMin || p > portMax {
		logger.Logger.Println("Return port not in range ", portMin, "-", portMax)
	}

	pd.availablePortList = append(pd.availablePortList, p)
}
