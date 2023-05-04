package models

import "time"

type PostsModel struct {
	Id       uint            `dbKit:"column:id, primaryKey"`
	Creator  UsersModel      `dbKit:"column:c_user_id, foreignKey:id"`
	Editor   UsersModel      `dbKit:"column:u_user_id, foreignKey:id"`
	Comments []CommentsModel `dbKit:"column:id, foreignKey:post_id"`
	Title    string          `dbKit:"column:title"`
	Content  string          `dbKit:"column:content"`
	Created  time.Time       `dbKit:"column:created_at"`
	Updated  time.Time       `dbKit:"column:updated_at"`
}

func (s *PostsModel) DatabaseName() string {
	return "acceptance"
}

func (s *PostsModel) TableName() string {
	return "posts"
}

func (s *PostsModel) ConnectorName() string {
	return "acceptance"
}
