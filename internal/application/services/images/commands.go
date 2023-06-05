package images

import (
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

//TODO: add images buffer, imageMergeVertical, imageMergeHorizontal,

func (d data) help(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	commandsList := "Image processing\n"
	commandsList += "Commands:\n"
	for _, commandsItem := range d.Commands() {
		if commandsItem.ListExclude {
			continue
		}
		commandsList += commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgCommands.Simple(msg.Chat.ID, commandsList)
}

func (d data) resize(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.WaitingWithText(msg.Chat.ID, "Send image, use text commands format '300'", "resizeImage")
}

func (d data) resizeImage(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	var err error
	var size = 100
	if msg.Photo == nil {
		return tgCommands.Simple(msg.Chat.ID, "image is empty or incorrect")
	}
	fileId := ""
	for _, photoItem := range msg.Photo {
		fileId = photoItem.FileID
	}
	buf, err := getTgFile(fileId)
	if err != nil {
		log.Println(err.Error())
		return tgCommands.Simple(msg.Chat.ID, err.Error())
	}
	if msg.Caption != "" {
		sizeNew, err := strconv.Atoi(msg.Caption)
		if err == nil {
			size = sizeNew
		}
	}
	fmt.Println("resize image with size", size)
	newImage, err := getMagic(buf.Bytes(), size)
	tgNewfile := tgbotapi.FileBytes{
		Name:  "photo.jpg",
		Bytes: newImage,
	}
	return tgCommands.PreparedCommand(tgbotapi.NewPhoto(msg.Chat.ID, tgNewfile))
}

func (d data) rotate(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	return tgCommands.WaitingWithText(msg.Chat.ID, "Send image, use text commands format '90' or '180'", "rotateImage")
}

func (d data) rotateImage(msg *tgbotapi.Message, commandName string, param string, params []string) tgCommands.HandlerResult {
	var err error
	var degrees float64 = 90
	if msg.Photo == nil {
		return tgCommands.Simple(msg.Chat.ID, "image is empty or incorrect")
	}
	fileId := ""
	for _, photoItem := range msg.Photo {
		fileId = photoItem.FileID
	}
	buf, err := getTgFile(fileId)
	if err != nil {
		log.Println(err.Error())
		return tgCommands.Simple(msg.Chat.ID, err.Error())
	}
	if msg.Caption != "" {
		degreesNew, err := strconv.ParseFloat(msg.Caption, 64)
		if err == nil {
			degrees = degreesNew
		}
	}
	fmt.Println("rotate image with degrees", degrees)
	newImage, err := rotate(buf.Bytes(), degrees)
	tgNewfile := tgbotapi.FileBytes{
		Name:  "rotated.jpg",
		Bytes: newImage,
	}
	return tgCommands.PreparedCommand(tgbotapi.NewPhoto(msg.Chat.ID, tgNewfile))
}