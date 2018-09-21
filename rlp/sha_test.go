package rlp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fc struct {
	Id     uint32
	Name   string
	Flag   bool
	Sex    byte
	Value  [4]uint64
	Values []int
}

type pair struct {
	T interface{}
	R interface{}
}

func TestSha256(t *testing.T) {
	fmt.Println(Sha256{})
	assert.Equal(t, [4]uint64{0}, Sha256{}.Hash_)

	h1 := Hash(2)
	fmt.Println(h1)
	assert.Equal(t, [4]uint64{0xb0a79775455db226, 0x10dd66f620963f46, 0x7c9605a573432caa, 0xce6e2d2a92708d7c},
		NewSha256("26b25d457597a7b0463f9620f666dd10aa2c4373a505967c7c8d70922a2d6ece").Hash_)
}

func TestHash(t *testing.T) {
	//int fc::sha256::hash<int>(2)
	fmt.Printf("int %x\n", Hash(2).Hash_, Hash(2).Bytes())
	assert.Equal(t, [4]uint64{0xb0a79775455db226, 0x10dd66f620963f46, 0x7c9605a573432caa, 0xce6e2d2a92708d7c}, Hash(2).Hash_)

	//string fc::sha256::hash<string>("encode")
	fmt.Printf("string %x\n", Hash("encode").Hash_)
	assert.Equal(t, [4]uint64{0x30752d3959341337, 0x43c58a183f07a4a6, 0x2a13f5f54922f828, 0x44b38a0c1ecbaf2}, Hash("encode").Hash_)

	//slice fc::sha256::hash<vector<int>>({1,0,0,8,6})
	fmt.Printf("slice %x\n", Hash([]int{1, 0, 0, 8, 6}).Hash_)
	assert.Equal(t, [4]uint64{0xa36ef456b37a9908, 0x55f10469a9354173, 0x8afb9e5ec5df83bf, 0xc7a4c2ad8c5dd1b2}, Hash([]int{1, 0, 0, 8, 6}).Hash_)

	//map fc::sha256::hash<map<int,int>>({{1,1},{2,3}})
	fmt.Printf("map %x\n", Hash(map[int]int{1: 1, 2: 3}).Hash_)
	//maybe not equal of unordered map
	//assert.Equal(t, [4]uint64{0x2a345a16b30e9ac0,0xc08e5c02c109d722,0x20a60382171ec7c5,0xfc79fbdce9986e41}, Hash(map[int]int{1:1,2:3}))

	//array fc::sha256::hash<int[4]>({1,2,3,4})
	fmt.Printf("array %x\n", Hash([4]int{1, 2, 3, 4}).Hash_)
	assert.Equal(t, [4]uint64{0xc6895e8c3f82d1af, 0xe67215330aa2b680, 0xe11dd3dd62c3b13a, 0xb3a379ee2dd73853}, Hash([4]int{1, 2, 3, 4}).Hash_)

	//struct
	var fcs = fc{1, "a", false, 'M', [4]uint64{1, 2, 3, 4}, []int{6, 7, 8}}
	fmt.Printf("struct: %x\n", Hash(fcs).Hash_)
	assert.Equal(t, [4]uint64{0x28ca16c792f63d06, 0x65ddc8cb182e3e1a, 0x3cf688a2caa80b54, 0xdbee41d1d901dc88}, Hash(fcs).Hash_)

	//pair fc::sha256::hash(make_pair(1,"a"))
	fmt.Printf("pair: %x\n", Hash(pair{1, "a"}).Hash_)
	assert.Equal(t, [4]uint64{0x8ef475122a81d373, 0xed1da097ad920e4d, 0x1d810c6eb257f, 0x64a37073101a8f65}, Hash(pair{1, "a"}).Hash_)

	//tuple3 pack*3 [1,"a",2]
	fmt.Printf("tuple3: %x\n", Hash(pair{pair{1, "a"}, 2}).Hash_)
	assert.Equal(t, [4]uint64{0x84f4091030117a8c, 0x33507294a01fbdeb, 0x37b9085b42c1d687, 0x1c65536bce56ea7d}, Hash(pair{pair{1, "a"}, 2}).Hash_)

	//tuple4 pack*4 [1,"a",2,['x','y']]
	fmt.Printf("tuple4: %x\n", Hash(pair{pair{1, "a"}, pair{2, []byte{'x', 'y'}}}).Hash_)
	assert.Equal(t, [4]uint64{0x1254b5d88231112, 0xb736a7d556d3bb30, 0x461a5062a842f8cf, 0xeee703f7e79d4dc}, Hash(pair{pair{1, "a"}, pair{2, []byte{'x', 'y'}}}).Hash_)

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


