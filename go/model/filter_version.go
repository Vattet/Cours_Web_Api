package model

import "fmt"

type FilterVersion int

const (
	FilterV1 FilterVersion = iota + 1
	FilterV2
	FilterV3
)

var FilterVersions = []FilterVersion{FilterV1, FilterV2, FilterV3}

func (v FilterVersion) String() string {
	return fmt.Sprintf("v%d", int(v))
}
