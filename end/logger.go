package end

import "hcc/vnc/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
