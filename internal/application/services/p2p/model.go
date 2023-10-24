package p2p

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type User struct {
	tgUser tgbotapi.User
	IsNew  bool
	IDStr  string
}

type promoCode struct {
	Code      string
	Type      string
	Unlimited bool
	Expected  time.Time
}

func newPromoCode(code, pType string) promoCode {
	return promoCode{
		Code:      code,
		Type:      pType,
		Unlimited: true,
	}
}

type promoCodes map[string]promoCode

func addPromoCodes(items ...promoCode) promoCodes {
	newList := make(map[string]promoCode)
	for _, item := range items {
		newList[item.Code] = item
	}
	return newList
}

func (pcs promoCodes) exist(code string) bool {
	_, exist := pcs[code]
	return exist
}

func (pcs promoCodes) get(code string) (promoCode, bool) {
	pc, exist := pcs[code]
	return pc, exist
}
