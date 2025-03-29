package repo

import (
	"TodoApp/internal/domain"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TaskRepository struct {
	Db *mongo.Client
}

var _ domain.TaskRepoImp = (*TaskRepository)(nil)

func DefaultTaskRepository(Db *mongo.Client) *TaskRepository {
	return &TaskRepository{Db: Db}
}

// Выполняет запись новой задачи в бд
func (repo *TaskRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	res, err := repo.Db.Database("To-Do").Collection("tasks").InsertOne(ctx, task)
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

// Получает задачу по уникальному id
func (repo *TaskRepository) GetTask(ctx context.Context, id bson.ObjectID) (domain.Task, error) {
	var task domain.Task
	collection := repo.Db.Database("To-Do").Collection("tasks")
	err := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&task)
	if err != nil {
		return task, err
	}
	return task, nil
}

// Получает список задач из бд
func (repo *TaskRepository) GetTasks(ctx context.Context, status string) ([]domain.Task, error) {
	var tasks []domain.Task
	collection := repo.Db.Database("To-Do").Collection("tasks")
	opt := options.Aggregate()
	var statement string
	if status == "active" {
		statement = "$lte"
	} else {
		statement = "$gt"
	}

	pipeline := mongo.Pipeline{
		// Фильтр по дате
		{{Key: "$match", Value: bson.D{
			{Key: "activeAt", Value: bson.D{{Key: statement, Value: time.Now()}}},
		}}},
		// Сортировка по activeAt
		{{Key: "$sort", Value: bson.D{{Key: "activeAt", Value: 1}}}},
		// Добавляем поле dayOfWeek (1 = воскресенье, 7 = суббота)
		{{Key: "$addFields", Value: bson.D{
			{Key: "dayOfWeek", Value: bson.D{{Key: "$dayOfWeek", Value: "$activeAt"}}},
		}}},
		// Добавляем isWeekend: true, если день недели = 1 (воскресенье) или 7 (суббота)
		{{Key: "$addFields", Value: bson.D{
			{Key: "isWeekend", Value: bson.D{
				{Key: "$in", Value: bson.A{
					bson.D{{Key: "$dayOfWeek", Value: "$activeAt"}},
					bson.A{1, 7}, // Проверка на выходной
				}},
			}},
		}}},
		// Если выходной — добавляем "ВЫХОДНОЙ-" к заголовку
		{{Key: "$addFields", Value: bson.D{
			{Key: "title", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: "$isWeekend"},
					{Key: "then", Value: bson.D{{Key: "$concat", Value: bson.A{"ВЫХОДНОЙ- ", "$title"}}}},
					{Key: "else", Value: "$title"},
				}},
			}},
		}}},
		// Выводим нужные поля
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "title", Value: 1},
			{Key: "status", Value: 1},
			{Key: "activeAtStr", Value: 1},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Обновляет существующую задачу в бд
func (repo *TaskRepository) UpdateTask(ctx context.Context, id bson.ObjectID, task domain.Task) error {
	collection := repo.Db.Database("To-Do").Collection("tasks")
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"activeAt":    task.ActiveDateTime,
			"activeAtStr": task.ActiveDateStr,
			"status":      "active",
		},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Удаляет задачу из бд
func (repo *TaskRepository) DeleteTask(ctx context.Context, id bson.ObjectID) error {
	collection := repo.Db.Database("To-Do").Collection("tasks")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Изменяет поле статуса задачи на done
func (repo *TaskRepository) ChangeStatus(ctx context.Context, id bson.ObjectID) error {
	collection := repo.Db.Database("To-Do").Collection("tasks")
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

// Проверка уникальности задачи
func (repo *TaskRepository) IsTaskUnique(ctx context.Context, task domain.Task) (bool, error) {
	collection := repo.Db.Database("To-Do").Collection("tasks")
	filter := bson.M{
		"$and": []bson.M{
			{"title": task.Title},
			{"activeAt": task.ActiveDateStr},
		},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
