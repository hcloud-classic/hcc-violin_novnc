package mysql

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
	errors "innogrid.com/hcloud-classic/hcc_errors"

	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/logger"
)

// Db : Pointer of mysql connection
var (
	ctx context.Context
	Db  *sql.DB
)

// Type alias
type Rows = sql.Rows
type Result = sql.Result

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() (*errors.HccError, func()) {
	var err error
	Db, err = sql.Open("mysql",
		config.Mysql.ID+":"+config.Mysql.Password+"@tcp("+
			config.Mysql.Address+":"+strconv.Itoa(int(config.Mysql.Port))+")/"+
			config.Mysql.Database+"?parseTime=true")
	if err != nil {
		return errors.NewHccError(errors.ViolinNoVNCInternalInitFail, "mysql open"), nil
	}

	timeTicker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	cancel := func() { done <- true }
	go func() {
		for true {
			select {
			case <-done:
				return
			case <-timeTicker.C:
				err = Db.Ping()
				if err != nil {
					logger.Logger.Println(
						errors.NewHccError(errors.ViolinNoVNCInternalConnectionFail,
							"mysql connection lost, retry...").Error())
				}
				break
			}
		}
		if err != nil {

		}
	}()

	return nil, cancel
}
