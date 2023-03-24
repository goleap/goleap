package models

type ExtraModel struct {
}

func (s ExtraModel) DatabaseName() string {
	return "acceptance"
}

func (s ExtraModel) TableName() string {
	return "extra"
}
