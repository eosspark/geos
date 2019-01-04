package unittests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type typeName = string

var maxSerializationTime = common.Seconds(1) // some test machines are very slow

// verify that round trip conversion, via bytes, reproduces the exact same data
func verifyByteRoundTripConversion(abis *abi.AbiSerializer, name typeName, data *common.Variants) common.Variant {
	bytes1 := abis.VariantToBinary(name, data, maxSerializationTime)
	var2 := abis.BinaryToVariant(name, bytes1, maxSerializationTime, false)
	//Rbytes, err := json.Marshal(var2)
	//if err != nil {
	//	fmt.Println("marshal is error: ", err)
	//}
	//fmt.Println(string(Rbytes))
	bytes2 := abis.VariantToBinary(name, &var2, maxSerializationTime)
	if bytes.Compare(bytes1, bytes2) != 0 {
		Throw(fmt.Errorf("It's not equal!!"))
	}
	return var2
}

func verifyRoundTripConversion5(abis *abi.AbiSerializer, name typeName, jsonStr string, hex string, expectedJSON string) {
	var variant1 common.Variants
	err := json.Unmarshal([]byte(jsonStr), &variant1)
	if err != nil {
		fmt.Println("verifyRoundTripConversion unmarshal variant1 is error: ", err)
		Throw(err)
	}

	bytes1 := abis.VariantToBinary(name, &variant1, maxSerializationTime)
	if bytes.Compare(bytes1, []byte(hex)) != 0 {
		fmt.Println("result:  ", bytes1, hex)
		Throw(fmt.Errorf("It's not equal!!"))
	}

	var2 := abis.BinaryToVariant(name, bytes1, maxSerializationTime, false)
	result, err := json.Marshal(var2)
	if err != nil {
		fmt.Println("verifyRoundTripConversion marshal var2 is error: ", err)
		Throw(err)
	}
	if strings.Compare(string(result), expectedJSON) != 0 {
		fmt.Println("result:  ", string(result), hex)
		Throw(fmt.Errorf("It's not equal!!"))
	}

	bytes2 := abis.VariantToBinary(name, &var2, maxSerializationTime)
	if bytes.Compare(bytes2, []byte(hex)) != 0 {
		fmt.Println("result:  ", bytes2, hex)
		Throw(fmt.Errorf("It's not equal!!"))
	}
}

func verifyRoundTripConversion4(abis *abi.AbiSerializer, name typeName, jsonStr string, hex string) {
	verifyRoundTripConversion5(abis, name, jsonStr, hex, jsonStr)
}

func getResolver() func(account common.AccountName) *abi.AbiSerializer {
	abiDef := abi.AbiDef{}
	return func(account common.AccountName) *abi.AbiSerializer {
		return abi.NewAbiSerializer(chain.EosioContractAbi(abiDef), maxSerializationTime)
	}
}

func verifyTypeRoundTripConversion(t types.ContractTypesInterface, abis *abi.AbiSerializer, name typeName, args *common.Variants) *common.Variants {
	var args2 common.Variants
	Try(func() {
		bytes1 := abis.VariantToBinary(name, args, maxSerializationTime)
		abi.FromVariantToActionData(&t, args, getResolver(), maxSerializationTime)

		abi.ToVariantFromActionData(&t, &args2, getResolver(), maxSerializationTime)

		r, err := json.Marshal(args2)
		fmt.Println(r, err)

		bytes2 := abis.VariantToBinary(name, &args2, maxSerializationTime)

		if bytes.Compare(bytes1, bytes2) != 0 {
			fmt.Println("result:  ", bytes1, bytes2)
			Throw(fmt.Errorf("It's not equal!!"))
		}

	}).FcLogAndRethrow().End()
	return &args2
}

func TestUintTypes(t *testing.T) {
	Try(func() {
		abiReader := strings.NewReader(currencyABI)
		abidef, err := abi.NewABI(abiReader)
		assert.NoError(t, err)

		var testData string = `{
		   "amount64" : 64,
		   "amount32" : 32,
		   "amount16" : 16,
		   "amount8"  : 8
		}`
		data := common.Variants{}
		strRead := strings.NewReader(testData)
		err = json.NewDecoder(strRead).Decode(&data)
		assert.NoError(t, err)
		verifyByteRoundTripConversion(abi.NewAbiSerializer(abidef, maxSerializationTime), "transfer", &data)

	}).FcLogAndRethrow().End()
}

