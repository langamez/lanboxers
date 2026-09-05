package domain

type RenderInfo struct {
	Cage   Cage
	Boxers Boxers
}

type RenderCommand struct {
	Boxer BaseBoxer
	Parts map[BodyPart][]Direction
}

type RenderChannels map[PlayerID]chan RenderCommand
