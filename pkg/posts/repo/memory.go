package repo

import (
	"errors"
	"golang-stepik-2022q1/reditclone/pkg/posts"
	"sync"
)

type MemRepo struct {
	sync.RWMutex
	data []*posts.Post
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		data: make([]*posts.Post, 0, 10),
	}
}

func (repo *MemRepo) GetAll() ([]*posts.Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	return repo.data, nil
}

func (repo *MemRepo) FilterByUserName(userName string) ([]*posts.Post, error) {
	repo.RLock()
	defer repo.RUnlock()

	userPosts := make([]*posts.Post, 0, 10)
	for _, post := range repo.data {
		if post.Author.Username == userName {
			userPosts = append(userPosts, post)
		}
	}
	return userPosts, nil
}

func (repo *MemRepo) Add(item *posts.Post) (*posts.Post, error) {
	repo.Lock()
	defer repo.Unlock()

	repo.data = append(repo.data, item)
	return item, nil
}

func (repo *MemRepo) GetById(id string) (*posts.Post, error) {
	repo.RLock()
	defer repo.RUnlock()
	return repo.getById(id), nil
}

func (repo *MemRepo) getById(id string) *posts.Post {
	for _, post := range repo.data {
		if post.ID == id {
			return post
		}
	}
	return nil
}

func (repo *MemRepo) Update(post *posts.Post) (*posts.Post, error) {
	// as soon as we store and share posts as struct reference,
	// there is no need to replace item in memory storage after it was updated.
	// so here we only check that item exists
	for _, item := range repo.data {
		if item == post {
			return item, nil
		}
	}
	return nil, errors.New("item not found")
}

func (repo *MemRepo) Delete(postId string) (*posts.Post, error) {
	repo.Lock()
	defer repo.Unlock()

	for idx, post := range repo.data {
		if post.ID == postId {
			repo.data[idx] = repo.data[len(repo.data)-1]
			repo.data = repo.data[:len(repo.data)-1]
			return post, nil
		}
	}
	return nil, nil
}

func (repo *MemRepo) AddComment(post *posts.Post, comment *posts.Comment) (*posts.Post, error) {
	post.Lock()
	defer post.Unlock()
	post.Comments = append(post.Comments, comment)
	return post, nil
}

func (repo *MemRepo) DeleteComment(post *posts.Post, commentId string) (*posts.Post, error) {
	post.Lock()
	defer post.Unlock()
	for idx, comment := range post.Comments {
		if comment.ID == commentId {
			post.Comments[idx] = post.Comments[len(post.Comments)-1]
			post.Comments = post.Comments[:len(post.Comments)-1]
			break
		}
	}
	return post, nil
}

func (repo *MemRepo) Vote(post *posts.Post, userId int, score int) (*posts.Post, error) {
	post.Lock()
	defer post.Unlock()

	exists := false
	for _, vote := range post.Votes {
		if vote.UserId == userId {
			vote.Vote = score
			calcStats(post, score*2)
			exists = true
		}
	}
	if !exists {
		post.Votes = append(post.Votes, &posts.Vote{userId, score})
		calcStats(post, score)
	}
	return post, nil
}

func (repo *MemRepo) DeleteVote(post *posts.Post, userId int) (*posts.Post, error) {
	post.Lock()
	defer post.Unlock()

	for idx, vote := range post.Votes {
		if vote.UserId == userId {
			post.Votes[idx] = post.Votes[len(post.Votes)-1]
			post.Votes = post.Votes[:len(post.Votes)-1]

			calcStats(post, -vote.Vote)
			break
		}
	}
	return post, nil
}

func (repo *MemRepo) IncViews(post *posts.Post) (*posts.Post, error) {
	post.Lock()
	defer post.Unlock()
	post.Views++
	return post, nil
}

func calcStats(post *posts.Post, vote int) {
	post.Score += vote
	post.UpvotePercentage = int(float32(post.Score) / float32(len(post.Votes)) * 100)
}
