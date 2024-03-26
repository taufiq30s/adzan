package calc

import "github.com/taufiq30s/adzan/internal/utils"

type Mazhab float32

const (
	SYAFI Mazhab = iota + 1
	HANAFI
)

var MazhabToShadowLengthMap = map[Mazhab]utils.ShadowLength{
	SYAFI:  utils.SINGLE,
	HANAFI: utils.DOUBLE,
}
