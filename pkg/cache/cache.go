package cache

import (
	"errors"
	"time"

	"github.com/DanilLagunov/diploma/pkg/models"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrItemExpired = errors.New("item expired")
)

type UserItem struct {
	Value      models.User
	Created    time.Time
	Expiration int64
}

type CourseItem struct {
	Value      models.Course
	Created    time.Time
	Expiration int64
}

type LessonItem struct {
	Value      models.Lesson
	Created    time.Time
	Expiration int64
}

// Cache interface.
type Cache interface {
	GetUser(key string) (models.User, error)
	GetCourse(key string) (models.Course, error)
	GetLesson(key string) (models.Lesson, error)
	SetUser(key string, value models.User, duration time.Duration)
	SetCourse(key string, value models.Course, duration time.Duration)
	SetLesson(key string, value models.Lesson, duration time.Duration)
}
