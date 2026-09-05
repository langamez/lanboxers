package game

import "github.com/langamez/lanboxers/domain"

func NewRenderChannels() map[domain.PlayerID]chan domain.RenderCommand {
	return map[domain.PlayerID]chan domain.RenderCommand{
		domain.PlayerMain:     make(chan domain.RenderCommand),
		domain.PlayerOpponent: make(chan domain.RenderCommand),
	}
}
