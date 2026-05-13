# Telegram Internet Monitoring

A high-performance automated network monitoring tool written in **Go**. This bot tracks internet connection quality, detects location data via reverse geocoding, and sends diagnostic reports to a designated Telegram chat.

---

## Features

- **Real-Time Diagnostics**: Measures Download speed, Upload speed, Latency (Ping), and Jitter.
- **Geolocation Intelligence**: Automatically detects ISP and physical location (City & District) using IP-based coordinates.
- **Automated Monitoring**: Scheduled reporting every 30 minutes to track network stability.
- **Anomaly Detection**: Visual indicators if the network falls below pre-defined thresholds.
- **Interactive Commands**: Control the bot directly via Telegram commands.
- **WIB Timezone**: Configured for Western Indonesian Time (Asia/Jakarta).

---

## Telegram Commands

| Command     | Description                                     |
| :---------- | :---------------------------------------------- |
| `/check`    | Triggers an immediate, manual speed test.       |
| `/startbot` | Enables the automated 30-minute status reports. |
| `/stopbot`  | Pauses the automated reporting cycle.           |
| `/menu`     | Displays the help menu and available actions.   |

---

## Getting Started

### 1. Prerequisites

- **Go** (1.18 or higher)
- **Telegram Bot Token**: Obtainable from @BotFather.
- **Telegram Chat ID**: The unique ID of the chat where reports will be sent.

### 2. Installation & Configuration

Install the necessary dependencies:

```bash
go get [github.com/go-telegram-bot-api/telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
go get [github.com/showwin/speedtest-go/speedtest](https://github.com/showwin/speedtest-go/speedtest)
Update the constants in main.go:

Go
const (
    botToken     = "YOUR_BOT_TOKEN"
    chatID       = YOUR_CHAT_ID
    autoInterval = 30 * time.Minute
)
Running in Background (Windows via NSSM)
To run the bot as a background service on Windows without an active terminal window, use NSSM (Non-Sucking Service Manager).

1. Build the Executable
Compile the Go code into a Windows binary:

PowerShell
go build -o network-guard.exe
2. Setup NSSM
Download NSSM from nssm.cc.

Copy win64/nssm.exe into your project folder.

Open Command Prompt or PowerShell as Administrator.

Run the installer:

PowerShell
.\nssm.exe install NetworkGuard
3. Service Configuration
In the GUI window:

Path: Select your network-guard.exe.

Startup directory: Set to your project folder.

I/O Tab: (Optional) Set stdout and stderr to a log.txt file to monitor output.

Click Install Service.

4. Service Management
Start: nssm start NetworkGuard

Stop: nssm stop NetworkGuard

Remove: nssm remove NetworkGuard

Sample Report Layout
SYSTEM OPERATIONAL
━━━━━━━━━━━━━━━━━━
DETECTED LOCATION
├ City : South Jakarta
├ District : Tebet
━━━━━━━━━━━━━━━━━━
NETWORK PROFILE
├ ISP : Biznet Home
└ Test Node : Jakarta (ID)
━━━━━━━━━━━━━━━━━━
DIAGNOSTIC RESULTS
📥 Download : 94.50 Mbps
📤 Upload : 92.10 Mbps
⚡ Latency : 5ms
━━━━━━━━━━━━━━━━━━
REPORT TIME
03 Apr 2026 | 08:45:00 WIB

Security Note
For security, it is recommended to use Environment Variables for the botToken and chatID rather than hardcoding them if the source code is hosted on public repositories.
```
