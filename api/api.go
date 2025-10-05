package api

// https://developer.valvesoftware.com/wiki/Steam_Web_API#GetPlayerSummaries_.28v0001.29
import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
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
	GamesCount int    `json:"game_count"`
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
	GameID                   string `json:"gameid"`
	GameExtraInfo            string `json:"gameextrainfo"`
	LocCountryCode           string `json:"loccountrycode"`
	LocCityID                string `json:"loccityid"`
}

type Result struct {
	Index   int
	Summary PlayerSummary
	Recent  RecentGames
}

var apiCalls = map[string]map[string]string{"GetFriendList": {"API Type": "ISteamUser", "version": "v0001", "steamid": "steamid"}, "GetPlayerSummaries": {"API Type": "ISteamUser", "version": "v0002", "steamid": "steamids"}, "GetOwnedGames": {"API Type": "IPlayerService", "version": "v0001", "steamid": "steamid"}, "GetRecentlyPlayedGames": {"API Type": "IPlayerService", "version": "v0001", "steamid": "steamid"}}

var key string = os.Getenv("STEAM_API_KEY")

func GetCommunityState(steamid string) string {
	result := GetPlayerSummary(steamid)
	state := CommunityVisibilityState(result.PlayerSummaryResponse.Players[0].CommunityVisibilityState)
	return state
}

func GetMostPlayed(steamid string) map[string]float64 {
	results := FriendListPlaytime(steamid)
	mostPlayed := make(map[string]float64, 0)
	for i, res := range results {
		var playtime int
		summary := res.Summary
		recent := res.Recent

		fmt.Printf("\n\n~~~~\nFriend ID: %v\n%v:\n", summary.PlayerSummaryResponse.Players[0].PersonaName, i)

		if len(recent.RecentGamesResponse.Games) == 0 {
			fmt.Println("No games played recently")
			continue
		}
		for i, game := range recent.RecentGamesResponse.Games {
			fmt.Printf("-------\nGame ID: %v\n%s\nPast 2 weeks: %v hours\nTotal Playtime: %v hours\n", i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
			playtime += game.Playtime2Week
		}
		mostPlayed[summary.PlayerSummaryResponse.Players[0].PersonaName] = float64(playtime) / 60

	}
	return mostPlayed
}

func FriendListPlaytime(steamid string) []Result {
	list := GetFriendList(steamid)
	resultList := make([]Result, len(list.FriendListResponse.FriendList))

	resultChan := make(chan Result, len(list.FriendListResponse.FriendList))
	var wg sync.WaitGroup
	for i, friend := range list.FriendListResponse.FriendList {

		wg.Add(1)
		go func(index int, steamid string) {
			defer wg.Done()
			summary := GetPlayerSummary(steamid)
			recent := GetRecentlyPlayed(steamid)

			resultChan <- Result{index, summary, recent}
		}(i, friend.FriendSteamID)
	}

	wg.Wait()

	close(resultChan)

	for result := range resultChan {
		resultList[result.Index] = result

	}
	return resultList

}

func UnixToTime(unix int64) string {
	loc, _ := time.LoadLocation("Europe/London")
	timeReturn := time.Unix(unix, 0).In(loc)
	result := timeReturn.Format("15:04 PM 02/01/06")

	return result
}

func GetPlayerSummary(steamid string) PlayerSummary {
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
	fmt.Println(result.PlayerSummaryResponse.Players[0].GameID)

	return result

}

func GetRecentlyPlayed(steamid string) RecentGames {
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

func GetFriendList(steamid string) FriendListResponse {
	var result FriendListResponse
	url := "http://api.steampowered.com/ISteamUser/GetFriendList/v0001/?key=" + key + "&steamid=" + steamid + "&relationship=friend"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	// fmt.Println(string(bodyBytes))
	json.Unmarshal(bodyBytes, &result)
	return result
}

func PersonaStateStr(personaState int) string {
	personas := map[int]string{
		0: "Offline",
		1: "Online",
		2: "Busy",
		3: "Away",
		4: "Snooze",
		5: "Looking2Trade",
		6: "Looking2Play",
	}
	return personas[personaState]
}

func CommunityVisibilityState(state int) string {
	switch state {
	case 1:
		return "Private"
	case 3:
		return "Public"
	default:
		return ""
	}
}

func GetMostPlayed2Weeks(friends []Friend) {
	// add concurrent API calls
	var maxPlaytime int
	var topPlayerName string

	for _, friend := range friends {
		sum := GetPlayerSummary(friend.FriendSteamID)
		totalPlaytime := 0
		recent := GetRecentlyPlayed(friend.FriendSteamID)
		for _, game := range recent.RecentGamesResponse.Games {
			totalPlaytime += game.Playtime2Week
		}
		if totalPlaytime > maxPlaytime {
			topPlayerName = sum.PlayerSummaryResponse.Players[0].PersonaName
			maxPlaytime = totalPlaytime
		}

	}
	fmt.Println(topPlayerName, maxPlaytime/60)
}

func GetOwnedGames(steamid string) RecentGames {
	var result RecentGames
	url := "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=" + key + "&steamid=" + steamid + "&format=json&include_appinfo=true&include_played_free_games=false"
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

func SteamAPI(service string, steamid []string, extras ...string) {
	// apiType := apiCalls[service]["API Type"]
	// version := apiCalls[service]["version"]
	// steamidtype := apiCalls[service]["steamid"]
	// url := fmt.Sprintf("http://api.steampowered.com/%s/%s/?key=%s&%s=%s ")
	// url := "http://api.steampowered.com/"
	// var result RecentGames
	// url := "http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=" + key + "&steamid=" + steamid + "&format=json&include_appinfo=true"
	// resp, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	//
	// err = json.Unmarshal(bodyBytes, &result)
	// if err != nil {
	// 	panic(err)
	// }
	// return result
}
