package posts

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

var (
	MissingTextError = errors.New("Text field missing for text type of post")
	MissingUrlError  = errors.New("Url field missing for link type of post")
)

type Author struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
}

type Comment struct {
	Created time.Time `json:"created"`
	Author  Author    `json:"author"`
	Body    string    `json:"body"`
	ID      string    `json:"id"`
}

type Vote struct {
	UserId int `json:"users"`
	Vote   int `json:"vote"`
}

type Post struct {
	sync.RWMutex
	MongoId          primitive.ObjectID `json:"-" bson:"_id"`
	ID               string             `json:"id"`
	Views            int                `json:"views"`
	Type             string             `json:"type"`
	Title            string             `json:"title"`
	Category         string             `json:"category"`
	Text             string             `json:"text"`
	Url              string             `json:"url"`
	Upvotes          int                `json:"-"`
	Downvotes        int                `json:"-"`
	Score            int                `json:"score"`
	UpvotePercentage int                `json:"upvotePercentage"`
	Votes            []*Vote            `json:"votes"`
	Author           Author             `json:"author"`
	Comments         []*Comment         `json:"comments"`
	Created          time.Time          `json:"created"`
}

type PostIn struct {
	Type     string `json:"type" valid:"in(link|text)"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Text     string `json:"text" valid:"optional"`
	Url      string `json:"url" valid:"optional"`
	Author   Author
}

func (in *PostIn) IsValid() error {
	if in.Type == "text" && in.Text == "" {
		return MissingTextError
	}
	if in.Type == "link" && in.Url == "" {
		return MissingUrlError
	}
	return nil
}

type CommentIn struct {
	Comment string `json:"comment"`
	Author  Author
}
