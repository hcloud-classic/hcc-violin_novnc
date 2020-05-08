package init

import "hcc/violin-novnc/lib/mysql"

func mysqlInit() error {
	err := mysql.Prepare()
	if err != nil {
		return err
	}

	return nil
}
