package eosapi

import (
	"context"
	"fmt"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type EOSApi struct {
}

func (api *EOSApi) GetInfo(ctx context.Context) *InfoResp {
	return &InfoResp{}
}

func (api *EOSApi) CreateKey() string {
	prikey, _ := ecc.NewRandomPrivateKey()
	pri := fmt.Sprintf("Private Key: %s\n", prikey.String())
	pub := fmt.Sprintf("Public Key: %s\n", prikey.PublicKey().String())
	out := pri + pub
	return out
}

//func (api *EOSApi)getInfoCli()  error {
//	resp, err := getInfo()
//	if err != nil {
//		return err
//	}
//	display, err := json.Marshal(resp)
//	if err != nil {
//		return err
//	}
//	fmt.Println(string(display))
//	return
//}
//func getInfo() (out *InfoResp, err error) {
//	variant, err := DoHttpCall(chainUrl, getInfoFunc, nil)
//	if err := json.Unmarshal(variant, &out); err != nil {
//		return nil, fmt.Errorf("Unmarshal: %s", err)
//	}
//	return
//}
