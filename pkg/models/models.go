package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDuplicateUsername  = errors.New("models: duplicate login")
)

type Resp struct {
	ErrResp ErrResp
	Posts   []Post
	User    User
	Notify 	[]Notify
	NewNotifyCount int
	Comment Comment
}

type ErrResp struct {
	IsAuthenticated      bool
	IsInvalidCredentials bool
	IsDuplicateEmail     bool
	IsDuplicateUsername  bool
	IsNoRecord           bool
	IsWeakPassword       bool
	IsDifferentPasswords bool
	IsAuthor 			 bool
}

type User struct {
	Login     string `json:"login,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  []byte `json:"password,omitempty"`
	CreatedAt string `json:"created_at"`
}

type Notify struct {
	Id          int
	AuthorId    string
	PostId      string
	Post        Post
	ActionType  string
	IsCommented bool
	IsLiked     bool
	IsDisliked  bool
	UserLogin   string
	UserId      string
	IsActive    bool
}

type Post struct {
	Id         int       `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	UserId     string    `json:"user_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	HumanDate  string
	ImageURL   string `json:"image_url,omitempty"`
	IsLiked    bool
	IsDisliked bool
	Rating     int       `json:"rating,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
	Category   []string  `json:"categories,omitempty"`
	UpdatedCategory string
	IsAuthor   bool
}

type Comment struct {
	Id         int    `json:"id,omitempty"`
	PostId     string `json:"post_id,omitempty"`
	UserId     string `json:"user_id,omitempty"`
	Username   string
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	HumanDate  string
	IsLiked    bool
	IsDisliked bool
	Rating     int `json:"rating,omitempty"`
	IsAuthor   bool
}

type RatingPost struct {
	PostId int
	UserId int
	Value  int
}

type RatingComment struct {
	CommentId int
	UserId    int
	Value     int
}

type CategoryPostLink struct {
	PostId     int
	CategoryId int
}
