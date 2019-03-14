package main

var fnTable = map[string]func([]string) int {
	"init": actionInit,
	"list-ix": actionListIx,
	"list-peer": actionListPeer,
	"add-peer": actionAddPeer,
	"apply": actionApply,
}
