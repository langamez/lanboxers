package render

import (
	"github.com/langamez/lanboxers/domain"
	"github.com/langamez/lanboxers/sprites"
)

func BoxerRenderer(ch <-chan domain.RenderCommand, previous domain.BaseBoxer) {
	for cmd := range ch {
		// draw cmd.Boxer using cmd.Parts
		color := "\033[31m"
		// Hit effect is beside the boxer body parts
		// so it will not remove the parts if it was a hit
		if cmd.Boxer.Situation.SituationType < domain.HeadHit {
			//if sit < HeadHit {
			// Not getting hit
			color = cmd.Boxer.Color
			ClearBoxer(previous, cmd.Parts)
		}
		DrawBoxer(cmd.Boxer, cmd.Parts, color)
		Frame(1)
		previous = cmd.Boxer
	}
}

func LoseEffect(
	color string,
	position domain.Position,
) {
	text := sprites.LoseText
	PrintOn(domain.Position{X: (position.X - len(text)) / 2, Y: position.Y / 2}, color, text)
	Frame(100)
}
