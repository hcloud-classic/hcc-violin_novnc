package pid

import (
	"hcc/violin-novnc/lib/fileutil"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

var ViolinNovncPIDFileLocation = "/var/run"
var ViolinNovncPIDFile = "/var/run/violin-novnc.pid"

// IsViolinNovncRunning : Check if ViolinNovnc is running
func IsViolinNovncRunning() (running bool, pid int, err error) {
	if _, err := os.Stat(ViolinNovncPIDFile); os.IsNotExist(err) {
		return false, 0, nil
	}

	pidStr, err := ioutil.ReadFile(ViolinNovncPIDFile)
	if err != nil {
		return false, 0, err
	}

	ViolinNovncPID, _ := strconv.Atoi(string(pidStr))

	proc, err := os.FindProcess(ViolinNovncPID)
	if err != nil {
		return false, 0, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, ViolinNovncPID, nil
	}

	return false, 0, nil
}

// WriteViolinNovncPID : Write ViolinNovnc PID to "/var/run/ViolinNovnc.pid"
func WriteViolinNovncPID() error {
	pid := os.Getpid()

	err := fileutil.CreateDirIfNotExist(ViolinNovncPIDFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(ViolinNovncPIDFile, strconv.Itoa(pid))
	if err != nil {
		return err
	}

	return nil
}

// DeleteViolinNovncPID : Delete the ViolinNovnc PID file
func DeleteViolinNovncPID() error {
	err := fileutil.DeleteFile(ViolinNovncPIDFile)
	if err != nil {
		return err
	}

	return nil
}
