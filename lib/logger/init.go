package logger

import "hcc/violin-novnc/lib/errors"

// Init : Prepare logger
func Init() *errors.HccError {
	if !Prepare() {
		errors.SetErrLogger(Logger)
		return errors.NewHccError(errors.ClarinetInternalInitFail, "logger")
	}

	errors.SetErrLogger(Logger)

	return nil
}

// End : Close logger
func End() {
	_ = FpLog.Close()
}
