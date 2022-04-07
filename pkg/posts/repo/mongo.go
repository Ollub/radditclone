package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-stepik-2022q1/reditclone/pkg/posts"
)

const (
	PostsDb         = "reddit"
	PostsCollection = "posts"
)

type MongoRepo struct {
	coll *mongo.Collection
}

func NewMongoRepo(mdb *mongo.Client) *MongoRepo {
	return &MongoRepo{mdb.Database(PostsDb).Collection(PostsCollection)}
}

func (repo *MongoRepo) GetAll() ([]*posts.Post, error) {
	items := make([]*posts.Post, 0, 10)
	res, err := repo.coll.Find(context.Background(), bson.M{})
	err = res.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *MongoRepo) FilterByUserName(userName string) ([]*posts.Post, error) {
	items := make([]*posts.Post, 0, 10)
	res, err := repo.coll.Find(context.Background(), bson.M{"author.username": userName})
	err = res.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *MongoRepo) Add(item *posts.Post) (*posts.Post, error) {
	item.MongoId = primitive.NewObjectID()
	item.ID = item.MongoId.Hex()

	_, err := repo.coll.InsertOne(context.Background(), item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *MongoRepo) GetById(id string) (*posts.Post, error) {
	item := &posts.Post{}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = repo.coll.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *MongoRepo) Delete(postId string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return 0, err
	}
	res, err := repo.coll.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (repo *MongoRepo) AddComment(post *posts.Post, comment *posts.Comment) (int64, error) {
	res, err := repo.coll.UpdateOne(context.Background(),
		bson.M{"_id": post.MongoId},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (repo *MongoRepo) DeleteComment(post *posts.Post, commentId string) (int64, error) {
	res, err := repo.coll.UpdateOne(context.Background(),
		bson.M{"_id": post.MongoId},
		bson.M{"$pull": bson.M{"comments": bson.M{"id": commentId}}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (repo *MongoRepo) Vote(postId string, vote *posts.Vote) (int64, error) {

	_, err := repo.DeleteVote(postId, vote.UserId)

	oid, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return 0, err
	}
	_, err = repo.deleteVote(oid, vote.UserId)
	if err != nil {
		return 0, err
	}
	res, err := repo.coll.UpdateOne(context.Background(),
		bson.M{"_id": oid},
		bson.M{"$pull": bson.M{"votes": vote}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, err
}

func (repo *MongoRepo) UpdateStat(postId string, upvote, score int) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return 0, err
	}
	res, err := repo.coll.UpdateOne(context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"score": score, "upvotePercentage": upvote}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, err
}

func (repo *MongoRepo) DeleteVote(postId string, userId int) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return 0, err
	}
	return repo.deleteVote(oid, userId)
}

func (repo *MongoRepo) deleteVote(postId primitive.ObjectID, userId int) (int64, error) {
	res, err := repo.coll.UpdateOne(context.Background(),
		bson.M{"_id": postId},
		bson.M{"$pull": bson.M{"votes": bson.M{"userId": userId}}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (repo *MongoRepo) IncViews(post *posts.Post) (*posts.Post, error) {
	return nil, nil
}
