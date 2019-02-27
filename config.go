package main

import "net"

var conf *config

type config struct {
	Format     string         `toml:"format"` // generated config format
	Asn        int64          `toml:"asn"`    // our ASN
	Interfaces []netInterface `toml:"interface"`

	// inherited
	Table string `toml:"table"`

	// internal
	tables []string
}

type netInterface struct {
	Name     string `toml:"name"` // OS name for the network interface
	Ixp      string `toml:"ixp"`  // IXP name
	Peers    []peer `toml:"peer"` // peers on that interface
	MultiHop int8   `toml:"multihop"`
	// multihop:
	// 0 is a "default" value (do not override application/template)
	// 1 is "explicitly direct"

	// inherited
	Table string `toml:"table"`
}

type peer struct {
	Asn      int64  `toml:"asn"`
	Template string `toml:"template"` // template name as in generated config

	// inherited
	Table    string `toml:"table"`
	MultiHop int8   `toml:"multihop"`

	// will be auto-generated
	BgpSessions []bgpSession `toml:"session"`
}

type bgpSession struct {
	PeerEndpoint net.IP `toml:"peer_endpoint"`

	// inherited
	Asn      int64  `toml:"asn"`
	Template string `toml:"template"`
	MultiHop int8   `toml:"multihop"`
	Table    string `toml:"table"`

	// will be auto-generated
	Name string `toml:"name"`

	// internal
	ipv6      bool
	processed bool
}
