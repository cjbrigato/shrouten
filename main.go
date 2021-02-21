package main

import (
	"flag"
)

func main() {
	addressPtr := flag.String("address", "127.0.0.1", "the shroutening happens on this address")
	portPtr := flag.String("port", "8000", "bonding the shroutening to this tcp prot")
	dbPtr := flag.String("db", "deso.le.db", "the shroutening persiztenza medium")
	flag.Parse()

	store := NewDB(*dbPtr)
	ServeShroutenerForever(store, *addressPtr, *portPtr)
}
