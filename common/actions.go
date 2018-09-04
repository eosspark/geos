package common

// Action
type Action struct {
	Account       AccountName       `json:"account"`
	Name          ActionName        `json:"name"`
	Authorization []PermissionLevel `json:"authorization,omitempty"`
	ActionData
}

type ActionData struct {
	HexData  HexBytes    `json:"hex_data,omitempty"`
	Data     interface{} `json:"data,omitempty" eos:"-"`
	abi      []byte      // TBD: we could use the ABI to decode in obj
	toServer bool
}

func NewActionData(obj interface{}) ActionData {
	return ActionData{
		HexData:  []byte{},
		Data:     obj,
		toServer: true,
	}
}