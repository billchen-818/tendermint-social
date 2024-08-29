package whisper

import (
	"encoding/json"

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

func (e *Envelope) Hash() bytes.HexBytes {
	if e.hash == nil {
		encoded, _ := json.Marshal(e)
		e.hash = tmhash.Sum(encoded)
	}

	return e.hash
}
