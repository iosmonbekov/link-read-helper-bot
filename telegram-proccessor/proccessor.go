package telegram_proccessor

import (
	"fmt"
	"link-saver-bot/client"
	"link-saver-bot/storage"
	errorstorage "link-saver-bot/storage/error-storage"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

type TelegramProccessor struct {
	Client  client.Client
	Storage storage.Storage
}

func New(client client.Client, storage storage.Storage) TelegramProccessor {
	return TelegramProccessor{
		Client:  client,
		Storage: storage,
	}
}

func (p *TelegramProccessor) Start() error {
	offset := 0
	limit := 100

	errorStorage := errorstorage.New("error-loger", 20)

	for {
		updates, err := p.Client.FetchUpdates(offset, limit)
		if err != nil {
			fmt.Printf("telegram_proccessor.Start: error on client.FetchUpdates -> %v\n", err)
			continue
		}

		if len(updates) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := p.Proccess(updates); err != nil {
			fmt.Printf("telegram_proccessor.Start: error on p.Proccess -> %v\n", err)

			errorStorage.Errors = append(errorStorage.Errors, err.Error())
			if len(errorStorage.Errors) >= int(errorStorage.Limit()) {
				errorStorage.Save([]byte(strings.Join(errorStorage.Errors, "\n")))
				fmt.Printf("Error count reached -> %v\n", errorStorage.Limit())
				os.Exit(1)
			}
			continue
		}

		errorStorage.Errors = []string{}
		offset = updates[len(updates)-1].UpdateId + 1
	}
}

func (p *TelegramProccessor) Proccess(updates []client.Update) error {
	for _, update := range updates {
		if update.Message.Text == "" {
			if err := p.Client.SendMessage(update, UNKNOWN_MESSAGE); err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
			}
			continue
		}

		if checkURL(update.Message.Text) {
			fmt.Println("<<"+update.Message.From.Username, update.Message.Text+">>\n")

			fileContent, err := p.Storage.Load(update.Message.From.Username)
			if err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on storage.Load -> %v\n", err)
			}

			if strings.Contains(string(fileContent), update.Message.Text) {
				if err := p.Client.SendMessage(update, LINK_ALREADY_SAVED); err != nil {
					return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
				}
				continue
			}

			fileContent = append(fileContent, []byte(update.Message.Text+"\n")...)

			if err := p.Storage.Save(update.Message.From.Username, fileContent); err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on storage.Save -> %v\n", err)
			}

			if err := p.Client.SendMessage(update, SUCCESS_LINK_SAVE); err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
			}

			continue
		}

		switch update.Message.Text {
		case "/start":
			if err := p.Client.SendMessage(update, HELLO_MESSAGE); err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
			}
		case "/help":
			if err := p.Client.SendMessage(update, HELP_MESSAGE); err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
			}
		case "/random":
			fileContent, err := p.Storage.Load(update.Message.From.Username)
			if err != nil {
				return fmt.Errorf("telegram_proccessor.Proccess: error on storage.Load -> %v\n", err)
			}

			if len(fileContent) == 0 {
				if err := p.Client.SendMessage(update, NO_LINK_MESSAGE); err != nil {
					return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
				}
			} else {
				links := strings.Split(string(fileContent), "\n")
				rand.Seed(time.Now().UnixNano())
				randomIndex := rand.Intn(len(links))

				toSendLink := links[randomIndex]
				links = append(links[0:randomIndex], links[randomIndex+1:]...)
				fileContent := []byte(strings.Join(links, "\n"))

				if err := p.Client.SendMessage(update, toSendLink); err != nil {
					return fmt.Errorf("telegram_proccessor.Proccess: error on client.SendMessage -> %v\n", err)
				}
				if err := p.Storage.Remove(update.Message.From.Username); err != nil {
					return fmt.Errorf("telegram_proccessor.Proccess: error on storage.Remove -> %v\n", err)
				}
				if err := p.Storage.Save(update.Message.From.Username, fileContent); err != nil {
					return fmt.Errorf("telegram_proccessor.Proccess: error on storage.Save -> %v\n", err)
				}
			}
		default:
			if err := p.Client.SendMessage(update, UNKNOWN_MESSAGE); err != nil {
				return fmt.Errorf("Error occured when sending message: %v", err)
			}
		}
	}

	return nil
}

func checkURL(text string) bool {
	_, err := url.ParseRequestURI(text)
	if err != nil {
		return false
	}
	return true
}
