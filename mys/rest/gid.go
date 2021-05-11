package rest

type GameType uint8

const (
	NoGame  GameType = 0
	Honkai3 GameType = 1
	Genshin GameType = 2
	Honkai2 GameType = 3
	WeiDing GameType = 4
	DaBieYe GameType = 5
)

var gamesTable = map[GameType]string{
	Honkai3: "崩坏3",
	Genshin: "原神",
	Honkai2: "崩坏学园2",
	WeiDing: "未定事件簿",
	DaBieYe: "大别野",
}

func (gt GameType) Name() string {
	return gamesTable[gt]
}
