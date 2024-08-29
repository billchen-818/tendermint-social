package whisper

import (
	"runtime"
	"time"

	"github.com/tendermint/tendermint/libs/log"
)

type Whisper struct {
	quit chan struct{}

	logger log.Logger
}

func New(l log.Logger) *Whisper {
	w := &Whisper{
		quit:   make(chan struct{}),
		logger: l,
	}

	return w
}

func (w *Whisper) Start() error {
	w.logger.Info("starting whisper service")

	numCPU := runtime.NumCPU()
	for i := 0; i < numCPU; i++ {
		go w.processQueue()
	}

	return nil
}

func (w *Whisper) processQueue() {
	for {
		w.logger.Info("process whisper message")
		time.Sleep(time.Second * 10)
	}
}
