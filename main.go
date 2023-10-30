package main

import (
	"flag"
	"fmt"

	"github.com/Archisman-Mridha/csi-plugin/pkg/driver"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "", "Endpoint the gRPC server will be listening at")
		region   = flag.String("region", "", "DigitalOcean region where volumes will be provisioned")
		token    = flag.String("token", "", "DigitalOcean api token")
	)
	flag.Parse()

	csiDriver := driver.NewCsiDriver(driver.NewCsiDriverArgs{
		Region: *region,
		Token:  *token,

		Endpoint: *endpoint,
	})
	if err := csiDriver.Run(); err != nil {
		fmt.Printf("Error running CSI driver : %v", err)
	}
}
