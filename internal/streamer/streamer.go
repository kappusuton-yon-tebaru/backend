package streamer

import (
	"bufio"
	"io"
	"sync"

	"github.com/google/uuid"
)

type Streamer struct {
	sync.Mutex
	subs map[uuid.UUID]chan string
}

func New() *Streamer {
	return &Streamer{
		subs: make(map[uuid.UUID]chan string),
	}
}

func (s *Streamer) StreamFromReader(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		s.Lock()
		for _, ch := range s.subs {
			ch <- line
		}
		s.Unlock()
	}

	s.Stop()
}

func (s *Streamer) Listen() (<-chan string, func()) {
	id := uuid.New()

	s.Lock()
	s.subs[id] = make(chan string)
	s.Unlock()

	unsub := func() {
		s.Lock()
		s.deleteSubscriber(id)
		s.Unlock()
	}

	return s.subs[id], unsub
}

func (s *Streamer) Stop() {
	s.Lock()
	for id := range s.subs {
		s.deleteSubscriber(id)
	}
	s.Unlock()
}

func (s *Streamer) deleteSubscriber(id uuid.UUID) {
	ch, ok := s.subs[id]
	if ok {
		delete(s.subs, id)
		close(ch)
	}
}
