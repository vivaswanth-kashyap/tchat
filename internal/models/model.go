package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Local SQLite primary key (uint)
	ServerID   string `json:"id" gorm:"column:server_id;unique"` // Server ID reference
	Username   string `json:"username" gorm:"not null"`
	Email      string `json:"email" gorm:"not null"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Message struct {
	gorm.Model
	ServerID  string  `json:"id" gorm:"column:server_id;unique"`
	Content   string  `json:"content" gorm:"not null"`
	UserID    uint    `json:"-"` // Local foreign key
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	ChannelID uint    `json:"-"` // Local foreign key
	Channel   Channel `json:"channel" gorm:"foreignKey:ChannelID"`
}

type Channel struct {
	gorm.Model
	ServerID    string    `json:"id" gorm:"column:server_id;unique"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Messages    []Message `json:"messages" gorm:"foreignKey:ChannelID"`
	Users       []User    `json:"users" gorm:"many2many:channel_users;"`
}
