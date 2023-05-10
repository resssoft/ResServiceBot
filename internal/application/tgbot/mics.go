package tgbot

import "strings"

func splitCommand(command string, separate string) ([]string, string) {
	if command == "" {
		return []string{}, ""
	}
	if separate == "" {
		separate = " "
	}
	result := strings.Split(command, separate)
	return result, strings.Replace(command, result[0]+separate, "", -1)
}
