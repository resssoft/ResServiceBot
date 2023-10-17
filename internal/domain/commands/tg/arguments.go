package tgModel

import "strings"

type CommandArguments struct {
	Raw    string
	Parsed bool
	List   []string
}

func (ca *CommandArguments) Parse() []string {
	if ca.Raw == "" {
		return nil
	}
	if ca.Parsed {
		return ca.List
	}
	ca.List = strings.Split(ca.Raw, " ")
	ca.Parsed = true
	return ca.List
}
