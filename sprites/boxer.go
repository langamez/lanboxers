package sprites

import "github.com/langamez/lanboxers/domain"

var BodyChars = map[domain.Situation]map[domain.BodyPart]map[domain.Direction]map[domain.Direction]string{
	domain.Idle: {
		domain.Head: {
			domain.Left:  {domain.Left: "(:/)"},
			domain.Right: {domain.Left: "(/:)"},
		},
		domain.Arm: {
			domain.Left:  {domain.Left: "╭╭══O", domain.Right: "╰╰══O"},
			domain.Right: {domain.Right: "O══╮╮", domain.Left: "O══╯╯"},
		},
		domain.Shoulder: {
			domain.Left:  {domain.Left: "\\\\", domain.Right: "//"},
			domain.Right: {domain.Right: "//", domain.Left: "\\\\"},
		},
	},
	domain.Punch: {
		domain.Shoulder: {
			domain.Left:  {domain.Left: "\\\\══O", domain.Right: "//══O"},
			domain.Right: {domain.Right: "O══//", domain.Left: "O══\\\\"},
		},
	},
	domain.PunchInit: {
		domain.Shoulder: {
			domain.Left:  {domain.Left: "\\\\═O", domain.Right: "//═O"},
			domain.Right: {domain.Right: "O═//", domain.Left: "O═\\\\"},
		},
	},
	domain.HeadHit: {
		domain.Head: {
			domain.Left:  {domain.Left: "**"},
			domain.Right: {domain.Right: "**"},
		},
		domain.Arm: {
			domain.Left:  {domain.Left: "* ", domain.Right: "* "},
			domain.Right: {domain.Right: " *", domain.Left: " *"},
		},
		domain.Shoulder: {
			domain.Left:  {domain.Left: "*", domain.Right: "*"},
			domain.Right: {domain.Right: "*", domain.Left: "*"},
		},
	},
	domain.HeadInitHit: {
		domain.Head: {
			domain.Left:  {domain.Left: "*"},
			domain.Right: {domain.Right: "*"},
		},
		domain.Shoulder: {
			domain.Left:  {domain.Left: "*", domain.Right: "*"},
			domain.Right: {domain.Right: "*", domain.Left: "*"},
		},
	},
	domain.ShoulderHit: {
		domain.Head: {
			domain.Left:  {domain.Left: "* "},
			domain.Right: {domain.Left: " *"},
		},
		domain.Arm: {
			domain.Left:  {domain.Left: "* ", domain.Right: "* "},
			domain.Right: {domain.Right: " *", domain.Left: " *"},
		},
		domain.Shoulder: {
			domain.Left:  {domain.Left: "**", domain.Right: "**"},
			domain.Right: {domain.Right: "**", domain.Left: "**"},
		},
	},
	domain.ShoulderInitHit: {
		domain.Shoulder: {
			domain.Left:  {domain.Left: "*", domain.Right: "*"},
			domain.Right: {domain.Right: "*", domain.Left: "*"},
		},
	},
}
