package repo

import (
	"TodoApp/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskRepository struct {
	Db *mongo.Client
}

func DefaultTaskRepository(Db *mongo.Client) *TaskRepository {
	return &TaskRepository{Db: Db}
}

func (repo *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	res, err := repo.Db.Database("To-do").Collection("tasks").InsertOne(ctx, task)
	if err != nil {
		return err
	}
	var ok bool
	task.Id, ok = res.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("objectID convert error")
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

/*
db.tasks.aggregate([
  {
    $match: {
      activeAt: { $lte: new Date() } // Фильтр: задачи с датой до сегодня
    }
  },
  {
    $sort: { createdAt: 1 } // Сортировка по дате создания (или activeAt)
  },
  {
    $addFields: {
      dayOfWeek: { $dayOfWeek: "$activeAt" }, // Получаем день недели (1=вс, 7=сб)
      isWeekend: { $in: [{ $dayOfWeek: "$activeAt" }, [1, 7]] } // Проверяем выходной
    }
  },
  {
    $addFields: {
      title: {
        $cond: {
          if: "$isWeekend",
          then: { $concat: ["ВЫХОДНОЙ- ", "$title"] },
          else: "$title"
        }
      }
    }
  },
  {
    $project: {
      _id: 1,
      title: 1,
      activeAt: 1
    }
  }
]);
*/

func (repo *TaskRepository) GetTasks(ctx context.Context, status string) ([]models.Task, error) {
	var tasks []models.Task
	collection := repo.Db.Database("To-do").Collection("tasks")
	collection.Aggregate(ctx, []interface{})
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
func (repo *TaskRepository) UpdateTask(ctx context.Context, id bson.ObjectID, task models.Task) error {
	collection := repo.Db.Database("To-do").Collection("tasks")
	filter := bson.M{"title": task.Title, "activeAt": task.ActiveDate}
	res, err := collection.UpdateByID(ctx, id, filter)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
