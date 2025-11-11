package data

import (
	"log"
	"os"
	"strings"
)

func GetSteamID() (string, error) {
	fileText, err := os.ReadFile("data/steamid.txt")
	if err != nil {
		return "", err
	} else {
		return string(fileText), nil
	}
}

func WriteSteamID(steamid string) {
	steamid = strings.TrimSpace(steamid)
	file, err := os.Create("data/steamid.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	_, err = file.WriteString(steamid)
	if err != nil {
		log.Fatalln(err)
	}
}
