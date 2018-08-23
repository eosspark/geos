package main

import (
	"encoding/hex"
	"flag"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/p2p"
	"log"
	"time"
)

var p2pAddr = flag.String("p2p-addr", "127.0.0.1:9876", "P2P socket connection")
var chainID = flag.String("chain-id", "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f", "Chain id")
var networkVersion = flag.Int("network-version", 1206, "Network version")
var httpServerAddress = flag.String("http-server-address", "http://127.0.0.1:8888", "The local IP and port to listen for incoming http connections; set blank to disable.")
var walletServerAddress = flag.String("wallet-server-address", "http://127.0.0.1:8900", "The local IP and port to listen for incoming http connections; set blank to disable.")

func main() {
	flag.Parse()
	cID, err := hex.DecodeString(*chainID)
	if err != nil {
		log.Fatal(err)
	}
	data, _ := common.DecodeIDTypeByte(cID)
	client := p2p.NewClient(*p2pAddr, common.ChainIDType(data), uint16(*networkVersion))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = client.StartConnect()
	if err != nil {
		log.Fatal(err)
	}

}
