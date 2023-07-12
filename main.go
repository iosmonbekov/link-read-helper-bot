package main

import (
	"flag"
	"fmt"
	telegram_client "link-saver-bot/client/telegram-client"
	sqlitestorage "link-saver-bot/storage/sqlite-storage"
	telegram_proccessor "link-saver-bot/telegram-proccessor"
	"log"
	"os"
)

func main() {
	tgClient := telegram_client.New("https://api.telegram.org", mustToken())
	// tgStorage, err := telegram_storage.New("link-saver-bot")
	// if err != nil {
	// 	fmt.Printf("main: error on telegram_storage.New -> %v\n", err)
	// 	os.Exit(1)
	// }
	tgStorage, err := sqlitestorage.New("link-saver-bot-sql-storage.db")
	if err != nil {
		fmt.Printf("main: error on telegram_storage.New -> %v\n", err)
		os.Exit(1)
	}
	tgStorage.Init()

	proccessor := telegram_proccessor.New(tgClient, tgStorage)

	if err := proccessor.Start(); err != nil {
		fmt.Printf("main: error on proccessor.Start -> %v\n", err)
	}
}

func mustToken() string {
	token := flag.String("token", "", "Token for accessing telegram bot api")
	flag.Parse()

	if *token == "" {
		log.Panic("Token not provided as flag when running the program")
	}
	return *token
}
