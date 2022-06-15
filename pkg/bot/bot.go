package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/DanilLagunov/diploma/pkg/cache"
	"github.com/DanilLagunov/diploma/pkg/cache/memcache"
	"github.com/DanilLagunov/diploma/pkg/db"
	"github.com/DanilLagunov/diploma/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	db    db.Database
	cache cache.Cache
}

func New(db db.Database, token string) (*Bot, error) {
	cache := memcache.NewMemCache(60*time.Second, 90*time.Second)
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = true

	bot := &Bot{
		api:   api,
		db:    db,
		cache: cache,
	}

	log.Printf("Authorized on account %s", bot.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.api.GetUpdatesChan(u)

	bot.updateController(updates)

	return bot, nil
}

func (b *Bot) updateController(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil {
			split := strings.Split(update.CallbackQuery.Data, " ")
			switch split[0] {
			case "/register":
				b.courseRegistrationCallback(update, split)
			case "/lessons":
				b.courseLessonsCallback(update, split)
			case "/view":
				b.viewLessonCallback(update, split)
			default:
				continue
			}
		} else if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				b.help(update)
			case "start":
				b.start(update)
			case "contacts":
				b.contacts(update)
			case "login":
				b.login(update)
			case "courses":
				b.courses(update)
			case "mycourses":
				b.myCourses(update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "Невідома команда. Для списку команд введіть /help"
			}
		} else if update.Message != nil {
			switch update.Message.Text {
			case "Допомога":
				b.help(update)
			case "start":
				b.start(update)
			case "Контакти":
				b.contacts(update)
			case "login":
				b.login(update)
			case "Всі курси":
				b.courses(update)
			case "Мої курси":
				b.myCourses(update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "Невідома команда. Для списку команд введіть /help"
			}
		}
	}
}

func (b *Bot) start(update tgbotapi.Update) {
	var userKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Всі курси"),
			tgbotapi.NewKeyboardButton("Мої курси"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Контакти"),
			tgbotapi.NewKeyboardButton("Допомога"),
		),
	)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Це бот онлайн школи ApCenter\nДля входу до облікового запису введіть /login [пошта] [пароль]\nДля допомоги введіть /help."
	msg.ReplyMarkup = userKeyboard
	_, err := b.api.Send(msg)
	handleError(err)
}

func (b *Bot) help(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = `Для керування ботом ви можете використовувати вбудоване меню або текстові команди:
	/login [пошта] [пароль] - Вхід до облікового запису
	/help - Допомога
	/contacts - Контакти школи
	/courses - Всі курси школи
	/mycourses - Курси, на які ви підписані
	`
	_, err := b.api.Send(msg)
	handleError(err)
}

func (b *Bot) contacts(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "Пошта: school@apcenter.com \nТелефон: +380000000000\nАдреса: проспект Дмитра Яворницького, 35"
	_, err := b.api.Send(msg)
	handleError(err)
}

func (b *Bot) login(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	split := strings.Split(update.Message.Text, " ")
	if len(split) < 3 {
		msg.Text = "Недостатньо аргументів, спробуйте ще раз."
		_, err := b.api.Send(msg)
		handleError(err)
		return
	}

	user, err := b.cache.GetUser(split[1])
	if err != nil {
		fmt.Printf("cache error: %s", err)
		user, err = b.db.GetUser(context.TODO(), split[1])
		if errors.Is(err, db.ErrNotFound) {
			msg.Text = "Користувача не знайдено, спробуйте ще раз."
			_, err = b.api.Send(msg)
			handleError(err)
			return
		}
	}
	b.cache.SetUser(split[1], user, 0)

	if !utils.CheckPasswordHash(split[2], user.Password) {
		fmt.Println("wrong password")
	}

	err = b.db.UpdateUser(context.TODO(), user.Email, update.Message.Chat.ID)
	handleError(err)

	msg.Text = "Авторизація успішна!"
	_, err = b.api.Send(msg)
	handleError(err)
}

