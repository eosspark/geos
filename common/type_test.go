package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_config(t *testing.T) {
	fmt.Println(DefaultConfig)
}

type fc struct {
	Id     uint32
	Name   string
	Flag   bool
	Sex    byte
	Value  [4]uint64
	Values []int
}

func Test_Hash(t *testing.T) {
	//int fc::sha256::hash<int>(2)
	fmt.Printf("int %x\n", Hash(2))
	assert.Equal(t, [4]uint64{0xb0a79775455db226, 0x10dd66f620963f46, 0x7c9605a573432caa, 0xce6e2d2a92708d7c}, Hash(2))

	//string fc::sha256::hash<string>("encode")
	fmt.Printf("string %x\n", Hash("encode"))
	assert.Equal(t, [4]uint64{0x30752d3959341337, 0x43c58a183f07a4a6, 0x2a13f5f54922f828, 0x44b38a0c1ecbaf2}, Hash("encode"))

	//slice fc::sha256::hash<vector<int>>({1,0,0,8,6})
	fmt.Printf("slice %x\n", Hash([]int{1, 0, 0, 8, 6}))
	assert.Equal(t, [4]uint64{0xa36ef456b37a9908, 0x55f10469a9354173, 0x8afb9e5ec5df83bf, 0xc7a4c2ad8c5dd1b2}, Hash([]int{1, 0, 0, 8, 6}))

	//map fc::sha256::hash<map<int,int>>({{1,1},{2,3}})
	fmt.Printf("map %x\n", Hash(map[int]int{1: 1, 2: 3}))
	//maybe not equal of unordered map
	//assert.Equal(t, [4]uint64{0x2a345a16b30e9ac0,0xc08e5c02c109d722,0x20a60382171ec7c5,0xfc79fbdce9986e41}, Hash(map[int]int{1:1,2:3}))

	//array fc::sha256::hash<int[4]>({1,2,3,4})
	fmt.Printf("array %x\n", Hash([4]int{1, 2, 3, 4}))
	assert.Equal(t, [4]uint64{0xc6895e8c3f82d1af, 0xe67215330aa2b680, 0xe11dd3dd62c3b13a, 0xb3a379ee2dd73853}, Hash([4]int{1, 2, 3, 4}))

	//struct
	var fcs = fc{1, "a", false, 'M', [4]uint64{1, 2, 3, 4}, []int{6, 7, 8}}
	fmt.Printf("struct: %x\n", Hash(fcs))
	assert.Equal(t, [4]uint64{0x28ca16c792f63d06, 0x65ddc8cb182e3e1a, 0x3cf688a2caa80b54, 0xdbee41d1d901dc88}, Hash(fcs))

	//pair fc::sha256::hash(make_pair(1,"a"))
	fmt.Printf("pair: %x\n", Hash(Pair{1, "a"}))
	assert.Equal(t, [4]uint64{0x8ef475122a81d373, 0xed1da097ad920e4d, 0x1d810c6eb257f, 0x64a37073101a8f65}, Hash(Pair{1, "a"}))

	//twice
	//b1,_ := rlp.EncodeToBytes(1)
	//b2,_ := rlp.EncodeToBytes(2)
	//
	////b3 := append(b1, b2...)
	//h3 := sha256.New()
	//_, _ = h3.Write(b1)
	//hashed := h3.Sum(b2)
	////h1 := sha256.New()
	////_, _ = h1.Write(b1)
	////hashed1 := h1.Sum(nil)
	////
	////h2 := sha256.New()
	////_, _ = h2.Write(b2)
	////hashed2 := h2.Sum(nil)
	////
	////hashed := append(hashed1, hashed2...)
	//
	////fmt.Println(hashed)
	//
	//var result [4]uint64
	//
	//result[0] = binary.LittleEndian.Uint64(hashed[:8])
	//result[1] = binary.LittleEndian.Uint64(hashed[8:16])
	//result[2] = binary.LittleEndian.Uint64(hashed[16:24])
	//result[3] = binary.LittleEndian.Uint64(hashed[24:32])
	//
	//fmt.Printf("twice: %x\n", result)
}
