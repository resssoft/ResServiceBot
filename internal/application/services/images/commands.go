package images

import (
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

//TODO: add images buffer, imageMergeVertical, imageMergeHorizontal,

func (d data) help(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	commandsList := "Image processing\n"
	commandsList += "Commands:\n"
	for _, commandsItem := range d.Commands() {
		if commandsItem.ListExclude {
			continue
		}
		commandsList += "/" + commandsItem.Command + " - " + commandsItem.Description + "\n"
	}
	return tgModel.Simple(msg.Chat.ID, commandsList)
}

func (d data) resize(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.DeferredWithText(msg.Chat.ID, "Send image, use text commands format '300'", "resizeImage", "", nil)
}

func (d data) resizeImage(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	var err error
	var size = 100
	if msg.Photo == nil {
		return tgModel.Simple(msg.Chat.ID, "image is empty or incorrect")
	}
	var fileData tgModel.FileCallbackData
	for _, photoItem := range msg.Photo {
		fileData.FileID = photoItem.FileID
		fileData.FileUID = photoItem.FileUniqueID
		fileData.Size = photoItem.FileSize
		fileData.Height = photoItem.Height
		fileData.Width = photoItem.Width
	}
	buf, err := command.FilesCallback(fileData)
	if err != nil {
		log.Println(err.Error())
		return tgModel.Simple(msg.Chat.ID, err.Error())
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
	return tgModel.PreparedCommand(tgbotapi.NewPhoto(msg.Chat.ID, tgNewfile))
}

func (d data) rotate(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	return tgModel.DeferredWithText(msg.Chat.ID, "Send image, use text commands format '90' or '180'", "rotateImage", "", nil)
}

func (d data) rotateImage(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	//fmt.Println("start rotateImage by service", msg)
	//zlog.Info().Any("msg", msg).Any("command", command).Send()
	var err error
	var degrees float64 = 90
	if msg.Photo == nil {
		return tgModel.Simple(msg.Chat.ID, "image is empty or incorrect")
	}
	var fileData tgModel.FileCallbackData
	for _, photoItem := range msg.Photo {
		fileData.FileID = photoItem.FileID
		fileData.FileUID = photoItem.FileUniqueID
		fileData.Size = photoItem.FileSize
		fileData.Height = photoItem.Height
		fileData.Width = photoItem.Width
	}
	buf, err := command.FilesCallback(fileData)

	if err != nil {
		log.Println(err.Error())
		return tgModel.Simple(msg.Chat.ID, err.Error())
	}
	if msg.Caption != "" {
		degreesNew, err := strconv.ParseFloat(msg.Caption, 64)
		if err == nil {
			degrees = degreesNew
		}
	}
	dataBytes := buf.Bytes()
	fmt.Println("rotate image with degrees", degrees, len(dataBytes))
	newImage, err := rotate(dataBytes, degrees)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, err.Error())
	}
	tgNewfile := tgbotapi.FileBytes{
		Name:  "rotated.jpg",
		Bytes: newImage,
	}
	return tgModel.PreparedCommand(tgbotapi.NewPhoto(msg.Chat.ID, tgNewfile))
}
