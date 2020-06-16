package cli

import (
	"fmt"
	"sync"
	"time"
)

var spinnerFrames = []string{
	"⠈⠁",
	"⠈⠑",
	"⠈⠱",
	"⠈⡱",
	"⢀⡱",
	"⢄⡱",
	"⢄⡱",
	"⢆⡱",
	"⢎⡱",
	"⢎⡰",
	"⢎⡠",
	"⢎⡀",
	"⢎⠁",
	"⠎⠁",
	"⠊⠁",
}

type Spinner struct {
	stop        chan struct{}
	stopped     chan struct{}
	mu          *sync.Mutex
	running     bool
	ticker      *time.Ticker
	prefix      string
	suffix      string
	frameFormat string
}

func NewSpinner() *Spinner {
	frameFormat := "\x1b[?7l\r%s%s%s\x1b[?7h"
	// toggling wrapping seems to behave poorly on windows
	// in general only the simplest escape codes behave well at the moment,
	// and only in newer shells
	return &Spinner{
		stop:        make(chan struct{}, 1),
		stopped:     make(chan struct{}),
		mu:          &sync.Mutex{},
		frameFormat: frameFormat,
	}
}

func (s *Spinner) SetPrefix(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = prefix
}

func (s *Spinner) SetSuffix(suffix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.suffix = suffix
}

func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return
	}
	s.running = true
	s.ticker = time.NewTicker(time.Millisecond * 1000)

	go func() {
		for {
			for _, frame := range spinnerFrames {
				select {
				case <-s.stop:
					func() {
						s.mu.Lock()
						defer s.mu.Unlock()
						s.ticker.Stop()
						s.running = false
						s.stopped <- struct{}{}
					}()
					return
				case <-s.ticker.C:
					func() {
						s.mu.Lock()
						defer s.mu.Unlock()
						fmt.Printf(s.frameFormat, s.prefix, frame, s.suffix)
					}()

				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.stop <- struct{}{}
	s.mu.Unlock()
	<-s.stopped
}
