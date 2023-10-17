package tgModel

import "bytes"

var PhotoMediaTypes = []string{
	"image/png", "image/jpg", "image/jpeg",
	"image/pjpeg", "image/svg+xml",
	"image/tiff", "image/vnd.microsoft.icon",
	"image/icon", "image/webp",
	"photo", //tg type
}

type FileHandlerFunc func(FileCallbackData) (*bytes.Buffer, error)

type FileCallback interface {
	GetFile(FileCallbackData)
}

type FileCallbackData struct {
	FileID  string `json:"file_id"`
	FileUID string `json:"file_unique_id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Size    int    `json:"file_size"`
}

type FileTypes []string

func (ft FileTypes) Has(fType string) bool {
	for _, item := range ft {
		if item == fType {
			return true
		}
	}
	return false
}

type TgFileInfo struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		FileId       string `json:"file_id,omitempty"`
		FileUniqueId string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	} `json:"result,omitempty"`
}
