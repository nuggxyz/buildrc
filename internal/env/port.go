//	port.go
//	golang/env
//
//	Created by walter on 2023-01-10.
//	Copyright © 2023, nugg.xyz LLC. All rights reserved.
//	---------------------------------------------------------------------
//	adapted from hashicorp/go-retryablehttp
//	Copyright © 2015, HashiCorp, Inc. MPL 2.0
//	---------------------------------------------------------------------
//

package env

import (
	"net"
)

// GetOpenLocalPort asks the kernel for a free open port that is ready to use.
func GetOpenLocalPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetPort is deprecated, use GetOpenLocalPort instead
// Ask the kernel for a free open port that is ready to use
func MustGetOpenLocalPort() int {
	port, err := GetOpenLocalPort()
	if err != nil {
		panic(err)
	}
	return port
}
