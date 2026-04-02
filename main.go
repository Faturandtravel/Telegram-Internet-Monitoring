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
	locationName = "PT DIGIMAX DAKSA NINRAKARA -  Podomoro Golf View, Jl. Mochamad Thohir No.10 Ruko Granada B3, Bojong Nangka, Kec. Gn. Putri, Kabupaten Bogor, Jawa Barat 16953"
	interval     = 30 * time.Second
	pingLimit    = 100 * time.Millisecond
	minDownload  = 10.0 
)

func runSpeedTest() string {
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
	
	statusHeader := "✅ SYSTEM STABLE"
	alertMessage := ""

	if isAnomalous {
		statusHeader = "🚨 ANOMALY DETECTED"
		alertMessage = "\n⚠️ *CRITICAL ALERT*:\n- High Latency or Low Bandwidth detected.\n- Please check Physical Cables/Fiber Optic.\n- Consider restarting the Gateway/ONT."
	}

	report := fmt.Sprintf(
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
			"🕒 *TIMESTAMP* : `%s`\n"+
			"%s",
		statusHeader,
		locationName,
		user.Isp,
		s.Name, s.Country,
		downloadMbps,
		uploadMbps,
		latency,
		time.Now().Format("Jan 02, 2006 | 15:04:05"),
		alertMessage,
	)

	return report
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Failed to initialize bot: ", err)
	}

	log.Printf("🚀 Network Guard active for location: %s", locationName)

	for {
		log.Println("🔄 Running speedtest for " + locationName)
		report := runSpeedTest()

		msg := tgbotapi.NewMessage(int64(chatID), report)
		msg.ParseMode = "Markdown"
		
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("❌ Failed to send Telegram message: %v", err)
		} else {
			log.Println("✅ Report delivered for " + locationName)
		}

		time.Sleep(interval)
	}
}