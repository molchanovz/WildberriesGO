package main

import (
	"bufio"
	"context"
	"firstProject/wildberriesFBS"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	// Store bot screaming status
	screaming = false
	bot       *tgbotapi.BotAPI

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

func main() {
	token, err := initEnv("variables.env", "TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text
	if user == nil {
		return
	}
	// Print to console
	log.Printf("%s написал %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message, text)
	} else if strings.HasPrefix(text, "WB") {
		sendTextMessage(message, "Готовлю файл")

		wildberriesKey, err := initEnv("variables.env", "API_KEY_WB")
		if err != nil {
			return
		}

		err = wildberriesFBS.GetReadyFile(wildberriesKey, text)
		if err != nil {
			sendTextMessage(message, err.Error())
		} else {
			sendMediaMessage(message, text)
			wildberriesFBS.Clean_files(text)
		}
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func handleCommand(message *tgbotapi.Message, command string) error {
	var err error
	switch command {
	case "/start":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Привет, "+message.From.FirstName)
		msg.Entities = message.Entities
		_, err = bot.Send(msg)
	case "/menu":
		err = sendMenu(message.Chat.ID)
		break
	}
	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == nextButton {
		text = secondMenu
		markup = secondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = firstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

func sendTextMessage(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.Entities = message.Entities
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения")
	}
}
func sendMediaMessage(message *tgbotapi.Message, text string) {
	// Открываем файл для чтения
	file, err := os.Open("wildberriesFBS/" + text + ".pdf")
	if err != nil {
		log.Fatal(err)
	}
	// Создаем объект для отправки документа
	document := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileReader{
		Name:   text + ".pdf",
		Reader: file,
	})
	// Отправляем документ
	_, err = bot.Send(document)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Документ успешно отправлен!")
}

func initEnv(path, name string) (string, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Ошибка загрузки файла %s: %v\n", path, err)
		return "", fmt.Errorf("ошибка загрузки файла " + path)
	}
	// Получаем значения переменных среды
	env := os.Getenv(name)

	if env == "" {
		return "", fmt.Errorf("переменная среды " + name + " не установлена")
	}
	return env, err
}
