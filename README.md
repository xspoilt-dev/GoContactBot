# GoContactBot

GoContactBot is a Telegram bot written in Go that allows an admin to manage user interactions. The bot can add users, block users, and facilitate communication between the admin and users.

## Features

- **Add Users**: Automatically adds users who send messages to the bot.
- **Block Users**: Allows the admin to block users from sending messages.
- **Reply to Users**: Admin can reply to users directly through the bot.
- **Persistent User Data**: User data is stored in a JSON file for persistence.

## Getting Started

### Prerequisites

- Go 1.23.2 or later
- A Telegram bot token from [BotFather](https://core.telegram.org/bots#botfather)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/xspoilt-dev/GoContactBot.git
    cd GoContactBot
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Create a `.user_data.json` file in the root directory:
    ```sh
    touch .user_data.json
    ```

4. Replace the `adminChatID` and `Token` placeholders in `main.go` with your actual admin chat ID and bot token:
    ```go
    var (
        adminChatID = int64(1234567890) // Replace with your admin chat ID
        bot, err = telebot.NewBot(telebot.Settings{
            Token:  "YOUR_BOT_TOKEN", // Replace with your bot token
            Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
        })
    )
    ```

### Usage

1. Run the bot:
    ```sh
    go run main.go
    ```

2. Interact with the bot on Telegram:
    - Send a message to the bot to be added as a user.
    - Admin can reply to users using the `/reply_<id>` command.
    - Admin can block users using the `/block_<id>` command.

### Customization

- **User Data File**: The user data is stored in `.user_data.json`. You can change the file name by modifying the `userDataFile` variable in `main.go`.
    ```go
    var userDataFile = ".user_data.json"
    ```

- **Admin Chat ID**: Set the `adminChatID` variable to your Telegram chat ID to receive messages from users.
    ```go
    var adminChatID = int64(1234567890) // Replace with your admin chat ID
    ```

- **Bot Token**: Set the `Token` field in the `telebot.Settings` struct to your bot token.
    ```go
    bot, err = telebot.NewBot(telebot.Settings{
        Token:  "YOUR_BOT_TOKEN", // Replace with your bot token
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    ```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
