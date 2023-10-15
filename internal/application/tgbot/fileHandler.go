package tgbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	zlog "github.com/rs/zerolog/log"
	"net/http"
)

func (d *Data) getTgFile(fileData tgModel.FileCallbackData) (*bytes.Buffer, error) {
	zlog.Info().Str("getTgFile with ID", fileData.FileID).Send()
	buf := new(bytes.Buffer)
	response, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", d.Token, fileData.FileID))
	if err != nil {
		return nil, errors.New("download TG photo error")
	}
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, errors.New("read file Body err")
	}
	result := buf.String()
	fileInfo := tgModel.TgFileInfo{}
	err = json.Unmarshal([]byte(result), &fileInfo)
	if err != nil {
		return nil, errors.New("decode fileInfo err")
	}
	fileUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", d.Token, fileInfo.Result.FilePath)
	response, err = http.Get(fileUrl)
	if err != nil {
		return nil, errors.New("download TG import file error")
	}
	buf = new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf, nil
}
