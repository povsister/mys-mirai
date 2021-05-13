package common

import "github.com/povsister/mys-mirai/mys/runtime"

type Rectangle struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

// represents a file size in byte
type SizeByte struct {
	Size runtime.Int `json:"size"`
}

type Crop struct {
	URL    string `json:"url"`
	Height uint32 `json:"h"`
	Width  uint32 `json:"w"`
	AxisX  int    `json:"x"`
	AxisY  int    `json:"y"`
}
