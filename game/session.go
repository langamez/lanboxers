package game

import (
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/internal/terminal"
	"github.com/langamez/lanboxers/network"
)

type Session struct {
	Host bool
	Game *Game
	Conn *network.Connection
}

func NewSession(g *Game) (*Session, error) {
	isHost, ip := terminal.GetConnectionInfo()
	conn, err := network.SetupConnection(isHost, ip, 7777)
	if err != nil {
		return nil, err
	}

	//if !isHost {
	//	g.Boxers[domain.PlayerMain].Direction = g.Boxers[domain.PlayerMain].Direction.Opposite()
	//	g.Boxers[domain.PlayerOpponent].Direction = g.Boxers[domain.PlayerOpponent].Direction.Opposite()
	//}

	return &Session{
		Game: g,
		Conn: conn,
		Host: isHost,
	}, nil
}

func (s *Session) Start() {
	var player domain.PlayerID
	localChan := make(chan domain.Act)
	remoteChan := make(chan domain.Act)

	if s.Host {
		player = domain.PlayerMain
	} else {
		player = domain.PlayerOpponent
	}
	go terminal.ReadActions(localChan)
	go s.handleLocalInput(player, localChan)

	go s.Conn.Listen(remoteChan)
	go s.handleRemoteInput(player.Opposite(), remoteChan)
}

func (s *Session) handleLocalInput(
	player domain.PlayerID,
	actions <-chan domain.Act,
) {
	for act := range actions {
		if err := s.Conn.SendAct(act); err != nil {
			s.Game.CloseGame()
		}

		s.Game.HandleAction(act, player)
	}
}

func (s *Session) handleRemoteInput(
	player domain.PlayerID,
	actions <-chan domain.Act,
) {
	for act := range actions {
		s.Game.HandleAction(act, player)
	}
}
