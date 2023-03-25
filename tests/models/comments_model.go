package models

import "time"

type CommentsModel struct {
	User    UsersModel     `dbKit:"column:user_id, foreignKey:id"`
	Id      uint           `dbKit:"column:id, primaryKey"`
	Post    PostsModel     `dbKit:"column:post_id, foreignKey:id"`
	Parent  *CommentsModel `dbKit:"column:parent_id, foreignKey:id"`
	Content string         `dbKit:"column:content"`
	Created time.Time      `dbKit:"column:created_at"`
	Updated time.Time      `dbKit:"column:updated_at"`
}

func (s *CommentsModel) DatabaseName() string {
	return "acceptance"
}

func (s *CommentsModel) TableName() string {
	return "comments"
}
