package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DanilLagunov/diploma/pkg/bot"
	"github.com/DanilLagunov/diploma/pkg/db/mongo"
)

func main() {
	db, err := mongo.New(
		fmt.Sprintf("mongodb+srv://%s:%s@diploma.g0hbuep.mongodb.net/?retryWrites=true&w=majority",
			os.Getenv("MONGO_USER"),
			os.Getenv("MONGO_PASSWORD")),
		os.Getenv("MONGO_DB"),
		os.Getenv("MONGO_USERS_COLLECTION"),
		os.Getenv("MONGO_COURSES_COLLECTION"),
		os.Getenv("MONGO_LESSONS_COLLECTION"))
	if err != nil {
		log.Panic(err)
	}

	// lessnon1, _ := db.CreateLesson(
	// 	context.TODO(),
	// 	"Склад числа.",
	// 	"https://youtu.be/a1M-tWI000k",
	// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
	// 	90,
	// )
	// lessnon2, _ := db.CreateLesson(
	// 	context.TODO(),
	// 	"Додавання та віднімання.",
	// 	"https://youtu.be/K37RnNkGGcI",
	// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
	// 	75,
	// )
	// lessnon3, _ := db.CreateLesson(
	// 	context.TODO(),
	// 	"Прості числа.",
	// 	"https://youtu.be/LQ4brH4zN1M",
	// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
	// 	80,
	// )

	// db.CreateUser(context.TODO(), "user1@gmail.com", "user")
	// db.CreateCourse(context.TODO(), "Математика для початківців", "Курс математики для учнів початкових класів.", []models.Lesson{lessnon1, lessnon2, lessnon3})

	token := os.Getenv("BOT_TOKEN")

	bot.New(db, token)
}

// lessnon1, _ := db.CreateLesson(
// 	context.TODO(),
// 	"Введення до вищої математики.",
// 	"https://youtu.be/Jkb7enPFW88",
// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
// 	90,
// )
// lessnon2, _ := db.CreateLesson(
// 	context.TODO(),
// 	"Інтеграли.",
// 	"https://youtu.be/j2FK5MGg35k",
// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
// 	75,
// )
// lessnon3, _ := db.CreateLesson(
// 	context.TODO(),
// 	"Будова організму. Клітини.",
// 	"https://youtu.be/HOG_jVteKEk",
// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
// 	80,
// )
// lessnon4, _ := db.CreateLesson(
// 	context.TODO(),
// 	"Будова організму. Тканини.",
// 	"https://youtu.be/HZyF2XTMDxk",
// 	"https://docs.google.com/forms/d/e/1FAIpQLSfOxe6ueK2op-wii75bDVcSdTKJNxWvyH8s-TFGFM2SVbgMvA/viewform",
// 	65,
// )

// db.CreateUser(context.TODO(), "raccoo@gmail.com", "raccoo")
// db.CreateCourse(context.TODO(), "Вища математика", "Курс вищої математики для студентів вищих навчальних закладів.", []models.Lesson{lessnon1, lessnon2})
// db.CreateCourse(context.TODO(), "Біологія", "Курс з біології для учнів середньої та старшої школи.", []models.Lesson{lessnon3, lessnon4})
