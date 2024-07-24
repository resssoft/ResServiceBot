package markovka

type Params struct {
	DataPath string
	BotName  string
	Chats    map[int64]ChatData
}

type ChatData struct {
	Id    int64
	Param int
}

type States struct {
	States []State
	Name   string
}

// Markov model
const noData = "word no exist on the memory"

type Markov struct{}

type State struct {
	Word       string
	Count      int
	Prob       float64
	NextStates []State
}

var markov Markov
