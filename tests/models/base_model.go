package models

import (
	"time"
)

type BaseModel struct {
	Id        uint        `dbKit:"column:id, primaryKey"`
	Recursive *BaseModel  `dbKit:"column:recursive_id, foreignKey:id"`
	Slice     []BaseModel `dbKit:"column:slice_id, foreignKey:id"`
	Extra     ExtraModel  `dbKit:"column:extra_id, foreignKey:id"`
	CreatedAt time.Time   `dbKit:","`

	BadTag string `dbKit:"column:bad_tag:bad_tag"`

	private bool
}

func (s *BaseModel) DatabaseName() string {
	// Only for staticcheck
	s.private = true

	return "test"
}

func (s *BaseModel) TableName() string {
	return "base"
}
