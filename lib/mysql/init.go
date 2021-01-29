package mysql

import (
	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var cancelHealthCheck func()

func Init() *errors.HccError {
	err, cancel := Prepare()
	cancelHealthCheck = cancel
	return err
}

func End() {
	cancelHealthCheck()
	_ = Db.Close()
}
