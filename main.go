package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/buger/jsonparser"
	"log"
	"net"
	"os"
	"strings"
)

var conf *config

func addTable(conf *config, table string) {
	if len(table) == 0 {
		return
	}
	if ok, _ := inArray(table, conf.Tables); !ok {
		conf.Tables = append(conf.Tables, table)
	}
}

func main() {
	var err error

	pwd, err := os.Getwd()
	hardFail(err)

	var configPath = flag.String("config", "config.toml", "config file")
	var outputPath = flag.String("output", pwd, "output folder")
	flag.Parse()

	log.Println("AutoPeer BGP configuration generator https://github.com/Jamesits/AutoPeer")

	conf = &config{}
	metaData, err := toml.DecodeFile(*configPath, conf)
	hardFail(err)

	// print unknown config variables
	for _, key := range metaData.Undecoded() {
		log.Printf("Warning: unknown option %q", key.String())
	}

	// normalize config
	// backend
	conf.Backend = strings.ToLower(conf.Backend)
	if conf.Backend != "bird1" {
		panic(errors.New("unknown backend"))
	}

	// router id
	var routerId net.IP

	if !govalidator.IsIPv4(conf.RouterID) {
		// try to find a interface name
		ifaces, _ := net.Interfaces()
		for _, iface := range ifaces {
			found := false

			if strings.EqualFold(iface.Name, conf.RouterID) {
				addrs, _ := iface.Addrs()
				for _, addr := range addrs {
					ip := addr.(*net.IPNet).IP
					if ip.To4() == nil || ip.IsLoopback() {
						// is not a IPv4 or is loopback
						// Note: there are some controversy on net.IP.To4() == nil, see:
						// https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6
						continue
					} else {
						routerId = ip
						conf.RouterID = routerId.String()
						found = true
					}
				}
			}

			if found {
				break
			}
		}
	} else {
		routerId = net.ParseIP(conf.RouterID)
	}
	log.Printf("Router ID = %s\n", conf.RouterID)

	// ASN
	log.Printf("ASN = %d\n", conf.Asn)

	// table
	addTable(conf, conf.Table)

	// iterate interfaces
	for ifaceIndex, ifaceDef := range conf.Interfaces {
		log.Printf("Processing interface #%d: %s, IXP %s\n", ifaceIndex, ifaceDef.Name, ifaceDef.Ixp)

		// interface default
		if ifaceDef.MultiHop == 0 {
			conf.Interfaces[ifaceIndex].MultiHop = 1
		}

		// interface inheritance
		if len(ifaceDef.Table) == 0 {
			conf.Interfaces[ifaceIndex].Table = conf.Table
		}

		// table
		addTable(conf, ifaceDef.Table)

		// iterate peers
		for peerIndex, peerDef := range ifaceDef.Peers {
			// get peer information
			peerNetId := getNetIdFromAsn(peerDef.Asn)
			if peerNetId == -1 {
				// no data
				continue
			}

			peerInfoObj := getNetInfoObject(peerNetId)

			name, err := jsonparser.GetString(peerInfoObj, "data", "[0]", "name")
			softFail(err)

			log.Printf("Processing peer #%d: %s (%d)\n", peerIndex, name, peerDef.Asn)

			// peer inheritance
			if len(peerDef.Table) == 0 {
				conf.Interfaces[ifaceIndex].Peers[peerIndex].Table = ifaceDef.Table
			}
			if peerDef.MultiHop == 0 {
				conf.Interfaces[ifaceIndex].Peers[peerIndex].MultiHop = ifaceDef.MultiHop
			}

			// table
			addTable(conf, peerDef.Table)

			// iterate possible sessions
			_, err = jsonparser.ArrayEach(peerInfoObj, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				// TODO: if is_rs_peer = true, ignore them

				// if is not current IXP, skip
				ixp, err := jsonparser.GetString(value, "name")
				hardFail(err)
				if strings.EqualFold(ixp, ifaceDef.Ixp) {
					asn, err := jsonparser.GetInt(value, "asn")
					hardFail(err)

					ipaddr4, err := jsonparser.GetString(value, "ipaddr4")
					softFail(err)

					ipaddr6, err := jsonparser.GetString(value, "ipaddr6")
					softFail(err)

					// add to database
					// TODO: check if overriding exist config
					// TODO: split inheritance
					// TODO: propagate table
					if len(ipaddr4) > 0 {
						log.Printf("Creating IPv4 session %d, %s\n", asn, ipaddr4)
						conf.Interfaces[ifaceIndex].Peers[peerIndex].BgpSessions = append(conf.Interfaces[ifaceIndex].Peers[peerIndex].BgpSessions, bgpSession{
							Name:         cleanString(fmt.Sprintf("%s_%d", name, getUid())),
							PeerEndpoint: net.ParseIP(ipaddr4),
							Asn:          peerDef.Asn,
							Template:     peerDef.Template,
							MultiHop:     peerDef.MultiHop,
							IPv6:         false,
							Table:        peerDef.Table,
						})
					}

					if len(ipaddr6) > 0 {
						log.Printf("Creating IPv6 session %d, %s\n", asn, ipaddr6)
						conf.Interfaces[ifaceIndex].Peers[peerIndex].BgpSessions = append(conf.Interfaces[ifaceIndex].Peers[peerIndex].BgpSessions, bgpSession{
							Name:         cleanString(fmt.Sprintf("%s_%d", name, getUid())),
							PeerEndpoint: net.ParseIP(ipaddr6),
							Asn:          peerDef.Asn,
							Template:     peerDef.Template,
							MultiHop:     peerDef.MultiHop,
							IPv6:         true,
							Table:        peerDef.Table,
						})
					}
				}

			}, "data", "[0]", "netixlan_set")
			softFail(err)
		}
	}

	// invoke generator function
	generator_bird1(conf, outputPath)

	log.Println("AutoPeer done.")
}
