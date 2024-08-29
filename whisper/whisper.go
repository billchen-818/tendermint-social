package whisper

import (
	"fmt"
	"runtime"
	"time"

	"github.com/tendermint/tendermint/libs/log"
)

type Whisper struct {
	envelopQueue chan *Envelope
	quit         chan struct{}

	logger log.Logger
}

func New(l log.Logger) *Whisper {
	w := &Whisper{
		envelopQueue: make(chan *Envelope, envelopeQueueSize),
		quit:         make(chan struct{}),
		logger:       l,
	}

	return w
}

func (w *Whisper) Send(e *Envelope) error {
	ok, err := w.add(e)
	if err == nil && !ok {
		return fmt.Errorf("failed to add envelope")
	}

	return err
}

func (w *Whisper) Start() error {
	w.logger.Info("starting whisper service")

	numCPU := runtime.NumCPU()
	for i := 0; i < numCPU; i++ {
		go w.processQueue()
	}

	return nil
}

func (w *Whisper) Stop() error {
	close(w.quit)
	w.logger.Info("closing whisper service")
	return nil
}

func (w *Whisper) add(e *Envelope) (bool, error) {
	now := uint64(time.Now().Unix())
	sent := e.Expiry - e.TTL
	if sent > now {
		return false, fmt.Errorf("envelope created in the future [%v]", e.Hash())
	}

	w.envelopQueue <- e
	return true, nil
}

func (w *Whisper) processQueue() {
	for {
		select {
		case <-w.quit:
			return
		case e := <-w.envelopQueue:
			fmt.Printf("received envelope: %v\n", e)
		}
	}
}
