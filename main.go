package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/showwin/speedtest-go/speedtest"
)

const (
	botToken     = "8680477535:AAGoZjz6DB9nIt_FZdICFOOAvB12xjFt0Ag"
	chatID       = -5275988185
	autoInterval = 30 * time.Minute
	pingLimit    = 50 * time.Millisecond
	minDownload  = 60.0
)

var (
	autoReportActive = true
	mu               sync.RWMutex
)

func getWIBTime() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now().Add(7*time.Hour).Format("02 Jan 2006 | 15:04:05") + " WIB"
	}
	return time.Now().In(loc).Format("02 Jan 2006 | 15:04:05") + " WIB"
}

func getLocationAddress(lat, lon string) (city, district string) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s&zoom=12&accept-language=en", lat, lon)

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "inetmonitor/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "Unknown", "Unknown"
	}
	defer resp.Body.Close()

	var result struct {
		Address struct {
			City          string `json:"city"`
			Town          string `json:"town"`
			Regency       string `json:"county"`
			Suburb        string `json:"suburb"`
			Village       string `json:"village"`
			District      string `json:"district"`
			StateDistrict string `json:"state_district"`
			State         string `json:"state"`
		} `json:"address"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Unknown", "Unknown"
	}

	city = result.Address.City
	if city == "" {
		city = result.Address.Town
	}
	if city == "" {
		city = result.Address.StateDistrict
	}
	if city == "" {
		city = result.Address.State
	}

	district = result.Address.District
	if district == "" {
		district = result.Address.Suburb
	}
	if district == "" {
		district = result.Address.Village
	}
	if district == "" {
		district = result.Address.Regency
	}

	return city, district
}

func runSpeedTest() string {
	log.Println("⏳ Starting Real-Time Diagnostics...")

	stClient := speedtest.New()
	user, err := stClient.FetchUserInfo()
	if err != nil {
		return "🚨 *SYSTEM ERROR*: Failed to detect location & ISP via IP."
	}

	serverList, _ := stClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	if len(targets) == 0 {
		return "🚨 *SYSTEM ERROR*: No test server responded."
	}

	s := targets[0]
	s.PingTest(nil)
	s.DownloadTest()
	s.UploadTest()

	downloadMbps := float64(s.DLSpeed) / 125000
	uploadMbps := float64(s.ULSpeed) / 125000
	latency := s.Latency.Round(time.Millisecond)

	isAnomalous := latency > pingLimit || downloadMbps < minDownload

	statusHeader := "🟢 SYSTEM OPERATIONAL"
	if isAnomalous {
		statusHeader = "🔴 NETWORK ANOMALY"
	}

	alertMsg := ""
	if isAnomalous {
		alertMsg = "\n⚠️ *ADVISORY*:\n- Performance below 100Mbps standard.\n- Check physical connection or ISP outage.\n"
	}

	city, district := getLocationAddress(user.Lat, user.Lon)
	currentTime := getWIBTime()

	return fmt.Sprintf(
		"🌐 *%s*\n"+
			"━━━━━━━━━━━━━━━━━━\n"+
			"📍 *DETECTED LOCATION*\n"+
			"├ *City* : `%s`\n"+
			"├ *District* : `%s`\n"+
			"└ *Coords* : `%s, %s`\n"+
			"━━━━━━━━━━━━━━━━━━\n"+
			"📡 *NETWORK PROFILE*\n"+
			"├ *ISP* : `%s`\n"+
			"├ *Public IP* : `%s`\n"+
			"└ *Test Node* : `%s (%s)`\n"+
			"━━━━━━━━━━━━━━━━━━\n"+
			"📊 *DIAGNOSTIC RESULTS*\n"+
			"📥 *Download* : `%.2f Mbps`\n"+
			"📤 *Upload* : `%.2f Mbps`\n"+
			"⚡ *Latency* : `%s` (Jitter: %s)\n"+
			"━━━━━━━━━━━━━━━━━━\n"+
			"%s"+
			"🕒 *REPORT TIME*\n`%s`",
		statusHeader,
		city, district, user.Lat, user.Lon,
		user.Isp, user.IP, s.Name, s.Country,
		downloadMbps, uploadMbps, latency, s.Jitter,
		alertMsg,
		currentTime,
	)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Failed to initialize bot: ", err)
	}

	log.Printf("🚀 Network Guard Active: %s", bot.Self.UserName)

	commands := []tgbotapi.BotCommand{
		{Command: "check", Description: "Run a speedtest now"},
		{Command: "startbot", Description: "Enable automatic reporting"},
		{Command: "stopbot", Description: "Disable automatic reporting"},
		{Command: "menu", Description: "Show help menu"},
	}
	bot.Request(tgbotapi.NewSetMyCommands(commands...))

	go func() {
		for {
			mu.RLock()
			active := autoReportActive
			mu.RUnlock()

			if active {
				report := runSpeedTest()
				msg := tgbotapi.NewMessage(int64(chatID), "📢 *AUTOMATED STATUS REPORT*\n\n"+report)
				msg.ParseMode = "Markdown"
				bot.Send(msg)
				log.Printf("✅ Automated report sent: %s", getWIBTime())
			} else {
				log.Println("💤 Auto-report is currently DISABLED.")
			}
			time.Sleep(autoInterval)
		}
	}()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {

			case "check":
				bot.Send(tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping))
				waitMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ *Analyzing network, please wait...*")
				waitMsg.ParseMode = "Markdown"
				sentWait, _ := bot.Send(waitMsg)

				report := runSpeedTest()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📑 *MANUAL DIAGNOSTIC REPORT*\n\n"+report)
				msg.ParseMode = "Markdown"
				bot.Send(msg)

				del := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentWait.MessageID)
				bot.Request(del)

			case "stopbot":
				mu.Lock()
				autoReportActive = false
				mu.Unlock()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🛑 *Automatic reporting has been DISABLED.*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "startbot":
				mu.Lock()
				autoReportActive = true
				mu.Unlock()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "▶️ *Automatic reporting has been ENABLED.*")
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "menu":
				helpText := "🛠️ *NETWORK GUARD INTERFACE*\n\n" +
					"🚀 `/check` - Run a manual speedtest.\n" +
					"▶️ `/startbot` - Enable 30m auto-reports.\n" +
					"🛑 `/stopbot` - Disable auto-reports.\n" +
					"📋 `/menu` - Show this help menu."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Unknown command.")
				bot.Send(msg)
			}
		}
	}
}
