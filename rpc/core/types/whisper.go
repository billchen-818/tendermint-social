package coretypes

import (
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/whisper"
)

type Envelope struct {
	TTL   uint64            `json:"ttl"`
	Topic whisper.TopicType `json:"topic"`
	Data  []byte            `json:"data"`
}

type ResultEnvelope struct {
	Hash      bytes.HexBytes `json:"hash"`
	Timestamp uint64         `json:"timestamp"`
	Data      *Envelope      `json:"data"`
}
