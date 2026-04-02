package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/showwin/speedtest-go/speedtest"
)

const (
	botToken     = "8680477535:AAGoZjz6DB9nIt_FZdICFOOAvB12xjFt0Ag"
	chatID       = -5275988185
	locationName = "PT DIGIMAX DAKSA NINRAKARA - Podomoro Golf View"
	autoInterval = 30 * time.Minute 
	pingLimit    = 100 * time.Millisecond
	minDownload  = 10.0
)

func runSpeedTest() string {
	log.Println("⏳ Starting Speedtest...")
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return "🚨 *SYSTEM ERROR*: Target ISP unreachable."
	}

	serverList, _ := speedtest.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	if len(targets) == 0 {
		return "🚨 *SYSTEM ERROR*: No responsive test servers found."
	}

	s := targets[0]
	s.PingTest(nil)
	s.DownloadTest()
	s.UploadTest()

	downloadMbps := s.DLSpeed / 1000000
	uploadMbps := s.ULSpeed / 1000000
	latency := s.Latency.Round(time.Millisecond)

	isAnomalous := latency > pingLimit || downloadMbps < minDownload
	statusHeader := "✅ SYSTEM MONITORING INET - OFFICE"
	alertMessage := ""

	if isAnomalous {
		statusHeader = "🚨 ANOMALY DETECTED"
		alertMessage = "\n⚠️ *CRITICAL ALERT*:\n- High Latency or Low Bandwidth detected.\n- Check physical cables/fiber optics.\n- Restart Gateway/ONT if necessary."
	}

	return fmt.Sprintf(
		"🌐 *%s*\n"+
			"───────────────────\n"+
			"📍 *LOCATION* : `%s`\n"+
			"📡 *NETWORK PROFILE*\n"+
			"• *ISP* : `%s`\n"+
			"• *Node* : `%s (%s)`\n"+
			"───────────────────\n"+
			"📊 *PERFORMANCE METRICS*\n"+
			"• *Download* : `%.2f Mbps`\n"+
			"• *Upload* : `%.2f Mbps`\n"+
			"• *Latency* : `%s`\n"+
			"───────────────────\n"+
			"🕒 *TIMESTAMP* : `%s`\n%s",
		statusHeader, locationName, user.Isp, s.Name, s.Country,
		downloadMbps, uploadMbps, latency,
		time.Now().Format("Jan 02, 2006 | 15:04:05"), alertMessage,
	)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Bot initialization failed: ", err)
	}

	log.Printf("🚀 Network Guard Active: %s", bot.Self.UserName)

	go func() {
		for {
			report := runSpeedTest()
			msg := tgbotapi.NewMessage(int64(chatID), "🕒 *SCHEDULED REPORT*\n\n"+report)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
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

		if update.Message.IsCommand() && update.Message.Command() == "cek" {
			log.Printf("📩 Manual request from: %s", update.Message.From.UserName)
			
			waitMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "⏳ *Testing network speed...* Please wait.")
			waitMsg.ParseMode = "Markdown"
			sentWait, _ := bot.Send(waitMsg)

			report := runSpeedTest()

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "📊 *MANUAL CHECK RESULT*\n\n"+report)
			msg.ParseMode = "Markdown"
			bot.Send(msg)

			del := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, sentWait.MessageID)
			bot.Request(del)
		}
	}
}