package models

type DebugModel struct {
}

func (s DebugModel) DatabaseName() string {
	return "acceptance"
}

func (s DebugModel) TableName() string {
	return "extra"
}
