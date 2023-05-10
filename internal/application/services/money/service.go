package financy

import (
	"encoding/json"
	"fmt"
	tgCommands "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type data struct {
	list  tgCommands.Commands
	token string
}

const dbDateFormatMonth = "2006-01-02"

var (
	urlTmp = "https://api.apilayer.com/fixer/convert?to=%s&from=%s&amount=%s"
)

type converterResult struct {
	Result  float64 `json:"result,omitempty"`
	Success bool    `json:"success,omitempty"`
}

//TODO: rate cache (to 6h) + use rate.am (run events to update)
// templates "fiat", "convert", "конверт", "кон", "из", "from"
// https://rate.am/calculator/rates.ashx?cr1=USD&hcr=AMD&cr2=RUR&orgId=466fe84c-197f-4174-bc97-e1dc7960edc7&rtype=1&tp=0&l=lang3&r=
// https://github.com/zaikin-andrew/rate-am-extension/search?q=rate.am

func New(token string) tgCommands.Service {
	result := data{
		token: token,
	}
	commandsList := make(tgCommands.Commands)
	commandsList["fiat"] = tgCommands.Command{
		Command:     "/fiat",
		Synonyms:    []string{"сколько будет", "фиат"},
		Description: "Convert one fiat to others (usd, rub, amd)",
		CommandType: "text",
		Permissions: tgCommands.FreePerms,
		Handler:     result.fiat,
	}

	result.list = commandsList
	return &result
}

func (d data) Commands() tgCommands.Commands {
	return d.list
}

func (d data) fiat(msg *tgbotapi.Message, commandName string, param string, params []string) (tgbotapi.Chattable, bool) {
	convertFrom := "AMD"
	convertTo1 := "RUB"
	convertTo2 := "USD"
	if len(params) > 2 {
		switch params[2] {
		case "a", "amd", "am", "ам", "амд", "дпам", "драм", "др":
			convertFrom = "AMD"
			convertTo1 = "RUB"
			convertTo2 = "USD"
		case "r", "ru", "rub", "rur", "ру", "р", "руб", "рублей":
			convertFrom = "RUB"
			convertTo1 = "AMD"
			convertTo2 = "USD"
		case "s", "us", "usd", "$", "дол", "до", "доларов", "долларов":
			convertFrom = "USD"
			convertTo1 = "AMD"
			convertTo2 = "RUB"
		}
	}
	_, err := strconv.Atoi(params[1])
	msgText := "-"
	if err != nil {
		msgText = "digit err"
	} else {
		msgText = fmt.Sprintf("Convert result from %s %s = \n%s %s \n%s %s \n[%s]",
			params[1], convertFrom,
			d.fiatConvert(convertFrom, convertTo1, params[1]), convertTo1,
			d.fiatConvert(convertFrom, convertTo2, params[1]), convertTo2,
			time.Now().Format(dbDateFormatMonth))
	}

	msgNew := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	msgNew.ReplyToMessageID = msg.MessageID
	return msgNew, true
}

func (d data) fiatConvert(from, to, amount string) string {
	time.Sleep(time.Millisecond * 100)
	url := fmt.Sprintf(urlTmp, to, from, amount)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", d.token)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	var resultObj converterResult
	err = json.Unmarshal(body, &resultObj)
	if err != nil {
		fmt.Println(err)
	} else {
		return fmt.Sprintf("%.2f", resultObj.Result)
	}
	return ""
}
