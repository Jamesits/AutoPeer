package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/buger/jsonparser"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var err error

	pwd, err := os.Getwd()
	hardFail(err)

	var configPath = flag.String("config", "autopeer.toml", "config file")
	var outputPath = flag.String("output", pwd, "output folder")
	var format = flag.String("format", "", "override output format")
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
	if len(*format) > 0 {
		// override
		conf.Format = *format
	}
	if len(conf.Format) == 0 {
		// default
		conf.Format = "bird1"
	}
	conf.Format = strings.ToLower(conf.Format)
	if conf.Format != "bird1" {
		panic(errors.New("unknown backend"))
	}

	// ASN
	log.Printf("ASN = %d\n", conf.Asn)

	// table
	addTable(conf, conf.Table)

	// iterate interfaces
	for ifaceIndex := range conf.Interfaces {
		ifaceDef := &conf.Interfaces[ifaceIndex]
		log.Printf("Processing interface #%d: %s, IXP %s\n", ifaceIndex, ifaceDef.Name, ifaceDef.Ixp)

		// interface inheritance
		if len(ifaceDef.Table) == 0 {
			ifaceDef.Table = conf.Table
		}

		// table
		addTable(conf, ifaceDef.Table)

		// iterate peers
		for peerIndex := range ifaceDef.Peers {
			peerDef := &ifaceDef.Peers[peerIndex]
			log.Printf("Processing peer #%d: %d\n", peerIndex, peerDef.Asn)

			// peer inheritance
			if len(peerDef.Table) == 0 {
				peerDef.Table = ifaceDef.Table
			}
			if peerDef.MultiHop == 0 {
				peerDef.MultiHop = ifaceDef.MultiHop
			}

			// table
			addTable(conf, peerDef.Table)

			// list for unprocessed sessions
			oldSessions := peerDef.BgpSessions

			// list for processed sessions
			peerDef.BgpSessions = nil

			// try to get peer information
			peerNetId := getNetIdFromAsn(peerDef.Asn)
			if peerNetId != -1 {
				// retrieved data from data source, fill config
				peerInfoObj := getNetInfoObject(peerNetId)

				name, err := jsonparser.GetString(peerInfoObj, "data", "[0]", "name")
				softFail(err)

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

						var ipaddrs []string = nil

						if len(ipaddr4) > 0 {
							// we have an IPv4 session
							ipaddrs = append(ipaddrs, ipaddr4)
						}

						if len(ipaddr6) > 0 {
							// we have an ipv6 session
							ipaddrs = append(ipaddrs, ipaddr6)
						}

						for _, ipaddr := range ipaddrs {
							log.Printf("Creating session %d, %s\n", asn, ipaddr)

							ip := net.ParseIP(ipaddr)

							// inheritance
							newSession := bgpSession{
								Name:         cleanString(fmt.Sprintf("%s_%d", name, getUid())),
								PeerEndpoint: ip,
								Asn:          peerDef.Asn,
								Template:     peerDef.Template,
								MultiHop:     peerDef.MultiHop,
								ipv6:         isIPv6(ip),
								Table:        peerDef.Table,
							}

							// try to find an existing config
							for key, session := range oldSessions {
								if session.PeerEndpoint.Equal(newSession.PeerEndpoint) {
									// we got a match
									log.Printf("Found existing config for %d, %s\n", asn, ipaddr)
									oldSessions[key].processed = true

									// override config
									if len(session.Name) > 0 {
										newSession.Name = session.Name
									}

									if session.Asn > 0 {
										newSession.Asn = session.Asn
									}

									if len(session.Template) > 0 {
										newSession.Template = session.Template
									}

									if session.MultiHop > 0 && session.MultiHop != ifaceDef.MultiHop {
										newSession.MultiHop = session.MultiHop
									}

									if len(session.Table) > 0 {
										newSession.Table = session.Table
									}
									break
								}
							}

							// propagate table
							addTable(conf, newSession.Table)

							// add back to list
							peerDef.BgpSessions = append(
								peerDef.BgpSessions,
								newSession,
							)
						}
					}

				}, "data", "[0]", "netixlan_set")
				softFail(err)
			}

			// add all unprocessed sessions back to list
			for _, session := range oldSessions {
				if session.processed == false {
					log.Printf("Found individual config %s\n", session.Name)
					session.processed = true

					// inheritance
					if session.Asn == 0 {
						session.Asn = peerDef.Asn
					}
					if len(session.Template) == 0 {
						session.Template = peerDef.Template
					}
					if session.MultiHop == 0 {
						session.MultiHop = peerDef.MultiHop
					}
					if len(session.Table) == 0 {
						session.Table = peerDef.Table
					}
					session.ipv6 = isIPv6(session.PeerEndpoint)

					// propagate table
					addTable(conf, session.Table)

					// add back to list
					peerDef.BgpSessions = append(
						peerDef.BgpSessions,
						session,
					)
				}
			}
		}
	}

	// invoke generator function
	generator_bird1(conf, outputPath)

	log.Println("AutoPeer done.")
}
