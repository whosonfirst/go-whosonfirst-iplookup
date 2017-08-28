package iplookup

import (
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net"
)

type Provider interface {
	QueryString(string) (spr.StandardPlacesResult, error)
	Query(net.IP) (spr.StandardPlacesResult, error)
}
