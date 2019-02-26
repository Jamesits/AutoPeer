package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
)

func generator_bird1(conf *config, outputPath *string) {
	var err error

	log.Println("Generating bird v1 config files...")

	// create output files
	configFile4, err := os.Create(path.Join(*outputPath, "bird.conf"))
	hardFail(err)
	defer configFile4.Close()
	configFile6, err := os.Create(path.Join(*outputPath, "bird6.conf"))
	hardFail(err)
	defer configFile6.Close()

	// write common prefix
	const prefix = `# Generated by AutoPeer
protocol device {}
`
	_, err = fmt.Fprint(configFile4, prefix)
	hardFail(err)
	_, err = fmt.Fprint(configFile6, prefix)
	hardFail(err)

	// write router id
	_, err = fmt.Fprintf(configFile4, "router id %s;\n", conf.RouterID)
	hardFail(err)
	_, err = fmt.Fprintf(configFile6, "router id %s;\n", conf.RouterID)
	hardFail(err)

	// write includes
	_, err = fmt.Fprint(configFile4, "include \"include4.conf\";\n")
	hardFail(err)
	_, err = fmt.Fprint(configFile6, "include \"include6.conf\";\n")
	hardFail(err)

	// write tables
	for _, table := range conf.Tables {
		_, err = fmt.Fprintf(configFile4, "table %s;\n", table)
		hardFail(err)
		_, err = fmt.Fprintf(configFile6, "table %s;\n", table)
		hardFail(err)
	}

	// write static announcement block
	_, err = fmt.Fprint(configFile4, "protocol static bgp_announcement {\n")
	hardFail(err)
	_, err = fmt.Fprint(configFile6, "protocol static bgp_announcement {\n")
	hardFail(err)
	if len(conf.Table) > 0 {
		_, err = fmt.Fprintf(configFile4, "    table %s;\n", conf.Table)
		hardFail(err)
		_, err = fmt.Fprintf(configFile6, "    table %s;\n", conf.Table)
		hardFail(err)
	}
	for _, route := range conf.Announcements {
		_, ipnet, err := net.ParseCIDR(route)
		hardFail(err)
		if ipnet.IP.To4() != nil {
			// IPv4
			_, err = fmt.Fprintf(configFile4, "    route %s blackhole;\n", ipnet.String())
			hardFail(err)
		} else {
			// IPv6
			_, err = fmt.Fprintf(configFile6, "    route %s blackhole;\n", ipnet.String())
			hardFail(err)
		}
	}
	_, err = fmt.Fprint(configFile4, "}\n")
	hardFail(err)
	_, err = fmt.Fprint(configFile6, "}\n")
	hardFail(err)

	// write BGP session info
	for _, iface := range conf.Interfaces {
		for _, peer := range iface.Peers {
			for _, session := range peer.BgpSessions {
				var configFile io.Writer
				if session.IPv6 == false {
					configFile = configFile4
				} else {
					configFile = configFile6
				}

				_, err = fmt.Fprintf(configFile, "protocol bgp %s ", session.Name)
				hardFail(err)
				if len(session.Template) > 0 {
					_, err = fmt.Fprintf(configFile, "from %s {\n", session.Template)
					hardFail(err)
				} else {
					_, err = fmt.Fprintf(configFile, "{\n")
					hardFail(err)
				}

				if len(session.Table) > 0 {
					_, err = fmt.Fprintf(configFile, "    table %s;\n", session.Table)
					hardFail(err)
				}

				_, err = fmt.Fprintf(configFile, "    local as %d;\n", conf.Asn)
				hardFail(err)

				_, err = fmt.Fprintf(configFile, "    neighbor %s as %d;\n", session.PeerEndpoint, session.Asn)
				hardFail(err)

				if session.MultiHop > 1 {
					_, err = fmt.Fprintf(configFile, "    multihop %d;\n", session.MultiHop)
					hardFail(err)
				}

				_, err = fmt.Fprint(configFile, "}\n\n")
				hardFail(err)
			}
		}
	}
}
