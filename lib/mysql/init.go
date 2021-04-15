package mysql

import (
	errors "innogrid.com/hcloud-classic/hcc_errors"
)

func Init() *errors.HccError {
	err := Prepare()
	return err
}

func End() {
	_ = Db.Close()
}
