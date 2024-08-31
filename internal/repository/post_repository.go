package repository

import (
	"context"
	"github.com/AdblkA/blogging/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository struct {
	Collection *mongo.Collection
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	res, err := r.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	err = res.All(context.Background(), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetByID(id string) (models.Post, error) {
	var post models.Post
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Post{}, err
	}
	err = r.Collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: objID}}).Decode(&post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *PostRepository) Create(post *models.Post) (interface{}, error) {
	res, err := r.Collection.InsertOne(context.Background(), post)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil

}

func (r *PostRepository) Update(postID primitive.ObjectID, updatedPost *models.Post) (int64, error) {
	updateFields := bson.D{}

	if updatedPost.Title != "" {
		updateFields = append(updateFields, bson.E{Key: "title", Value: updatedPost.Title})
	}
	if len(updatedPost.Content) > 0 {
		updateFields = append(updateFields, bson.E{Key: "content", Value: updatedPost.Content})
	}
	if len(updatedPost.Category) > 0 {
		updateFields = append(updateFields, bson.E{Key: "category", Value: updatedPost.Category})
	}
	if len(updatedPost.Tags) > 0 {
		updateFields = append(updateFields, bson.E{Key: "tags", Value: updatedPost.Tags})
	}
	if !updatedPost.UpdatedAt.IsZero() {
		updateFields = append(updateFields, bson.E{Key: "updated_at", Value: updatedPost.UpdatedAt})
	}

	if len(updateFields) == 0 {
		return 0, nil
	}
	res, err := r.Collection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: postID}}, bson.D{{Key: "$set", Value: updateFields}})
	if err != nil {
		return 0, err
	}

	return res.ModifiedCount, nil
}

func (r *PostRepository) Delete(postID string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return 0, err
	}
	res, err := r.Collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: objID}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
