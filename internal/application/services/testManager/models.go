package testManager

type testParent struct {
	Name     string
	Code     string
	Children []*testChild
}

type testChild struct {
	Question  string
	Type      string
	Answers   []string
	Variants  []string
	Incorrect string
	Correct   string

	Answer  string
	Current bool
}
