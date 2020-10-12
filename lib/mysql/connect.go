package mysql

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql

	"hcc/violin-novnc/lib/config"
	"hcc/violin-novnc/lib/errors"
)

// Db : Pointer of mysql connection
var Db *sql.DB

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() *errors.HccError {
	var err error
	Db, err = sql.Open("mysql",
		config.Mysql.ID+":"+config.Mysql.Password+"@tcp("+
			config.Mysql.Address+":"+strconv.Itoa(int(config.Mysql.Port))+")/"+
			config.Mysql.Database+"?parseTime=true")
	if err != nil {
		return errors.NewHccError(errors.ViolinNoVNCInternalInitFail, "mysql open")
	}

	err = Db.Ping()
	if err != nil {
		return errors.NewHccError(errors.ViolinNoVNCInternalInitFail, "mysql ping")
	}

	logger.Logger.Println("db is connected")

	return nil
}
