package meta

// 论坛板块
type Forum uint32

const (
	YsTavern   Forum = 26
	YsGuide    Forum = 43
	YsDoujin   Forum = 29
	YsOfficial Forum = 28
	YsCoser    Forum = 49
	YsHardcore Forum = 50
)

var forumsTable = map[Forum]string{
	YsTavern:   "酒馆",
	YsGuide:    "攻略",
	YsDoujin:   "同人图",
	YsOfficial: "官方",
	YsCoser:    "COS",
	YsHardcore: "硬核",
}

func (ft Forum) Name() string {
	return forumsTable[ft]
}

// 话题
type Topic uint32

const (
	YsWaterSlime Topic = 180
)

var topicsTable = map[Topic]string{
	YsWaterSlime: "水史莱姆聚集地",
}

func (tt Topic) Name() string {
	return topicsTable[tt]
}
