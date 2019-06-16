package hh

import (
	"encoding/json"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/keeper"
	"time"
)


var (
	_ ibc.Packet = &hhIBCPacket{}
)

type hhIBCPacket struct {
	NFT
}

func (this *hhIBCPacket) Timeout() uint64 {

	return uint64(time.Hour.Nanoseconds())
}

func (this *hhIBCPacket) Commit() []byte {
	b, err:=json.Marshal(this)
	if err!=nil {
		panic(err)
	}
	return b
}
