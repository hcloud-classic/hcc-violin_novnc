package init

import (
	"errors"
	"hcc/violin-novnc/lib/logger"
)

// LoggerInit : Init logger
func LoggerInit() error {
	if !logger.Prepare() {
		return errors.New("error occurred while preparing logger")
	}

	return nil
}
