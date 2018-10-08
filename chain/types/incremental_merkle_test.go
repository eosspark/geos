package types

import (
	"github.com/stretchr/testify/assert"
	"github.com/eosspark/eos-go/crypto"
	"testing"
)

func Test_append(t *testing.T) {
	merkel := new(IncrementalMerkle)

	var merkelAdd = func(t *testing.T, merkel *IncrementalMerkle, o interface{}) {
		test := crypto.Hash256(o)
		assert.Equal(t, merkel.Append(test), merkel.GetRoot())
	}

	merkelAdd(t, merkel, "eos")
	assert.Equal(t, len(merkel.ActiveNodes), 1)
	merkelAdd(t, merkel, "test")
	assert.Equal(t, len(merkel.ActiveNodes), 1)
	merkelAdd(t, merkel, "merkel")
	assert.Equal(t, len(merkel.ActiveNodes), 3)
	merkelAdd(t, merkel, "go")
	assert.Equal(t, len(merkel.ActiveNodes), 1)
	merkelAdd(t, merkel, "golang")
	assert.Equal(t, len(merkel.ActiveNodes), 3)
	merkelAdd(t, merkel, "cpp")
	assert.Equal(t, len(merkel.ActiveNodes), 3)
	merkelAdd(t, merkel, "incremental")
	assert.Equal(t, len(merkel.ActiveNodes), 4)
	merkelAdd(t, merkel, "append")
	assert.Equal(t, len(merkel.ActiveNodes), 1)
}

func Test_calculateMaxDepth(t *testing.T) {
	depth := make(map[int]int)
	for i := 0; i <= 32; i++ {
		depth[calculateMaxDepth(uint64(i))]++
	}

	assert.Equal(t, depth[2], depth[1])
	assert.Equal(t, depth[3], 2*depth[2])
	assert.Equal(t, depth[4], 2*depth[3])
	assert.Equal(t, depth[5], 2*depth[4])
	assert.Equal(t, depth[6], 2*depth[5])

	np2 := make(map[uint64]int)
	for i := 0; i <= 32; i++ {
		np2[nextPowerOf2(uint64(i))]++
	}
	assert.Equal(t, np2[2], np2[1])
	assert.Equal(t, np2[4], 2*np2[2])
	assert.Equal(t, np2[8], 2*np2[4])
	assert.Equal(t, np2[16], 2*np2[8])
	assert.Equal(t, np2[32], 2*np2[16])

	for i, j := 2, 1; i <= 1024; i *= 2 {
		assert.Equal(t, j, clzPower2(uint64(i)))
		j++
	}
}

func Test_moveNodes(t *testing.T) {
	src := make([]crypto.Sha256, 1)
	src[0] = *crypto.NewSha256String("26b25d457597a7b0463f9620f666dd10aa2c4373a505967c7c8d70922a2d6ece")
	dst := []crypto.Sha256{{}}
	moveNodes(&dst, &src)
	assert.Equal(t, [][4]uint64{{0xb0a79775455db226, 0x10dd66f620963f46, 0x7c9605a573432caa, 0xce6e2d2a92708d7c}}, dst)
}
