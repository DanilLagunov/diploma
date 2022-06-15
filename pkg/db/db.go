package db

import (
	"context"
	"errors"

	"github.com/DanilLagunov/diploma/pkg/models"
)

var ErrNotFound = errors.New("not found")

type Database interface {
	CreateUser(ctx context.Context, email, password string) error
	GetUser(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, email string, chatID int64) error
	GetUserCourses(ctx context.Context, chatId int64) ([]models.Course, error)
	UpdateUserCourses(ctx context.Context, chatId int64, course models.Course) error
	GetCourse(ctx context.Context, id string) (models.Course, error)
	GetCourses(ctx context.Context) ([]models.Course, error)
	CreateCourse(ctx context.Context, title, description string, lessons []models.Lesson) error
	GetCourseLessons(ctx context.Context, id string) ([]models.Lesson, error)
	CreateLesson(ctx context.Context, title, lection, task string, estimated int) (models.Lesson, error)
	GetLesson(ctx context.Context, id string) (models.Lesson, error)
}
