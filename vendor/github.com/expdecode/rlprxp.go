package main

import (
	"encoding/json"
	"fmt"
	"github.com/expdecode/rlp"
)

type ProducerKey struct {
	ProducerName    uint64   //account_name
	BlockSigningKey [33]byte //public_key_type
}
type ProducerScheduleType struct {
	Version   uint32
	Producers []*ProducerKey
}
type OptionalPST struct {
	// ProducerValid bool //eos raw.hpp 278
	// Pst
	Pst ProducerScheduleType
}
type Optionalstr struct {
	Shape        int16
	Name         string
	NewProducers *OptionalPST `eos:"optional"`
}

func main() {
	// rlp.Debug = true
	producekey := &ProducerKey{
		ProducerName:    19,
		BlockSigningKey: [33]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9},
	}
	produce := ProducerScheduleType{
		Version:   56,
		Producers: []*ProducerKey{producekey},
	}
	optionalpst := &OptionalPST{
		Pst: produce,
	}
	msg := Optionalstr{
		Shape:        99,
		Name:         "walker",
		NewProducers: optionalpst,
	}
	data, err := rlp.EncodeToBytes(msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	// var kk Optionalstr

	kk := &Optionalstr{}
	err = rlp.DecodeBytes(data, kk)
	fmt.Printf("%v,%v\n", err, kk)

	aa, _ := json.Marshal(kk)
	fmt.Println(string(aa))

}

// [99 0 6 119 97 108 107 101 114 1 56 0 0 0 1 19 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 9]
// {
// 	"Shape": 99,
// 	"Name": "walker",
// 	"NewProducers": {
// 		"Pst": {
// 			"Version": 0,
// 			"Producers": [{
// 				"ProducerName": 19,
// 				"BlockSigningKey": [1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9]
// 			}]
// 		}
// 	}
// }
