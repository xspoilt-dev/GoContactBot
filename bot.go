package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)

type User struct {
	UserName string `json:"user_name"`
	UserID   int    `json:"user_id"`
	Blocked  bool   `json:"blocked"`
}

var (
	userDataFile = ".user_data.json"
	adminChatID  = int64(1234567890) // Replace with your admin chat ID
	bot          *telebot.Bot
	replyQueue   = make(map[int64]int64)
)

func initializeUserDataFile() {
	if _, err := os.Stat(userDataFile); os.IsNotExist(err) {
		file, err := os.Create(userDataFile)
		if err != nil {
			log.Fatalf("Failed to create user data file: %v", err)
		}
		defer file.Close()
		file.Write([]byte("[]"))
	}
}

func getUserData() []User {
	var users []User
	file, err := os.Open(userDataFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil && err != io.EOF {
		log.Fatalf("Failed to decode user data: %v", err)
	}
	return users
}

func saveUserData(users []User) {
	file, err := os.Create(userDataFile)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		log.Fatalf("Failed to encode user data: %v", err)
	}
}

func addUser(userName string, userID int) {
	users := getUserData()
	for _, u := range users {
		if u.UserID == userID {
			return
		}
	}
	users = append(users, User{UserName: userName, UserID: userID})
	saveUserData(users)
}

func blockUser(userID int64) {
	users := getUserData()
	for i, u := range users {
		if u.UserID == int(userID) {
			users[i].Blocked = true
			saveUserData(users)
			return
		}
	}
}

func isUserBlocked(userID int) bool {
	users := getUserData()
	for _, u := range users {
		if u.UserID == userID && u.Blocked {
			return true
		}
	}
	return false
}

func main() {
	initializeUserDataFile()

	var err error
	bot, err = telebot.NewBot(telebot.Settings{
		Token:  "", // Replace with your bot token
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if m.Chat.ID == adminChatID || isUserBlocked(m.Sender.ID) {
			return
		}

		addUser(m.Sender.Username, m.Sender.ID)

		adminMessage := fmt.Sprintf(
			"📨 *Message from User:*\n\n👤 *Username:* `%s`\n🆔 *ID:* `%d`\n💬 *Message:* `%s`\n\n💬 /reply_%d\n🚫 /block_%d",
			m.Sender.Username, m.Sender.ID, m.Text, m.Sender.ID, m.Sender.ID,
		)

		_, err := bot.Send(&telebot.Chat{ID: adminChatID}, adminMessage, telebot.ModeMarkdown)
		if err != nil {
			log.Printf("Failed to send message to admin: %v", err)
		}
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if m.Chat.ID != adminChatID {
			return
		}

		if pendingUserID, exists := replyQueue[m.Chat.ID]; exists {
			delete(replyQueue, m.Chat.ID)
			_, err := bot.Send(&telebot.Chat{ID: pendingUserID}, "📬 *Admin Reply:* "+m.Text, telebot.ModeMarkdown)
			if err != nil {
				bot.Send(m.Chat, "Failed to send the reply to the user.")
				log.Printf("Failed to send reply: %v", err)
			} else {
				bot.Send(m.Chat, "Reply sent to the user.")
			}
			return
		}

		if strings.HasPrefix(m.Text, "/reply_") {
			userIDStr := strings.TrimPrefix(m.Text, "/reply_")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				bot.Send(m.Chat, "Invalid command format. Use /reply_<id>.")
				return
			}
			replyQueue[m.Chat.ID] = userID
			bot.Send(m.Chat, fmt.Sprintf("Now replying to user `%d`. Please send your message.", userID), telebot.ModeMarkdown)
		} else if strings.HasPrefix(m.Text, "/block_") {
			userIDStr := strings.TrimPrefix(m.Text, "/block_")
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				bot.Send(m.Chat, "Invalid command format. Use /block_<id>.")
				return
			}
			blockUser(userID)
			bot.Send(m.Chat, fmt.Sprintf("🚫 User `%d` has been blocked.", userID), telebot.ModeMarkdown)
		}
	})

	fmt.Println("Bot is running...")
	bot.Start()
}
