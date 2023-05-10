package tgbot

type TgFileInfo struct {
	Ok     bool `json:"ok,omitempty"`
	Result struct {
		FileId       string `json:"file_id,omitempty"`
		FileUniqueId string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	} `json:"result,omitempty"`
}