func (b *Bot) courses(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	courses, err := b.db.GetCourses(context.TODO())
	handleError(err)

	for _, course := range courses {
		var courseRegistration = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Приєднатися!", fmt.Sprintf("/register %s", course.ID)),
			),
		)
		msg.Text = fmt.Sprintln(course.Title + "\n" + course.Description + "\nКількість уроків: " + strconv.Itoa(len(course.Lessons)))
		msg.ReplyMarkup = courseRegistration
		_, err = b.api.Send(msg)
		handleError(err)
	}
}

func (b *Bot) myCourses(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	user, err := b.cache.GetUser(fmt.Sprint(update.Message.Chat.ID))
	if err != nil {
		fmt.Printf("cache error: %s", err)
		user, err = b.db.GetUser(context.TODO(), fmt.Sprint(update.Message.Chat.ID))
		if errors.Is(err, db.ErrNotFound) {
			msg.Text = "Ви не авторизовані!"
			_, err = b.api.Send(msg)
			handleError(err)
			return
		}
	}
	courses, err := b.db.GetUserCourses(context.TODO(), user.ChatID)
	handleError(err)

	for _, course := range courses {
		var courseLessons = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Список уроків", fmt.Sprintf("/lessons %s", course.ID)),
			),
		)
		msg.Text = fmt.Sprintln(course.Title + "\n" + course.Description + "\n" + strconv.Itoa(len(course.Lessons)))
		msg.ReplyMarkup = courseLessons
		_, err = b.api.Send(msg)
		handleError(err)
	}
}

func (b *Bot) courseRegistrationCallback(update tgbotapi.Update, split []string) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	user, err := b.cache.GetUser(fmt.Sprint(update.CallbackQuery.From.ID))
	if err != nil {
		fmt.Printf("cache error: %s", err)
		user, err = b.db.GetUser(context.TODO(), fmt.Sprint(update.CallbackQuery.From.ID))
		if errors.Is(err, db.ErrNotFound) {
			msg.Text = "Ви не авторизовані!"
			_, err = b.api.Send(msg)
			handleError(err)
			return
		}
	}
	course, err := b.cache.GetCourse(split[1])
	if err != nil {
		fmt.Printf("cache error: %s", err)
		course, err = b.db.GetCourse(context.TODO(), split[1])
		if errors.Is(err, db.ErrNotFound) {
			handleError(err)
			return
		}
	}
	b.cache.SetCourse(split[1], course, 0)

	err = b.db.UpdateUserCourses(context.TODO(), user.ChatID, course)
	handleError(err)
}

func (b *Bot) courseLessonsCallback(update tgbotapi.Update, split []string) {
	lessons, err := b.db.GetCourseLessons(context.TODO(), split[1])
	handleError(err)

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	for i, lesson := range lessons {
		var viewLesson = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Перейти до уроку", fmt.Sprintf("/view %s %s", lesson.ID, split[1])),
			),
		)
		msg.Text = fmt.Sprintf("Урок %d. %s", i+1, lesson.Title)
		msg.ReplyMarkup = viewLesson
		_, err = b.api.Send(msg)
		handleError(err)
	}
}

func (b *Bot) viewLessonCallback(update tgbotapi.Update, split []string) {
	lesson, err := b.cache.GetLesson(split[1])
	if err != nil {
		fmt.Printf("cache error: %s", err)
		lesson, err = b.db.GetLesson(context.TODO(), split[1])
		if errors.Is(err, db.ErrNotFound) {
			handleError(err)
			return
		}
	}
	b.cache.SetLesson(split[1], lesson, 0)

	msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, "")
	var back = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", fmt.Sprintf("/lessons %s", split[2])),
		),
	)
	msg.Text = fmt.Sprintf("%s\nЛекція: %s\nЗавдання: %s\nЧас виконання: %d хвилин", lesson.Title, lesson.Lection, lesson.Task, lesson.EstimatedTime)
	msg.ReplyMarkup = back
	_, err = b.api.Send(msg)
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
}
