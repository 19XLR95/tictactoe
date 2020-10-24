package main

import (
	"encoding/base64"
	"math/rand"
	"time"
)

func generateRandomString(len int) string {
	s := ""

	ns := rand.NewSource(time.Now().UnixNano())
	nr := rand.New(ns)

	for i := 0; i < len; i++ {
		n := nr.Intn(126-32) + 32

		s += string(n)
	}

	return s
}

func convertBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func tie(gameStateArr [][]int) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if gameStateArr[i][j] == 0 {
				return false
			}
		}
	}

	return true
}

func getEmptySlots(gameStateArr [][]int) ([]int, []int) {
	var rows []int
	var cols []int

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if gameStateArr[i][j] == 0 {
				rows = append(rows, i)
				cols = append(cols, j)
			}
		}
	}

	return rows, cols
}
