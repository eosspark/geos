package types

import (
	"fmt"
	"testing"

	. "github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/stretchr/testify/assert"
)

func Test_nextPowerOf2(t *testing.T) {
	assert.EqualValues(t, 1, nextPowerOf2(1))
	assert.EqualValues(t, 2, nextPowerOf2(2))
	assert.EqualValues(t, 4, nextPowerOf2(3))
	assert.EqualValues(t, 4, nextPowerOf2(4))
	assert.EqualValues(t, 8, nextPowerOf2(5))
}

func Test_clsPower2(t *testing.T) {
	assert.EqualValues(t, 10, clzPower2(1<<10))
	assert.EqualValues(t, 5, clzPower2(1<<5))
	assert.EqualValues(t, 1, clzPower2(2))
	assert.EqualValues(t, 0, clzPower2(1))
}

func Test_calculateMaxDepth(t *testing.T) {
	assert.EqualValues(t, 0, calculateMaxDepth(0))
	assert.EqualValues(t, 1, calculateMaxDepth(1))
	assert.EqualValues(t, 2, calculateMaxDepth(2))
	assert.EqualValues(t, 3, calculateMaxDepth(3))
	assert.EqualValues(t, 3, calculateMaxDepth(4))
	assert.EqualValues(t, 4, calculateMaxDepth(5))
	assert.EqualValues(t, 4, calculateMaxDepth(6))
	assert.EqualValues(t, 4, calculateMaxDepth(7))
	assert.EqualValues(t, 4, calculateMaxDepth(8))
}

func Test_moveNodes(t *testing.T) {
	src := make([]DigestType, 1)
	src[0] = *crypto.NewSha256String("26b25d457597a7b0463f9620f666dd10aa2c4373a505967c7c8d70922a2d6ece")
	dst := []DigestType{{}}
	moveNodes(&dst, &src)
	assert.Equal(t, *crypto.NewSha256String("26b25d457597a7b0463f9620f666dd10aa2c4373a505967c7c8d70922a2d6ece"), dst[0])
}

func TestMerkle(t *testing.T) {
	ids := []DigestType{*crypto.NewSha256String("00000043df9347b6d053a03a78499bc420acb05c5c3bec6acbd8d37a68b3f195"),
		*crypto.NewSha256String("00000042332b70edb826b578924218b8b509c3dcb2011608829f9f5f85d983b0")}

	assert.Equal(t, "ad2b926aa8d86b1ed9dc70b47b48778c43a2406eabba236d81580286640f78d5", Merkle(ids).String())
}

func TestIncrementalMerkle_GetRoot(t *testing.T) {
	m := new(IncrementalMerkle)
	m.Append(*crypto.NewSha256String("00000043df9347b6d053a03a78499bc420acb05c5c3bec6acbd8d37a68b3f195"))
	m.Append(*crypto.NewSha256String("00000042332b70edb826b578924218b8b509c3dcb2011608829f9f5f85d983b0"))

	assert.Equal(t, "ad2b926aa8d86b1ed9dc70b47b48778c43a2406eabba236d81580286640f78d5", m.GetRoot().String())
}

func TestIncrementalMerkle_Append(t *testing.T) {
	m := new(IncrementalMerkle)
	ids := []DigestType{*crypto.NewSha256String("00000043df9347b6d053a03a78499bc420acb05c5c3bec6acbd8d37a68b3f195"),
		*crypto.NewSha256String("00000042332b70edb826b578924218b8b509c3dcb2011608829f9f5f85d983b0"),
		*crypto.NewSha256String("0000002033c7b70cdeadd2e71aa0f4adf38579904656754d3f4d14b68444bc08"),
		*crypto.NewSha256String("00000030cf4390afd0b1755c166a83ac3c49a9c017ffb1c5ffb41df2e03e614b")}
	m.Append(ids[0])
	assert.Equal(t, "[00000043df9347b6d053a03a78499bc420acb05c5c3bec6acbd8d37a68b3f195]", fmt.Sprintf("%s", m.ActiveNodes))
	m.Append(ids[1])
	assert.Equal(t, "[ad2b926aa8d86b1ed9dc70b47b48778c43a2406eabba236d81580286640f78d5]", fmt.Sprintf("%s", m.ActiveNodes))
	m.Append(ids[2])
	assert.Equal(t, "[0000002033c7b70cdeadd2e71aa0f4adf38579904656754d3f4d14b68444bc08 ad2b926aa8d86b1ed9dc70b47b48778c43a2406eabba236d81580286640f78d5 d4fa7805b9676057cc98e248fad11e7405b6c419a626e7b88cf64fd85c6dd160]", fmt.Sprintf("%s", m.ActiveNodes))
	m.Append(ids[3])
	assert.Equal(t, "[80e1b3d420a7ff9653f965867b42b4127f9a58ce2432e960548f0fa82377fe80]", fmt.Sprintf("%s", m.ActiveNodes))

	merkel := new(IncrementalMerkle)

	var merkelAdd = func(t *testing.T, merkel *IncrementalMerkle, o interface{}) {
		test := crypto.Hash256(o)
		assert.Equal(t, merkel.Append(*test), merkel.GetRoot())
	}

	merkelAdd(t, merkel, "eos")
	assert.Equal(t, len(merkel.ActiveNodes), 1)
	merkelAdd(t, merkel, "test")
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