func abiDefFromString(abiStr string) *abi.AbiDef {
	abiReader := strings.NewReader(abiStr)
	abiDef, err := abi.NewABI(abiReader)
	if err != nil {
		panic(err.Error())
	}
	return abiDef
}

func TestGeneral(t *testing.T) {
	Try(func() {
		abiDef := chain.EosioContractAbi(*abiDefFromString(myABI))

		abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)
		fmt.Println(abis)

		data := common.Variants{}
		strRead := strings.NewReader(myOther)
		err := json.NewDecoder(strRead).Decode(&data)
		assert.NoError(t, err)
		fmt.Println(data)
		verifyByteRoundTripConversion(abis, "A", &data)

	}).Catch(func(e interface{}) {
		panic(e)
	}).End()
}

func TestAbiCycle(t *testing.T) {
	Try(func() {
		//abiDef := chain.EosioContractAbi(*abiDefFromString(typedefCycleABI))
		//abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)
		//fmt.Println(abis)

	}).Catch(func(e interface{}) {
		panic(e)
	}).End()

}

func TestLinkauth(t *testing.T) {
	Try(func() {
		abiDef := chain.EosioContractAbi(abi.AbiDef{})
		abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)
		var testData string = `{
            "account":     "lnkauth.acct",
			"code":        "lnkauth.code",
			"type":        "lnkauth.type",
			"requirement": "lnkauth.rqm"
		}`

		strRead := strings.NewReader(testData)
		linkAuth := chain.LinkAuth{}
		err := json.NewDecoder(strRead).Decode(&linkAuth)
		assert.NoError(t, err)
		assert.Equal(t, "lnkauth.acct", linkAuth.Account.String())
		assert.Equal(t, "lnkauth.code", linkAuth.Code.String())
		assert.Equal(t, "lnkauth.type", linkAuth.Type.String())
		assert.Equal(t, "lnkauth.rqm", linkAuth.Requirement.String())

		data := common.Variants{}
		strRead = strings.NewReader(testData)
		err = json.NewDecoder(strRead).Decode(&data)
		assert.NoError(t, err)
		var2 := verifyByteRoundTripConversion(abis, "linkauth", &data)

		var2Byte, err := json.Marshal(var2)
		var linkAuth2 chain.LinkAuth
		err = json.Unmarshal(var2Byte, &linkAuth2)
		assert.NoError(t, err)

		assert.Equal(t, linkAuth.Account, linkAuth2.Account)
		assert.Equal(t, linkAuth.Code, linkAuth2.Code)
		assert.Equal(t, linkAuth.Type, linkAuth2.Type)
		assert.Equal(t, linkAuth.Requirement, linkAuth2.Requirement)

		//var linkAuth3 chain.LinkAuth
		//verifyTypeRoundTripConversion(linkAuth3, abis, "linkauth", &data)

	}).Catch(func(e interface{}) {
		panic(e)
	}).End()
}

