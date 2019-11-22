package end

import "hcc/violin-novnc/lib/mysql"

func mysqlEnd() {
	if mysql.Db != nil {
		_ = mysql.Db.Close()
	}
}
