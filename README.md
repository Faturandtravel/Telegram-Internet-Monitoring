## 📡 Network Guard Bot

A high-precision Telegram bot built with **Go (Golang)** designed to monitor network performance. It provides automated speed tests every 30 minutes and real-time diagnostics on demand, optimized for high-speed fiber connections (100Mbps+).

## ✨ Key Features

- **High-Accuracy Testing**: Uses multi-threading and smart server selection to match official Speedtest standards.
- **Real-Time Geolocation**: Automatically detects City, Country, and GPS Coordinates based on your Public IP.
- **WIB Timezone Support**: Reports are timestamped to Western Indonesia Time (Asia/Jakarta).
- **Automated Monitoring**: Periodic network health checks sent directly to your Telegram chat.
- **Anomaly Alerts**: Visual indicators (🔴) trigger if the download speed drops below your defined threshold.
- **Manual Diagnostics**: Run an instant check anytime using the `/checknow` command.

## 🚀 Getting Started

### 1. Prerequisites

Ensure you have [Go](https://golang.org/doc/install) installed on your machine.

### 2. Installation

Clone the repository and install the required dependencies:

```bash
git clone [https://github.com/your-username/network-guard-bot.git](https://github.com/your-username/network-guard-bot.git)
cd network-guard-bot
go get [github.com/go-telegram-bot-api/telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
go get [github.com/showwin/speedtest-go/speedtest](https://github.com/showwin/speedtest-go/speedtest)
3. ConfigurationOpen main.go and update the following constants:botToken: Your unique bot token from @BotFather.chatID: Your Telegram User ID or Group ID.minDownload: Set your minimum expected speed (e.g., 60.0 for a 100Mbps plan).4. Running the BotBashgo run main.go
🛠️ Bot CommandsCommandDescription/checknowExecute an immediate high-precision speed test/menuDisplay the help menu and available commands📊 Sample Report OutputPlaintext🌐 🟢 SYSTEM OPERATIONAL
━━━━━━━━━━━━━━━━━━
📍 DETECTED LOCATION
├ City : Jakarta
├ Country : Indonesia
└ Coordinates : -6.2146, 106.8451
━━━━━━━━━━━━━━━━━━
📡 NETWORK PROFILE
├ ISP : PT Telekomunikasi Indonesia
├ Public IP : 180.242.xx.xx
└ Test Node : Jakarta (ID)
━━━━━━━━━━━━━━━━━━
📊 PERFORMANCE METRICS
📥 Download : 94.45 Mbps
📤 Upload : 30.12 Mbps
⚡ Latency : 4ms (Jitter: 1ms)
━━━━━━━━━━━━━━━━━━
🕒 REPORT TIMESTAMP
02 Apr 2026 | 19:45:00 WIB
⚠️ Security WarningDo not commit your botToken to public repositories. It is recommended to use environment variables for sensitive credentials in production environments.📝 LicenseThis project is licensed under the MIT License.
### Quick Tips for your GitHub Repo:
1.  **File Name**: Save the text above as `README.md`.
2.  **License**: You can add a `LICENSE` file with the MIT license text to make it official.
3.  **Bot Token**: If you already pushed the code with the token visible, **revoke the token** via BotFather and generate a new one immediately.
```
