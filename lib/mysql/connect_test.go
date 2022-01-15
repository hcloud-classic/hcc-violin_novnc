package mysql

import (
	"hcc/violin-novnc/action/grpc/client"
	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
	"testing"

	"innogrid.com/hcloud-classic/hcc_errors"
)

func Test_DB_Prepare(t *testing.T) {
	err := logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.PiccoloInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	defer func() {
		logger.End()
	}()

	config.Init()

	err = client.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = Init()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		End()
	}()
}
