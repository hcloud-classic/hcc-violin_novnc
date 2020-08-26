package mysql

import (
	"hcc/violin-novnc/lib/errors"
)

func Init() *errors.HccError {
	err := Prepare()
	return err
}

func End() {
	_ = Db.Close()
}
