package logger

import errors "github.com/hcloud-classic/hcc_errors"

// Init : Prepare logger
func Init() *errors.HccError {
	if !Prepare() {
		errors.SetErrLogger(Logger)
		return errors.NewHccError(errors.ViolinNoVNCInternalInitFail, "logger")
	}

	errors.SetErrLogger(Logger)

	return nil
}

// End : Close logger
func End() {
	_ = FpLog.Close()
}
