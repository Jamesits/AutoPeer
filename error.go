package main

import "log"

// check every fucking err
func hardFail(e error) {
	if e != nil {
		panic(e)
	}
}

// check every fucking err and forget them
func softFail(e error) {
	if e != nil {
		log.Printf("[ERROR] %s", e)
	}
}
