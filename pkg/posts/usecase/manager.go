package usecase

import (
	"context"
	"github.com/google/uuid"
	"golang-stepik-2022q1/reditclone/pkg/errors"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/posts"
	"golang-stepik-2022q1/reditclone/pkg/posts/repo"
	"time"
)

type Repo interface {
	GetAll() ([]*posts.Post, error)
	FilterByUserName(userName string) ([]*posts.Post, error)
	Add(*posts.Post) (*posts.Post, error)
	GetById(string) (*posts.Post, error)
	Delete(postId string) (int64, error)
	AddComment(post *posts.Post, comment *posts.Comment) (int64, error)
	DeleteComment(post *posts.Post, commentId string) (int64, error)
	Vote(post *posts.Post, userId string, score int) (*posts.Post, error)
	DeleteVote(post *posts.Post, userId string) (*posts.Post, error)
	IncViews(post *posts.Post) (*posts.Post, error)
}

type Manager struct {
	repo *repo.MongoRepo
}

func NewManager(repo *repo.MongoRepo) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) GetAll(ctx context.Context) ([]*posts.Post, error) {
	items, err := m.repo.GetAll()
	if err != nil {
		log.Clog(ctx).Error("Cant fetch posts", log.Fields{"error": err.Error()})
	}
	return items, err
}

func (m *Manager) FilterByUser(ctx context.Context, userName string) ([]*posts.Post, error) {
	items, err := m.repo.FilterByUserName(userName)
	if err != nil {
		log.Clog(ctx).Error("Cant fetch posts", log.Fields{"error": err.Error()})
	}
	return items, err
}

func (m *Manager) Create(ctx context.Context, in *posts.PostIn) (*posts.Post, error) {
	post := &posts.Post{
		ID:       uuid.New().String(),
		Type:     in.Type,
		Title:    in.Title,
		Category: in.Category,
		Text:     in.Text,
		Url:      in.Url,
		Author:   in.Author,
		Votes:    []*posts.Vote{},
		Comments: []*posts.Comment{},
		Created:  time.Now(),
	}
	_, err := m.repo.Add(post)
	if err != nil {
		log.Clog(ctx).Error("Repo error during post creation", log.Fields{"error": err.Error()})
		return nil, err
	}
	return post, nil
}

func (m *Manager) Get(ctx context.Context, postId string) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if post == nil {
		log.Clog(ctx).Info("Item not found")
		return nil, ItemNotFound
	}
	_, err = m.repo.IncViews(post)
	if err != nil {
		log.Clog(ctx).Warn("Cant increment post views")
		// here is better to return post without fixed view than return error
	}
	return post, nil
}

func (m *Manager) CreateComment(ctx context.Context, postId string, commentIn *posts.CommentIn) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if err != nil {
		log.Clog(ctx).Error("Repo error during post fetching", log.Fields{"error": err.Error()})
		return nil, errors.InternalError{err.Error()}
	}
	if post == nil {
		log.Clog(ctx).Info("Item not found", log.Fields{"id": postId})
		return post, ItemNotFound
	}

	comment := &posts.Comment{
		Created: time.Now(),
		Author:  commentIn.Author,
		Body:    commentIn.Comment,
		ID:      uuid.New().String(),
	}
	_, err = m.repo.AddComment(post, comment)
	if err != nil {
		return post, errors.InternalError{err.Error()}
	}
	post.Comments = append(post.Comments, comment)
	return post, nil
}

func (m *Manager) DeleteComment(ctx context.Context, postId, commentId string) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if err != nil {
		log.Clog(ctx).Error("Repo error", log.Fields{"error": err.Error()})
		return post, err
	}
	if post == nil {
		log.Clog(ctx).Info("Post not found", log.Fields{"postId": postId})
		return post, ItemNotFound
	}

	_, err = m.repo.DeleteComment(post, commentId)
	if err != nil {
		return post, errors.InternalError{err.Error()}
	}
	return post, nil
}

func (m *Manager) DeletePost(ctx context.Context, postId string) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ItemNotFound
	}
	deletedCount, err := m.repo.Delete(postId)
	if err != nil {
		return nil, err
	}
	if deletedCount == 0 {
		return nil, ItemNotFound
	}
	return post, nil
}

func (m *Manager) Upvote(postId string, userId int) (*posts.Post, error) {
	return m.vote(postId, &posts.Vote{UserId: userId, Vote: 1})
}

func (m *Manager) Downvote(postId string, userId int) (*posts.Post, error) {
	return m.vote(postId, &posts.Vote{UserId: userId, Vote: -1})
}

func (m *Manager) Unvote(postId string, userId int) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return post, ItemNotFound
	}

	var vote *posts.Vote
	for _, v := range post.Votes {
		if v.UserId == userId {
			vote = v
		}
	}
	if vote == nil {
		return post, nil
	}
	// recalculate statistics
	if vote.Vote == 1 {
		post.Upvotes -= 1
		post.Score -= 1
	} else {
		post.Downvotes -= 1
		post.Score += 1
	}
	post.UpvotePercentage = int(float32(post.Upvotes) / float32(post.Upvotes+post.Downvotes) * 100)

	_, err = m.repo.DeleteVote(postId, userId)
	if err != nil {
		return post, errors.InternalError{"Cant update post"}
	}
	m.repo.UpdateStat(postId, post.UpvotePercentage, post.Score)
	return post, nil
}

func (m *Manager) vote(postId string, vote *posts.Vote) (*posts.Post, error) {
	post, err := m.repo.GetById(postId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return post, ItemNotFound
	}

	var existingVote *posts.Vote
	for _, v := range post.Votes {
		if v.UserId == vote.UserId {
			existingVote = v
		}
	}

	// if existing vote is the same -> do nothing
	if existingVote != nil && existingVote.Vote == vote.Vote {
		log.Debug("Vote not changed")
		return post, nil
	}

	var delta int
	if existingVote == nil {
		delta = vote.Vote
		post.Votes = append(post.Votes, vote)
	} else {
		// existing vote is opposite
		// f.e. if we change vote from -1 to 1, we should increment score with 2
		delta = 2
		// we also change existing vote to return proper post data
		existingVote.Vote = vote.Vote
	}
	// remove existing vote from stats, add new vote and recalculate percentage
	if vote.Vote == 1 {
		post.Upvotes += delta
	} else {
		post.Downvotes += delta
	}
	post.Score = post.Upvotes - post.Downvotes
	post.UpvotePercentage = int(float32(post.Upvotes) / float32(post.Upvotes+post.Downvotes) * 100)

	// now record/update vote and stats
	_, err = m.repo.Vote(postId, vote)
	if err != nil {
		return post, errors.InternalError{"Cant update post"}
	}
	m.repo.UpdateStat(postId, post.UpvotePercentage, post.Score)
	return post, nil
}
