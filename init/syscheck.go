package init

import "hcc/violin-scheduler/lib/syscheck"

func syscheckInit() error {
	err := syscheck.CheckRoot()
	if err != nil {
		return err
	}

	return nil
}
