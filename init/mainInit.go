package init

import "hcc/violin-novnc/lib/config"

// MainInit : Main initialization function
func MainInit() error {
	err := syscheckInit()
	if err != nil {
		return err
	}

	err = LoggerInit()
	if err != nil {
		return err
	}

	config.Parser()

	err = mysqlInit()
	if err != nil {
		return err
	}
	return nil
}
