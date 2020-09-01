package main

import (
	"net"

	"github.com/coolray-dev/raydash/docs"
	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/setting"
)

func setupSwagger() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Log.WithError(err).Fatal("Error Getting Interface Addr")
	}

	for _, address := range addrs {

		// Check if IP is a loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				docs.SwaggerInfo.Host = ipnet.IP.String() + ":" + setting.Config.GetString("app.port")
			}

		}
	}
}
