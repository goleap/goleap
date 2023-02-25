package testmodels

import (
	"time"
)

type BaseModel struct {
	Id        uint        `goleap:"column:id, primaryKey"`
	Recursive *BaseModel  `goleap:"column:recursive_id, foreignKey:id"`
	Slice     []BaseModel `goleap:"column:slice_id, foreignKey:id"`
	Extra     ExtraModel  `goleap:"column:extra_id, foreignKey:id"`
	CreatedAt time.Time   `goleap:","`

	BadTag string `goleap:"column:bad_tag:bad_tag"`

	// private bool
}

func (s *BaseModel) DatabaseName() string {
	return "test"
}

func (s *BaseModel) TableName() string {
	return "base"
}
