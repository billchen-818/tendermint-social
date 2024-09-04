package whisper

import (
	crand "crypto/rand"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type MessageParam struct {
	TTL     uint64
	Src     crypto.PrivKey
	Dst     crypto.PubKey
	Topic   TopicType
	Payload []byte
}

type sentMessage struct {
	Raw []byte
}

func NewSentMessage(param *MessageParam) (*sentMessage, error) {
	const payloadSizeFieldMaxSize = 4
	msg := sentMessage{}
	msg.Raw = make([]byte, 1, flagsLength+payloadSizeFieldMaxSize+len(param.Payload)+signatureLength+PubKeyLength)
	msg.Raw[0] = 0
	msg.addPayloadSizeField(param.Payload)
	msg.Raw = append(msg.Raw, param.Payload...)

	return &msg, nil
}

func (msg *sentMessage) Wrap(option *MessageParam) (*Envelope, error) {
	if option.TTL == 0 {
		option.TTL = DefaultTTL
	}
	if option.Src != nil {
		if err := msg.sign(option.Src); err != nil {
			return nil, err
		}
	}
	if option.Dst != nil {
		if err := msg.enCrypto(option.Dst); err != nil {
			return nil, err
		}
	}

	env := NewEnvelope(option.TTL, option.Topic, msg)
	return env, nil
}

func (msg *sentMessage) sign(privkey crypto.PrivKey) error {
	if isSigned(msg.Raw[0]) {
		return nil
	}

	msg.Raw[0] |= signatureFlag
	signature, err := privkey.Sign(msg.Raw)
	if err != nil {
		return err
	}

	msg.Raw = append(msg.Raw, signature...)
	pubkey := privkey.PubKey().(secp256k1.PubKey)
	msg.Raw = append(msg.Raw, pubkey[:]...)

	return nil
}

func (msg *sentMessage) enCrypto(pubkey crypto.PubKey) error {
	encrypted, err := ecies.Encrypt(crand.Reader, ecies.ImportECDSAPublic(pubkey.(secp256k1.PubKey).ToECDSA()), msg.Raw, nil, nil)
	if err != nil {
		return err
	}

	msg.Raw = encrypted

	return nil
}

func (msg *sentMessage) addPayloadSizeField(payload []byte) {
	fieldSize := getSizeOfPayloadSizeField(payload)
	field := make([]byte, 4)
	binary.LittleEndian.PutUint32(field, uint32(len(payload)))
	field = field[:fieldSize]
	msg.Raw = append(msg.Raw, field...)
	msg.Raw[0] |= byte(fieldSize)
}

func getSizeOfPayloadSizeField(payload []byte) int {
	s := 1
	for i := len(payload); i >= 256; i /= 256 {
		s++
	}
	return s
}

func isSigned(flag byte) bool {
	return (flag & signatureFlag) != 0
}
