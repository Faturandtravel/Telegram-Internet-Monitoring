package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/showwin/speedtest-go/speedtest"
)

const (
	botToken     = "8680477535:AAGoZjz6DB9nIt_FZdICFOOAvB12xjFt0Ag"
	chatID       = -5275988185
	autoInterval = 30 * time.Minute

	pingLimit   = 50 * time.Millisecond
	minDownload = 60.0
)

func getWIBTime() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now().Add(7*time.Hour).Format("02 Jan 2006 | 15:04:05") + " WIB"
	}
	return time.Now().In(loc).Format("02 Jan 2006 | 15:04:05") + " WIB"
}

// getLocationAddress melakukan reverse geocoding dari koordinat ke nama kota dan kecamatan
func getLocationAddress(lat, lon string) (kota, kecamatan string) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%s&lon=%s&zoom=12&accept-language=id", lat, lon)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "Tidak Diketahui", "Tidak Diketahui"
	}
	req.Header.Set("User-Agent", "inetmonitor/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "Tidak Diketahui", "Tidak Diketahui"
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
		return "Tidak Diketahui", "Tidak Diketahui"
	}

	// Tentukan nama kota (prioritas: city > town > state_district)
	kota = result.Address.City
	if kota == "" {
		kota = result.Address.Town
	}
	if kota == "" {
		kota = result.Address.StateDistrict
	}
	if kota == "" {
		kota = result.Address.State
	}
	if kota == "" {
		kota = "Tidak Diketahui"
	}

	// Tentukan nama kecamatan (prioritas: district > suburb > village)
	kecamatan = result.Address.District
	if kecamatan == "" {
		kecamatan = result.Address.Suburb
	}
	if kecamatan == "" {
		kecamatan = result.Address.Village
	}
	if kecamatan == "" {
		kecamatan = result.Address.Regency
	}
	if kecamatan == "" {
		kecamatan = "Tidak Diketahui"
	}

	return kota, kecamatan
}

func runSpeedTest() string {
	log.Println("⏳ Starting Real-Time Diagnostics (High Accuracy)...")

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

	kota, kecamatan := getLocationAddress(user.Lat, user.Lon)
	currentTime := getWIBTime()

	return fmt.Sprintf(
		"🌐 *%s*\n"+
			"━━━━━━━━━━━━━━━━━━\n"+
			"📍 *DETECTED LOCATION*\n"+
			"├ *Kota* : `%s`\n"+
			"└ *Kecamatan* : `%s`\n"+
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
		kota, kecamatan,
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
		{Command: "cekdong", Description: "Run a speedtest now"},
		{Command: "daftarmenu", Description: "Show help menu"},
	}
	bot.Request(tgbotapi.NewSetMyCommands(commands...))

	go func() {
		for {
			report := runSpeedTest()
			msg := tgbotapi.NewMessage(int64(chatID), "📢 *AUTOMATED STATUS REPORT*\n\n"+report)
			msg.ParseMode = "Markdown"
			bot.Send(msg)

			log.Printf("✅ Automated report sent: %s", getWIBTime())
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

			case "cekdong":
				log.Printf("📩 Manual request from: %s", update.Message.From.UserName)

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

			case "daftarmenu":
				helpText := "🛠️ *NETWORK GUARD INTERFACE*\n\n" +
					"🚀 `/cekdong` - Run a manual speedtest.\n" +
					"📋 `/daftarmenu` - Show this help menu.\n\n" +
					"💡 _System automatically reports every 30 minutes._"

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
