package tgModel

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SentMessages chan<- tgbotapi.Chattable

const (
	BotNameParam       BotParamRequest = "BotName"
	BotNAdminIdParam   BotParamRequest = "AdminId"
	BotAdminLoginParam BotParamRequest = "AdminLogin"
)

type BotParamRequest string
type BotParamResponse struct {
	StrVal   string
	Int64Val int64
	IntVal   int
	BoolVal  bool
	NotFound bool
	Error    error
}

type ParamHandlerFunc func(BotParamRequest) BotParamResponse

func (bpr BotParamResponse) Str() string {
	return bpr.StrVal
}

func (bpr BotParamResponse) Int() int {
	return bpr.IntVal
}

func (bpr BotParamResponse) Int64() int64 {
	return bpr.Int64Val
}

func BotParamStr(param string) BotParamResponse {
	return BotParamResponse{
		StrVal: param,
	}
}

func BotParamInt(param int) BotParamResponse {
	return BotParamResponse{
		IntVal: param,
	}
}

func BotParamInt64(param int64) BotParamResponse {
	return BotParamResponse{
		Int64Val: param,
	}
}

func BotParamBool(param bool) BotParamResponse {
	return BotParamResponse{
		BoolVal: param,
	}
}

func BotParamNotFound() BotParamResponse {
	return BotParamResponse{
		NotFound: true,
	}
}
