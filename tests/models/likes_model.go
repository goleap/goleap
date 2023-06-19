package models

import "time"

type LikeModel struct {
	Id        uint       `dbKit:"column:id, primaryKey"`
	UserId    uint       `dbKit:"column:user_id"`
	User      UsersModel `dbKit:"column:user_id, foreignKey:id"`
	Post      PostsModel `dbKit:"column:post_id, foreignKey:id"`
	CreatedAt time.Time  `dbKit:"column:created_at"`
}

func (model *LikeModel) DatabaseName() string {
	return "acceptance"
}

func (model *LikeModel) TableName() string {
	return "likes"
}

func (model *LikeModel) ConnectorName() string {
	return "acceptance_extend"
}
