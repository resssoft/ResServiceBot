package testManager

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "fun-coice/pkg/telegram-bot-api"
)

func (d *data) add(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	newTest := &testParent{
		Name:     "",
		Code:     fmt.Sprintf("%v", msg.Chat.ID),
		Children: nil,
	}
	d.setTest(newTest)
	d.setActiveTest(msg.From.ID, newTest.Code)
	return tgModel.DeferredWithText(msg.Chat.ID, "Send test code (for code='my_code', will be created command /test_run_my_code ) and Name test on the next line", "event:test_set_name", "", nil)
}

func (d *data) setName(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	testData := strings.Split(msg.Text, "\n")
	if len(testData) < 2 {
		return tgModel.DeferredWithText(msg.Chat.ID, "Incorrect! Set code and name on the different lines", "event:test_set_name", "", nil)
	}
	newTest := &testParent{
		Name:     testData[1],
		Code:     testData[0],
		Children: nil,
	}
	d.setTest(newTest)
	d.setActiveTest(msg.From.ID, newTest.Code)
	return tgModel.DeferredWithText(msg.Chat.ID, "Write question and answer, new line for split. For done send: /test_done", "event:test_append", "", nil)
}

func (d *data) append(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	if msg.Command() == "test_done" {
		return d.done(msg, command)
	}
	activeTest, err := d.getActiveTest(msg.From.ID)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get active test error:"+err.Error())
	}
	err = d.appendTestRows(msg.Text, activeTest)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "Please send 2 or more lines for question and one answer for minimal: "+err.Error())
	}

	return tgModel.DeferredWithText(msg.Chat.ID, "Write question and answer, new line for split. Right answers can be more than one. For done send: /test_done", "event:test_append", "", nil)
}

func (d *data) done(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	activeTest, err := d.getActiveTest(msg.From.ID)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get active test error:"+err.Error())
	}
	d.delActiveTest(msg.Chat.ID)
	return tgModel.Simple(msg.Chat.ID, "Test done. For run used: /test_run_"+activeTest.Code)
}

func (d *data) run(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	foundedTest, err := d.getTest(strings.ReplaceAll(msg.Command(), "test_run_", ""))
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get active test error:"+err.Error())
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	children := foundedTest.Children
	r.Shuffle(len(children), func(i, j int) { children[i], children[j] = children[j], children[i] })
	foundedTest.Children = children
	firstQuestion := foundedTest.Children[0]
	firstQuestion.Current = true
	d.setAUserTest(msg.Chat.ID, foundedTest)
	msg.Text = ""
	return tgModel.Simple(msg.Chat.ID, "Started test "+foundedTest.Name).WithRedirect("event:test_question", msg)
}

func (d *data) question(msg *tgbotapi.Message, command *tgModel.Command) *tgModel.HandlerResult {
	foundedTest, err := d.getUserTest(msg.Chat.ID)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get user test error:"+err.Error())
	}
	answerCount := 0
	for _, q := range foundedTest.Children {
		if q.Answer != "" {
			answerCount++
		}
		if q.Current {
			return tgModel.DeferredWithText(msg.Chat.ID, q.Question, "event:test_answer", "", nil)
		}
	}
	if answerCount == len(foundedTest.Children) {
		return d.ending(msg, command)
	}
	return d.ending(msg, command)
}

func (d *data) answer(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	foundedTest, err := d.getUserTest(msg.Chat.ID)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get user test error:"+err.Error())
	}
	msgResult := ""

	for i, q := range foundedTest.Children {
		if q.Current {
			msgResult = q.Incorrect
			q.Answer = msg.Text
			q.Current = false
			for _, answer := range q.Answers {
				if answer == msg.Text {
					msgResult = q.Correct
				}
			}
			foundedTest.Children[i] = q
		}
	}
	for i, q := range foundedTest.Children {
		if q.Answer == "" {
			q.Current = true
			foundedTest.Children[i] = q
			break
		}
	}
	return tgModel.Simple(msg.Chat.ID, msgResult).WithRedirect("event:test_question", msg)
}

func (d *data) ending(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	foundedTest, err := d.getUserTest(msg.Chat.ID)
	if err != nil {
		return tgModel.Simple(msg.Chat.ID, "get user test error:"+err.Error())
	}
	correctCount := 0.0
	for _, q := range foundedTest.Children {
		for _, answer := range q.Answers {
			if answer == q.Answer {
				correctCount++
			}
		}
	}
	msgResult := fmt.Sprintf("Done. Result: %.0f is correct (%.2f)",
		correctCount,
		correctCount*100.0/float64(len(foundedTest.Children)))
	d.delUserTest(msg.Chat.ID)
	return tgModel.Simple(msg.Chat.ID, msgResult)
}

func (d *data) importTest(msg *tgbotapi.Message, _ *tgModel.Command) *tgModel.HandlerResult {
	//TODO:implement
	return tgModel.Simple(msg.Chat.ID, "Format: \n"+
		"First line: code (for code='my_code', will be created command /test_run_my_code ) \n"+
		"Second line test name - it will be printed before start")
}
