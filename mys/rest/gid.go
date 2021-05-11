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

type gidDetail struct {
	Name, URL string
}

var gamesTable = map[GameType]gidDetail{
	Honkai3: {"崩坏3", "bh3"},
	Genshin: {"原神", "ys"},
	Honkai2: {"崩坏学园2", "bh2"},
	WeiDing: {"未定事件簿", "wd"},
	DaBieYe: {"大别野", "dby"},
	NoGame:  {"未知", ""},
}

func (gt GameType) Name() string {
	return gamesTable[gt].Name
}

var gidReverseTable = map[string]GameType{
	"":    NoGame,
	"bh3": Honkai3,
	"ys":  Genshin,
	"bh2": Honkai2,
	"wd":  WeiDing,
	"dby": DaBieYe,
}

func FromGidStr(s string) GameType {
	return gidReverseTable[s]
}
