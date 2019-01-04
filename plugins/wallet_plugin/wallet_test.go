package wallet_plugin

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletPassword(t *testing.T) {
	password := "PW5J6XpRE6Lur3Crv7QVsGoX1hMk1QPMGfyoT24kVfMTnarZ524xv"
	re := hash512(password)

	str := "9edf1073e71bc2ffbc23ae7a8faac576d8b7c703ce863224e270f1a5e6c4fde40c75b459d1cfe852434ad18d3d32881f48859ded61e5c55b4b43d8261a360245"
	bytes, _ := hex.DecodeString(str)

	assert.Equal(t, bytes, re)
}
