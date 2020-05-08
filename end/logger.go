package end

import "hcc/violin-novnc/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
