package mongo

import (
	"context"
	"time"

	"github.com/DanilLagunov/diploma/pkg/db"
	"github.com/DanilLagunov/diploma/pkg/models"
	"github.com/DanilLagunov/diploma/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database struct.
type Database struct {
	client            *mongo.Client
	usersCollection   *mongo.Collection
	coursesCollection *mongo.Collection
	lessonsCollection *mongo.Collection
}

// NewDatabase creating a new Database object.
func New(uri, dbName, usersCollectionName, coursesCollectionName, lessonsCollectionName string) (*Database, error) {
	var db Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	db.client = client
	db.usersCollection = client.Database(dbName).Collection(usersCollectionName)
	db.coursesCollection = client.Database(dbName).Collection(coursesCollectionName)
	db.lessonsCollection = client.Database(dbName).Collection(lessonsCollectionName)

	return &db, err
}

// USER DB HANDLERS

func (d *Database) CreateUser(ctx context.Context, email, password string) error {
	id := primitive.NewObjectID().Hex()
	hashedPassword, _ := utils.HashPassword(password)
	user := models.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
	}
	_, err := d.usersCollection.InsertOne(ctx, user)
	return err
}

func (d *Database) GetUser(ctx context.Context, email string) (models.User, error) {
	filter := bson.M{"email": email}
	var user models.User
	err := d.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, db.ErrNotFound
	}
	return user, err
}

func (d *Database) UpdateUser(ctx context.Context, email string, chatID int64) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"chat_id": chatID}}
	_, err := d.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetUserCourses(ctx context.Context, chatId int64) ([]models.Course, error) {
	filter := bson.M{"chat_id": chatId}
	var user models.User
	err := d.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return []models.Course{}, db.ErrNotFound
	}
	return user.Courses, nil
}

func (d *Database) UpdateUserCourses(ctx context.Context, chatId int64, course models.Course) error {
	filter := bson.M{"chat_id": chatId}
	var user models.User
	err := d.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return db.ErrNotFound
	}
	for _, c := range user.Courses {
		if c.ID == course.ID {
			return nil
		}
	}
	user.Courses = append(user.Courses, course)
	update := bson.M{"$set": bson.M{"courses": user.Courses}}
	_, err = d.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// COURSES DB HANDLERS

func (d *Database) GetCourse(ctx context.Context, id string) (models.Course, error) {
	filter := bson.M{"_id": id}
	var course models.Course
	err := d.coursesCollection.FindOne(ctx, filter).Decode(&course)
	if err == mongo.ErrNoDocuments {
		return course, db.ErrNotFound
	}
	return course, err
}

func (d *Database) GetCourses(ctx context.Context) ([]models.Course, error) {
	cur, err := d.coursesCollection.Find(ctx, bson.M{})
	if err != nil {
		return []models.Course{}, err
	}
	defer cur.Close(ctx)
	result := []models.Course{}
	if err := cur.All(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (d *Database) CreateCourse(ctx context.Context, title, description string, lessons []models.Lesson) error {
	id := primitive.NewObjectID().Hex()
	course := models.Course{
		ID:          id,
		Title:       title,
		Description: description,
		Lessons:     lessons,
	}
	_, err := d.coursesCollection.InsertOne(ctx, course)
	return err
}

func (d *Database) GetCourseLessons(ctx context.Context, id string) ([]models.Lesson, error) {
	filter := bson.M{"_id": id}
	var course models.Course
	err := d.coursesCollection.FindOne(ctx, filter).Decode(&course)
	if err == mongo.ErrNoDocuments {
		return []models.Lesson{}, db.ErrNotFound
	}
	return course.Lessons, nil
}

// LESSONS DB HANDLERS

func (d *Database) CreateLesson(ctx context.Context, title, lection, task string, estimated int) (models.Lesson, error) {
	id := primitive.NewObjectID().Hex()
	lesson := models.Lesson{
		ID:            id,
		Title:         title,
		Lection:       lection,
		Task:          task,
		EstimatedTime: estimated,
	}
	_, err := d.lessonsCollection.InsertOne(ctx, lesson)
	return lesson, err
}

func (d *Database) GetLesson(ctx context.Context, id string) (models.Lesson, error) {
	filter := bson.M{"_id": id}
	var lesson models.Lesson
	err := d.lessonsCollection.FindOne(ctx, filter).Decode(&lesson)
	if err == mongo.ErrNoDocuments {
		return lesson, db.ErrNotFound
	}
	return lesson, err
}
