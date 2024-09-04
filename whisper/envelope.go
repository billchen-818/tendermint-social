package whisper

import (
	"encoding/json"
	"time"

	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/bytes"
)

type Envelope struct {
	Expiry uint64
	TTL    uint64
	Topic  TopicType
	Data   []byte

	hash bytes.HexBytes
}

func NewEnvelope(ttl uint64, topic TopicType, msg *sentMessage) *Envelope {
	env := Envelope{
		Expiry: uint64(time.Now().Add(time.Duration(ttl)).Unix()),
		TTL:    ttl,
		Topic:  topic,
		Data:   msg.Raw,
	}

	return &env
}

func (e *Envelope) Hash() bytes.HexBytes {
	if e.hash == nil {
		encoded, _ := json.Marshal(e)
		e.hash = tmhash.Sum(encoded)
	}

	return e.hash
}
