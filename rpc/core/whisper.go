package core

import (
	"errors"
	"time"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	"github.com/tendermint/tendermint/whisper"
)

var (
	ErrNoTopics = errors.New("missing topics")
)

func PublishEnvelope(ctx *rpctypes.Context, envelop coretypes.Envelope) (*coretypes.ResultEnvelope, error) {
	if envelop.Topic == (whisper.TopicType{}) {
		return nil, ErrNoTopics
	}

	if envelop.TTL == 0 {
		envelop.TTL = whisper.DefaultTTL
	}

	e := &whisper.Envelope{
		Expiry: uint64(time.Now().Unix()) + envelop.TTL,
		TTL:    envelop.TTL,
		Topic:  envelop.Topic,
		Data:   envelop.Data,
	}

	var result = coretypes.ResultEnvelope{}

	err := env.Whisper.Send(e)
	if err == nil {
		result.Hash = e.Hash()
		result.Timestamp = e.Expiry - e.TTL
		result.Data = nil
	}

	return &result, err
}
