package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PriceResponse struct {
	Time       TimeData `json:"time"`
	Disclaimer string   `json:"disclaimer"`
	BPI        BPIData  `json:"bpi"`
}
type TimeData struct {
	Updated    string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
	UpdatedUK  string `json:"updateduk"`
}

type BPIData struct {
	USD Currency `json:"USD"`
	BTC Currency `json:"BTC"`
}

type Currency struct {
	Code        string  `json:"code"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if len(os.Getenv("TELEGRAM_APITOKEN")) < 16 {
		fmt.Println("Debe setear la variable de entorno TELEGRAM_APITOKEN con el siguiente comando")
		fmt.Println("(reemplazando <api_token> por el TOKEN provisto por Telegram)")
		fmt.Println("set TELEGRAM_APITOKEN=<api_token>")
	}
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Autorizado para la cuenta %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignorar actualizaciones que no sean mensajes
			continue
		}

		if !update.Message.IsCommand() { // ignorar mensajes que no sean comandos
			continue
		}

		// Crear un nuevo MessageConfig vacio
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extraer el comando del mensaje
		switch update.Message.Command() {
		case "ayuda":
			//TODO agregar todos los comandos validos al mensaje de ayuda
			msg.Text = "Los comandos validos son: /ayuda, /screenshot..."
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		case "screenshot":
			msg.Text = ""
			photo := tgbotapi.NewPhoto(update.Message.From.ID, tgbotapi.FilePath("img/screenshot.jpeg"))
			if _, err = bot.Send(photo); err != nil {
				log.Fatalln(err)
			}

			//TODO Completar el mensaje para este comando
		case "autor":
			msg.Text = "Este bot fue creado por ... (nombre y apellido del alumno)"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

			//TODO borrar este comando de ejemplo
		case "ejemplo":
			msg.Text = "Esta es la respuesta para un comando de ejemplo"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

			//TODO agregar 5 comandos nuevos

		case "bitcoin":
			msg.Text = GetBitcoinPrice()
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		default:
			msg.Text = "No entiendo ese comando"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}

	}
}

func GetBitcoinPrice() string {
	url := "https://api.coindesk.com/v1/bpi/currentprice/BTC.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al hacer la solicitud HTTP:", err)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error al leer la respuesta HTTP:", err)
		return ""
	}
	var priceData PriceResponse
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		fmt.Println("Error al decodificar el JSON:", err)
		return ""
	}
	btcPrice := priceData.BPI.USD.RateFloat
	resp.Body.Close()
	result := fmt.Sprint("USD ", btcPrice)
	return result
}