func TestUpdateauth(t *testing.T) {
	Try(func() {
		var testData string = `{
	    	"account" : "updauth.acct",
		    "permission" : "updauth.prm",
		    "parent" : "updauth.prnt",
		    "auth" : {
		        "threshold" : 2147483145,
			    "keys" : [ {"key" : "EOS65rXebLhtk2aTTzP4e9x1AQZs7c5NNXJp89W8R3HyaA6Zyd4im", "weight" : 57005},
		                   {"key" : "EOS5eVr9TVnqwnUBNwf9kwMTbrHvX5aPyyEG97dz2b2TNeqWRzbJf", "weight" : 57605} ],
                "accounts" : [ {"permission" : {"actor" : "prm.acct1", "permission" : "prm.prm1"}, "weight" : 53005 },
	                       {"permission" : {"actor" : "prm.acct2", "permission" : "prm.prm2"}, "weight" : 53405 } ],
                "waits" : []
	        }
	   }`

		strRead := strings.NewReader(testData)
		data := common.Variants{}
		err := json.NewDecoder(strRead).Decode(&data)
		strRead = strings.NewReader(testData)
		updauth := chain.UpdateAuth{}
		err = json.NewDecoder(strRead).Decode(&updauth)
		assert.NoError(t, err)

		assert.Equal(t, "updauth.acct", updauth.Account.String())
		assert.Equal(t, "updauth.prm", updauth.Permission.String())
		assert.Equal(t, "updauth.prnt", updauth.Parent.String())
		assert.Equal(t, uint32(2147483145), updauth.Auth.Threshold)
		assert.Equal(t, 2, len(updauth.Auth.Keys))
		assert.Equal(t, "EOS65rXebLhtk2aTTzP4e9x1AQZs7c5NNXJp89W8R3HyaA6Zyd4im", updauth.Auth.Keys[0].Key.String())
		assert.Equal(t, types.WeightType(57005), updauth.Auth.Keys[0].Weight)
		assert.Equal(t, "EOS5eVr9TVnqwnUBNwf9kwMTbrHvX5aPyyEG97dz2b2TNeqWRzbJf", updauth.Auth.Keys[1].Key.String())
		assert.Equal(t, types.WeightType(57605), updauth.Auth.Keys[1].Weight)
		assert.Equal(t, 2, len(updauth.Auth.Accounts))
		assert.Equal(t, "prm.acct1", updauth.Auth.Accounts[0].Permission.Actor.String())
		assert.Equal(t, "prm.prm1", updauth.Auth.Accounts[0].Permission.Permission.String())
		assert.Equal(t, types.WeightType(53005), updauth.Auth.Accounts[0].Weight)
		assert.Equal(t, "prm.acct2", updauth.Auth.Accounts[1].Permission.Actor.String())
		assert.Equal(t, "prm.prm2", updauth.Auth.Accounts[1].Permission.Permission.String())
		assert.Equal(t, types.WeightType(53405), updauth.Auth.Accounts[1].Weight)

		abiDef := chain.EosioContractAbi(abi.AbiDef{})
		abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)

		var2 := verifyByteRoundTripConversion(abis, "updateauth", &data)

		var2Byte, err := json.Marshal(var2)
		var updateauth2 chain.UpdateAuth
		err = json.Unmarshal(var2Byte, &updateauth2)
		assert.NoError(t, err)
		assert.Equal(t, updateauth2.Account.String(), updauth.Account.String())
		assert.Equal(t, updateauth2.Permission.String(), updauth.Permission.String())
		assert.Equal(t, updateauth2.Parent.String(), updauth.Parent.String())
		assert.Equal(t, updateauth2.Auth.Threshold, updauth.Auth.Threshold)

		assert.Equal(t, len(updateauth2.Auth.Keys), len(updauth.Auth.Keys))
		assert.Equal(t, updateauth2.Auth.Keys[0].Key.String(), updauth.Auth.Keys[0].Key.String())
		assert.Equal(t, updateauth2.Auth.Keys[0].Weight, updauth.Auth.Keys[0].Weight)
		assert.Equal(t, updateauth2.Auth.Keys[1].Key.String(), updauth.Auth.Keys[1].Key.String())
		assert.Equal(t, updateauth2.Auth.Keys[1].Weight, updauth.Auth.Keys[1].Weight)
		assert.Equal(t, len(updateauth2.Auth.Accounts), len(updauth.Auth.Accounts))
		assert.Equal(t, updateauth2.Auth.Accounts[0].Permission.Actor.String(), updauth.Auth.Accounts[0].Permission.Actor.String())
		assert.Equal(t, updateauth2.Auth.Accounts[0].Permission.Permission.String(), updauth.Auth.Accounts[0].Permission.Permission.String())
		assert.Equal(t, updateauth2.Auth.Accounts[0].Weight, updauth.Auth.Accounts[0].Weight)
		assert.Equal(t, updateauth2.Auth.Accounts[1].Permission.Actor.String(), updauth.Auth.Accounts[1].Permission.Actor.String())
		assert.Equal(t, updateauth2.Auth.Accounts[1].Permission.Permission.String(), updauth.Auth.Accounts[1].Permission.Permission.String())
		assert.Equal(t, updateauth2.Auth.Accounts[1].Weight, updauth.Auth.Accounts[1].Weight)

		//var updateAuth3 chain.UpdateAuth
		//verifyTypeRoundTripConversion(updateAuth3, abis, "updateauth", &data)
	}).Catch(func(e interface{}) {
		panic(e)
	}).End()
}

