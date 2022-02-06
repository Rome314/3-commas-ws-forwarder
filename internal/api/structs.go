package api

import (
	"encoding/json"
)

type User struct {
	ApiKey    string `json:"api_key"`
	Signature string `json:"signature"`
}

type Identifier struct {
	Channel string `json:"channel"`
	Users   []User `json:"users"`
}

func (i *Identifier) String() string {
	bts, _ := json.Marshal(i)
	return string(bts)
}

type Command struct {
	Command    string `json:"command"`
	Identifier string `json:"identifier"`
}
