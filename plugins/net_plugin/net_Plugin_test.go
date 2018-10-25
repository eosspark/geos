package net_plugin

import (
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"testing"
)

func makeArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

func TestNetPluginInitialize(t *testing.T) {

	makeArguments("--p2p-listen-endpoint", "127.0.0.1:8100", "--peer-private-key",
		`["EOS5kLsEcZL6ME32rcYLxrWoEzZsHxrqqmFjWAMtzkRNKnS7UpQNR","5Jz3wuG2nitWtU2E8JJXcv9vBTC1rWUtKNAPoc2UU3zX5TBia1x"]`,
		"--p2p-server-address", "127.0.0.1:8100", "--p2p-peer-address", "127.0.0.1:9876", "--p2p-peer-address", "127.0.0.1:7777",
	)

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	netPlugin := NewNetPlugin()
	netPlugin.NetPluginInitialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	prikey, _ := ecc.NewPrivateKey("5Jz3wuG2nitWtU2E8JJXcv9vBTC1rWUtKNAPoc2UU3zX5TBia1x")
	assert.Equal(t, "EOS Test Agent", netPlugin.my.userAgentName)
	assert.Equal(t, prikey.String(), netPlugin.my.privateKeys[prikey.PublicKey()].String())
	assert.Equal(t, "127.0.0.1:8100", netPlugin.my.ListenEndpoint)
	assert.Equal(t, "127.0.0.1:9876", netPlugin.my.suppliedPeers[0])
	assert.Equal(t, "127.0.0.1:8100", netPlugin.my.p2PAddress)
}

func TestNetPlugin(t *testing.T) {
	makeArguments("--p2p-listen-endpoint", "127.0.0.1:8100",
		"--p2p-server-address", "127.0.0.1:8100",
		"--p2p-peer-address", "127.0.0.1:9876",
	)

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	netPlugin := NewNetPlugin()
	netPlugin.NetPluginInitialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	netPlugin.PluginStartup()

}
