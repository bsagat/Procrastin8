package repo

import (
	"TodoApp/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TaskRepository struct {
	Db *mongo.Client
}

func DefaultTaskRepository(Db *mongo.Client) *TaskRepository {
	return &TaskRepository{Db: Db}
}

func (repo *TaskRepository) CreateTask(ctx context.Context, task models.Task) error {
	_, err := repo.Db.Database("To-do").Collection("tasks").InsertOne(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepository) IsTaskUnique(ctx context.Context, task models.Task) (bool, error) {
	collection := repo.Db.Database("To-do").Collection("tasks")
	filter := bson.M{
		"$and": []bson.M{
			{"title": task.Title},
			{"activeAt": task.ActiveDate},
		},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (repo *TaskRepository) GetTasks(ctx context.Context, status string) ([]models.Task, error) {
	var tasks []models.Task
	collection := repo.Db.Database("To-do").Collection("tasks")
	findoptions := options.Find()
	cursor, err := collection.Find(ctx, bson.M{"status": status}, findoptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *TaskRepository) ChangeStatus(ctx context.Context, id bson.ObjectID) error {
	collection := repo.Db.Database("To-do").Collection("tasks")
	update := bson.M{"$set": bson.M{"status": "done"}}
	res, err := collection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (repo *TaskRepository) DeleteTask(ctx context.Context, id bson.ObjectID) error {
	collection := repo.Db.Database("To-do").Collection("tasks")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
