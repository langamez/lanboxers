package render

import (
	"lanBox/domain"
	"strings"
)

func drawHud(renderInfo domain.RenderInfo) {
	name(renderInfo, domain.PlayerMain)
	name(renderInfo, domain.PlayerOpponent)

	DrawHealth(renderInfo, domain.PlayerMain)
	DrawHealth(renderInfo, domain.PlayerOpponent)
}

func DrawHealth(hudInfo domain.RenderInfo, playerID domain.PlayerID) {
	var (
		hpChar string
		hpPos  = domain.Position{Y: hudInfo.Cage.Area.Min.Y - 2}
	)
	//Get cage length
	cageLength := hudInfo.Cage.Area.Max.X - hudInfo.Cage.Area.Min.X
	//Calculate each player health bar length
	hpLength := cageLength / 3
	//Calculate space between health bar
	indent := hpLength - (2 * domain.WallLength)
	//Convert to character
	filled := hudInfo.Boxers[playerID].Health * hpLength / 100
	hpChar = strings.Repeat("█", filled) +
		strings.Repeat("░", hpLength-filled)
	//Calculate horizontal health bar position
	if playerID == domain.PlayerMain {
		hpPos.X = hudInfo.Cage.Area.Min.X + domain.WallLength
	} else {
		hpPos.X = hudInfo.Cage.Area.Min.X + domain.WallLength + hpLength + indent
	}
	//Print
	PrintOn(hpPos, hudInfo.Boxers[playerID].Color, hpChar)
}

func name(hudInfo domain.RenderInfo, playerID domain.PlayerID) {
	// todo correct name and hp bar position
	var namePos = domain.Position{Y: hudInfo.Cage.Area.Min.Y - 3}
	//Get cage length
	cageLength := hudInfo.Cage.Area.Max.X - hudInfo.Cage.Area.Min.X
	//Calculate each player health bar length
	hpLength := cageLength / 3
	//Calculate space between health bar
	indent := hpLength - (2 * domain.WallLength)
	//Calculate horizontal health bar position
	if playerID == domain.PlayerMain {
		namePos.X = hudInfo.Cage.Area.Min.X + domain.WallLength
	} else {
		namePos.X = hudInfo.Cage.Area.Min.X + domain.WallLength + hpLength + indent
	}
	//Print
	PrintOn(namePos, hudInfo.Boxers[playerID].Color, hudInfo.Boxers[playerID].Name)
}
