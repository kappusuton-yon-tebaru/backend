package hub

import (
	"github.com/kappusuton-yon-tebaru/backend/internal/streamer"
)

type Hub struct {
	streamers map[string]*streamer.Streamer
}

func New() *Hub {
	return &Hub{
		streamers: make(map[string]*streamer.Streamer),
	}
}

func (h *Hub) RegisterStreamer(id string, s *streamer.Streamer) {
	h.streamers[id] = s
}

func (h *Hub) UnregisterStreamer(id string) {
	if _, ok := h.streamers[id]; ok {
		delete(h.streamers, id)
	}
}

func (h *Hub) GetOrRegisterStreamer(id string) *streamer.Streamer {
	s, ok := h.streamers[id]
	if !ok {
		return streamer.New()
	}

	return s
}
