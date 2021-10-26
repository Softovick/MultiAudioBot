package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xfrr/goffmpeg/transcoder"
	telebot "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var voiceFiles = make([]string, 0) //имена файлов из папки voices
	//используем LookupEnv, чтобы определить видим ли мы переменную окружения
	tokenTelegram, existsToken := os.LookupEnv("TOKEN_TELEGRAM")
	if !existsToken {
		log.Fatal("TOKEN_TELEGRAM not found")
		return
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  tokenTelegram,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/hello", func(msg *telebot.Message) {
		bot.Send(msg.Sender, "Hello friend")
	})

	bot.Handle(telebot.OnVoice, func(msg *telebot.Message) {
		//по умолчанию телеграмм обрабатывает аудио в формате ogg
		voiceFile := "voices\\" + msg.Voice.UniqueID + ".ogg"
		bot.Download(&msg.Voice.File, voiceFile)
		voiceFiles = append(voiceFiles, voiceFile)
		fmt.Printf("%#v\n", voiceFiles)
		bot.Send(msg.Sender, "Get voice, and save to folder 'voices'")
	})

	bot.Handle("/merge", func(msg *telebot.Message) {
		count := len(voiceFiles)
		if count == 0 {
			bot.Send(msg.Sender, "Обрабатывать нечего")
			return
		}
		transCoder := new(transcoder.Transcoder)
		//входящая строка подсмотрена из решения для утилиты ffmpeg
		inputSource := "concat:"
		//"concat:01.ogg|02.ogg|03.ogg|04.ogg|05.ogg|07.ogg|07.ogg"
		for i := 0; i < len(voiceFiles); i++ {
			inputSource += voiceFiles[i] + "|"
		}
		inputSource = inputSource[:len(inputSource)-1]
		//inputSource := voiceFiles[0]
		outputFile := hex.EncodeToString([]byte(time.Stamp)) + ".ogg"
		//выводим информацию в консоль еще и потому, чтобы дать время заполнить данными переменные для горутин
		fmt.Printf("Из файла %s получаем файл %s", inputSource, outputFile)
		err := transCoder.Initialize(inputSource, outputFile)
		if err != nil {
			log.Fatal(err)
			return
		}
		//скорее всего .Run - goroutune
		//при использовании progress=true в случае ошибки получил развернутое описание
		done := transCoder.Run(true)
		progress := transCoder.Output()
		for message := range progress {
			fmt.Println(message)
		}
		err = <-done
		if err != nil {
			log.Fatal(err)
			return
		}
		bot.Send(msg.Sender, "Файл обработан")
		fileVoice := &telebot.Voice{File: telebot.FromDisk(outputFile)}
		bot.Send(msg.Sender, fileVoice)
		os.Remove(outputFile)
		for i := 0; i < len(voiceFiles); i++ {
			os.Remove(voiceFiles[i])
		}
		voiceFiles = voiceFiles[:0]
	})

	bot.Start()
}
