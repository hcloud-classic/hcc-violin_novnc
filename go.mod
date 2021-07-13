module hcc/violin-novnc

go 1.13

require (
	github.com/Terry-Mao/goconf v0.0.0-20161115082538-13cb73d70c44
	github.com/go-sql-driver/mysql v1.5.0
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	google.golang.org/grpc v1.39.0
	innogrid.com/hcloud-classic/hcc_errors v1.1.2
	innogrid.com/hcloud-classic/pb v0.0.0
)

replace (
	innogrid.com/hcloud-classic/hcc_errors => ../hcc_errors
	innogrid.com/hcloud-classic/pb => ../pb
)
