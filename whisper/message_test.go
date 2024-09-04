package whisper

import (
	"fmt"
	"testing"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestSentMessage(t *testing.T) {
	param := MessageParam{
		TTL:     50,
		Topic:   BytesToTopic([]byte("12345678")),
		Payload: []byte("hello"),
	}

	src := secp256k1.GenPrivKey()
	dst := secp256k1.GenPrivKey().PubKey()
	param.Src = src
	param.Dst = dst

	sentMsg, err := NewSentMessage(&param)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(len(sentMsg.Raw))
		fmt.Println(cap(sentMsg.Raw))
		fmt.Println(sentMsg.Raw)
		fmt.Println([]byte("hello"))

	}
}
