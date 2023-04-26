package main

type TGUser struct {
	UserID  int64
	ChatId  int64
	Login   string
	Name    string
	IsAdmin bool
}

type KeyBoardTG struct {
	Rows []KeyBoardRowTG
}

type KeyBoardRowTG struct {
	Buttons []KeyBoardButtonTG
}

type KeyBoardButtonTG struct {
	Text string
	Data string
}

type ChatUser struct {
	ChatId      int64
	ChatName    string
	ContentType string
	CustomRole  string
	VoteCount   int
	User        TGUser
}
