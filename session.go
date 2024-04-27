package authz

import (
	"time"
)

type Model[ID any] struct {
	ID       ID
	Metadata map[string]string
}

type Session[ID any] struct {
	*Model[ID]
	Tokens []tokenItem
}

func (s *Session[ID]) mergeMetadata(metadata map[string]string) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}

	for k, v := range metadata {
		s.Metadata[k] = v
	}
}

func (s *Session[ID]) addToken(token string, plat string) {
	s.Tokens = append(s.Tokens, tokenItem{
		Value:     token,
		Platform:  plat,
		CreatedAt: time.Now(),
	})
}

type tokenItem struct {
	Value     string
	Platform  string
	CreatedAt time.Time
}
