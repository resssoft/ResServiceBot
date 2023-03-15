package main

import (
	"fmt"
	"github.com/mnogu/go-calculator"
	"log"
	"math"
)

func calcFromStr(data string) string {
	log.Println("calcFromStr", data)
	val, err := calculator.Calculate(data)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	resultText := fmt.Sprintf("%.2f", val)
	intPart, floatPart := math.Modf(val)
	if floatPart == 0 {
		resultText = fmt.Sprintf("%.0f", intPart)
	}
	if val < 0.01 {
		resultText = fmt.Sprintf("%.5f", val)
	}
	return resultText
}
