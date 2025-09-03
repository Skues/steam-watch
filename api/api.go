package api

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const mySteamID string = "76561198198664702"
const CFRSteamID string = "76561198082191202"

type FriendListResponse struct {
	FriendListResponse FriendList `json:"friendslist"`
}
type FriendList struct {
	FriendList []Friend `json:"friends"`
}

type Friend struct {
	FriendSteamID string `json:"steamid"`
	Relationship  string `json:"relationship"`
	FriendSince   int64  `json:"friend_since"`
}

type RecentGames struct {
	RecentGamesResponse RecentGamesResult `json:"response"`
}
type RecentGamesResult struct {
	TotalCount int    `json:"total_count"`
	Games      []Game `json:"games"`
}
type Game struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Week   int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIcon         string `json:"img_icon_url"`
}

type PlayerSummary struct {
	PlayerSummaryResponse PlayerList `json:"response"`
}
type PlayerList struct {
	Players []Player `json:"players"`
}

type Player struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	AvatarHash               string `json:"avatarhash"`
	LastLogoff               int64  `json:"lastlogoff"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            string `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LocCountryCode           string `json:"loccountrycode"`
}

var key string = os.Getenv("STEAM_API_KEY")

func GetInfo() {
	list := friendList(mySteamID)
	fmt.Println(list)
	id := "76561198082191202"
	playerInfo := playerSummary(id)
	fmt.Println(playerInfo.PlayerSummaryResponse.Players[0].PersonaName)
	fmt.Printf("Last logged in on: %v", unixToTime(playerInfo.PlayerSummaryResponse.Players[0].LastLogoff))

	result := recentlyPlayed(id)
	// fmt.Println(result)
	for i, game := range result.RecentGamesResponse.Games {
		fmt.Printf("\n-------\nID: %v\n%s\nPast 2 weeks: %v hours\nTotal Playtime: %v hours", i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)

	}

}

func unixToTime(unix int64) (timeReturn time.Time) {
	loc, _ := time.LoadLocation("Europe/London")
	timeReturn = time.Unix(unix, 0).In(loc)

	return timeReturn
}

func playerSummary(steamid string) PlayerSummary {
	var result PlayerSummary
	url := "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=" + key + "&steamids=" + steamid
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		panic(err)
	}

	return result

}

func recentlyPlayed(steamid string) RecentGames {
	var result RecentGames
	url := "http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=" + key + "&steamid=" + steamid + "&format=json"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		panic(err)
	}
	return result

}

func friendList(steamid string) FriendListResponse {
	var result FriendListResponse
	url := "http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=" + key + "&steamid=" + steamid + "&relationship=friend"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	json.Unmarshal(bodyBytes, &result)
	return result

}