func TestDeleteauth(t *testing.T) {
	Try(func() {
		var testData string = `{
			"account" : "delauth.acct",
			"permission" : "delauth.prm"
		}`

		strRead := strings.NewReader(testData)
		deleteAuth := chain.DeleteAuth{}
		err := json.NewDecoder(strRead).Decode(&deleteAuth)
		assert.NoError(t, err)
		assert.Equal(t, "delauth.acct", deleteAuth.Account.String())
		assert.Equal(t, "delauth.prm", deleteAuth.Permission.String())

		abiDef := chain.EosioContractAbi(abi.AbiDef{})
		abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)

		strRead = strings.NewReader(testData)
		data := common.Variants{}
		err = json.NewDecoder(strRead).Decode(&data)
		assert.NoError(t, err)
		var2 := verifyByteRoundTripConversion(abis, "deleteauth", &data)
		var2Byte, err := json.Marshal(var2)
		var deletAuth2 chain.DeleteAuth
		err = json.Unmarshal(var2Byte, &deletAuth2)
		assert.NoError(t, err)
		assert.Equal(t, deleteAuth.Account, deletAuth2.Account)
		assert.Equal(t, deleteAuth.Permission, deletAuth2.Permission)

		//var deleteAuth3 chain.DeleteAuth
		//verifyTypeRoundTripConversion(&deleteAuth3,abis,"deleteauth",&data)

	}).Catch(func(e interface{}) {
		panic(e)
	}).End()
}

//func TestSetcode(t *testing.T) {
//	Try(func() {
//		var testData string =`{
//			"account" : "setcode.acc",
//			"vmtype" : 0,
//			"vmversion" : 0,
//			"code" : "0061736d0100000001390a60037e7e7f017f60047e7e7f7f017f60017e0060057e7e7e7f7f"
//		}`
//		strRead := strings.NewReader(testData)
//		setCode := chain.SetCode{}
//		err := json.NewDecoder(strRead).Decode(&setCode)
//		assert.NoError(t, err)
//		assert.Equal(t, "setcode.acc", setCode.Account.String())
//		assert.Equal(t, uint8(0), setCode.VmType)
//		assert.Equal(t,uint8(0),setCode.VmVersion)
//		assert.Equal(t,"0061736d0100000001390a60037e7e7f017f60047e7e7f7f017f60017e0060057e7e7e7f7f",string(setCode.Code))
//
//
//		abiDef := chain.EosioContractAbi(abi.AbiDef{})
//		abis := abi.NewAbiSerializer(abiDef, maxSerializationTime)
//
//		strRead = strings.NewReader(testData)
//		data := common.Variants{}
//		err = json.NewDecoder(strRead).Decode(&data)
//		assert.NoError(t, err)
//		var2 := verifyByteRoundTripConversion(abis, "setcode", &data)
//		var2Byte, err := json.Marshal(var2)
//		var setCode2 chain.SetCode
//		err = json.Unmarshal(var2Byte, &setCode2)
//		assert.NoError(t, err)
//		assert.Equal(t, setCode.Account, setCode2.Account)
//		assert.Equal(t, setCode.VmType, setCode2.VmType)
//		assert.Equal(t, setCode.VmVersion, setCode2.VmVersion)
//		assert.Equal(t, setCode.Code, setCode2.Code)
//
//		//var setCode3 chain.SetCode
//		//verifyTypeRoundTripConversion(&setCode3,abis,"setcode",&data)
//	}).Catch(func(e interface{}) {
//	panic(e)
//	}).End()
//}

//Try(func() {
//
//}).Catch(func(e interface{}) {
//panic(e)
//}).End()
