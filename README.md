# gift-buyer

An automated Telegram userbot for monitoring and purchasing limited edition Telegram star gifts. The bot continuously monitors Telegram's star gift catalog and automatically purchases limited gifts when they become available.

## Features

- **Real-time Monitoring**: Continuously monitors Telegram's star gift catalog for changes
- **Auto-purchase**: Automatically purchases limited edition gifts based on availability thresholds
- **Activity Monitoring**: Built-in crash detection that restarts the bot if no activity is detected
- **Configurable Polling**: Adjustable monitoring intervals
- **Admin Notifications**: Sends real-time updates to configured admin about gift availability
- **Session Management**: Secure session handling with persistent authentication
- **Comprehensive Logging**: Debug-level logging with caller information and timestamps

## Prerequisites

- Go 1.24 or higher
- Telegram account
- Telegram API credentials (App ID and App Hash)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/chaindead/gift-buyer.git
cd gift-buyer
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o gift-buyer ./bot
```

## Configuration

### Environment Variables

Create a `.env` file or set the following environment variables:

```bash
# Required: Telegram admin username to receive notifications
TG_ADMIN=your_username

# Required: Telegram API credentials
TG_APP_ID=your_app_id
TG_API_HASH=your_api_hash

# Optional: Session string for authentication (generated via auth command)
TG_SESSION=your_session_string

# Optional: Polling interval (default: 1s)
POLL_INTERVAL=1s
```

### Telegram API Setup

To obtain your Telegram API credentials:

1. Go to [my.telegram.org](https://my.telegram.org)
2. Log in with your phone number
3. Go to **API Development Tools**
4. Create a new application
5. Note down your `App ID` and `App Hash`

For detailed setup instructions, refer to the [Telegram MCP Setup Guide](https://raw.githubusercontent.com/chaindead/telegram-mcp/refs/heads/main/README.md).

## Usage

### 1. Authentication (First Time Setup)

Before running the bot, you need to authenticate and obtain a session string:

```bash
./gift-buyer -auth
```

This will:
- Prompt you to enter your phone number
- Send you a verification code via Telegram
- Generate a session string for future use
- Display the session string to add to your environment variables

### 2. Running the Bot

Once authenticated, run the bot with:

```bash
./gift-buyer
```

The bot will:
- Connect to Telegram using your session
- Start monitoring the star gift catalog
- Send notifications to the configured admin
- Automatically purchase limited gifts based on availability criteria

### Command Line Options

- `-auth`: Run authentication mode to obtain session string

## Project Structure

```
gift-buyer/
├── auth/           # Authentication handling
│   └── auth.go     # Session authentication logic
├── bot/            # Main bot functionality
│   ├── bot.go      # Main bot logic and gift monitoring
│   ├── send.go     # Gift purchasing and session management
│   └── watcher.go  # Gift upgrade monitoring
├── config/         # Configuration management
│   └── config.go   # Environment variable parsing
├── log/            # Logging setup
│   └── init.go     # Zerolog initialization
├── go.mod          # Go module dependencies
└── README.md       # This file
```

## Configuration Details

### Gift Purchase Criteria

The bot automatically purchases gifts that meet these criteria:
- Total availability ≤ 50,000 units
- Remaining availability > 0
- Sorted by total count (smallest first)
- Maximum 100 purchases per gift

### Monitoring Behavior

- **Polling Interval**: Configurable via `POLL_INTERVAL` (default: 1 second)
- **Activity Monitoring**: Crashes if no activity detected for 1 minute
- **Hash-based Updates**: Only processes actual changes in gift catalog
- **Notification Format**: HTML-formatted messages with gift details

## Logging

The application uses structured logging with:
- Console output with timestamps
- Caller information for debugging
- Stack traces for errors
- Debug level logging enabled

## Security Considerations

- **Session Security**: Keep your session string secure and never commit it to version control
- **API Limits**: Be mindful of Telegram's rate limits
- **Auto-purchase**: Monitor your star balance as the bot will automatically spend stars

## Dependencies

- `github.com/amarnathcjd/gogram` - Telegram client library
- `github.com/caarlos0/env/v11` - Environment variable parsing
- `github.com/rs/zerolog` - Structured logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Troubleshooting

### Common Issues

1. **Authentication Errors**: Ensure your App ID and App Hash are correct
2. **Session Expired**: Re-run with `-auth` flag to generate a new session
3. **No Activity**: Check your internet connection and Telegram API access
4. **Permission Denied**: Ensure the admin username is correct and accessible

### Debug Mode

The application runs in debug mode by default. Check the console output for detailed information about:
- Gift catalog updates
- Purchase attempts
- API responses
- Error messages

## Support

For issues related to:
- Telegram API: Check [Telegram API Documentation](https://core.telegram.org/api)
- Gogram Library: Visit [Gogram GitHub](https://github.com/amarnathcjd/gogram)
- This Project: Open an issue in this repository
