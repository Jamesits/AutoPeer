package main

import "net"

type config struct {
	Backend       string         `toml:"backend"`   // generated config format
	RouterID      string         `toml:"router_id"` // Router ID, either an IP or an interface name
	Asn           int64          `toml:"asn"`       // our ASN
	Interfaces    []netInterface `toml:"interface"`
	Announcements []string       `toml:"announcements"`

	// will be auto-filled
	Tables []string

	// inherited
	Table string `toml:"table"`
}

type netInterface struct {
	Name     string `toml:"name"` // OS name for the network interface
	Ixp      string `toml:"ixp"`  // IXP name
	Peers    []peer `toml:"peer"` // peers on that interface
	MultiHop int8   `toml:"multihop"`

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
	IPv6 bool   `toml:"ipv6"`
}
