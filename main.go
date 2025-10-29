package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"steam/code/api"
	"strings"
)

type kv struct {
	Key   string
	Value float64
}

const space = "~~~~~~~~~~~~~~~"

// Main function
func main() {
	steamid := flag.String("id", "", "Enter user's SteamID")
	friendListFlag := flag.NewFlagSet("FL", flag.ExitOnError)
	summary := friendListFlag.Bool("s", false, "Shows a summary of a friendlist")
	filter := friendListFlag.String("f", "", "Filter when showing friend list (online, offline) ")

	functionCmd := flag.NewFlagSet("function", flag.ExitOnError)
	playerSummaryCmd := functionCmd.Bool("PS", false, "Get player summary")
	recentlyPlayedCmd := functionCmd.Bool("RS", false, "Get games recently played")
	ownedGamesCmd := functionCmd.Bool("OG", false, "Get overall played games")
	mostPlayedListCmd := functionCmd.Bool("MP", false, "Get most played of friend list")

	if len(os.Args) < 2 {
		fmt.Println("You must enter one flag.")
		os.Exit(1)
	}

	flag.Parse()

	if *steamid != "" {
		fmt.Fprintf(os.Stdout, "Loaded Steam ID: %s", *steamid)
	} else {
		fileText, err := os.ReadFile("steamid.txt")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		*steamid = string(fileText)
		fmt.Fprintf(os.Stdout, "Loaded Steam ID: %s", *steamid)
	}
	if *steamid == "" {
		fmt.Fprintln(os.Stderr, "No SteamID found from flag or local save.")
		os.Exit(1)
	}
	if os.Args[1] == "FL" {
		var output string
		friendListFlag.Parse(os.Args[2:])
		friendList := api.FriendListData(*steamid)
		var newFriendList api.DetailedFriendList
		if *summary {
			for i, friend := range friendList.DetailedFriendList {
				details := friend.FriendSummary.PlayerSummaryResponse.Players[0]
				output += fmt.Sprintf("%s\nNumber: %v\nName: %s\nFriend since: %s\nCurrently: %s\nRelationship: %s\nLast Logoff: %s\n", space, i, details.PersonaName, api.UnixToTime(friend.FriendDetails.FriendSince), api.PersonaStateStr(details.PersonaState), friend.FriendDetails.Relationship, api.UnixToTime(details.LastLogoff))
				if details.CommunityVisibilityState == 3 {
					output += fmt.Sprintf("\nTime Created: %s\nCurrently playing: %s\nLocation: %s\n%s", api.UnixToTime(details.TimeCreated), details.GameExtraInfo, details.LocCountryCode, space)
				}
			}

		} else if *filter != "" {
			switch strings.ToLower(*filter) {
			case "online":
				output += "Showing whos online:\n"
				newFriendList = SpecificFriendStatus(friendList, isOnline)
			case "offline":
				output += "Showing whos offline:\n"
				newFriendList = SpecificFriendStatus(friendList, isOffline)
			case "playing":
				output += "Showing whos in-game:\n"
				newFriendList = SpecificFriendStatus(friendList, isPlaying)
			}
			for _, friend := range newFriendList.DetailedFriendList {
				output += fmt.Sprintf("%s\nName: %s\nCurrently: %s\nRelationship: %s\n", space, friend.FriendSummary.PlayerSummaryResponse.Players[0].PersonaName, api.PersonaStateStr(friend.FriendSummary.PlayerSummaryResponse.Players[0].PersonaState), friend.FriendDetails.Relationship)
			}
		}
		fmt.Fprintln(os.Stdout, output)
	} else if os.Args[1] == "function" {
		functionCmd.Parse(os.Args[2:])
		if *playerSummaryCmd {
			playerSummary := api.GetPlayerSummary(*steamid)
			player := playerSummary.PlayerSummaryResponse.Players[0]
			state := api.CommunityVisibilityState(player.CommunityVisibilityState)

			output := fmt.Sprintf("%s:\n%s\nLast Logoff: %s", player.PersonaName, api.PersonaStateStr(player.PersonaState), api.UnixToTime(player.LastLogoff))
			if state == "Public" {
				output += fmt.Sprintf("\nTime Created: %s\nCurrently playing: %s\nLocation: %s", api.UnixToTime(player.TimeCreated), player.GameExtraInfo, player.LocCountryCode)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *recentlyPlayedCmd {
			recentlyPlayed := api.GetRecentlyPlayed(*steamid)
			fmt.Fprintf(os.Stdout, "Total Games:%v\n", recentlyPlayed.RecentGamesResponse.TotalCount)
			output := ""

			for i, game := range recentlyPlayed.RecentGamesResponse.Games {
				output += fmt.Sprintf("%s\n%v\n%s:\nPlaytime 2 Weeks: %v hours\nPlaytime Overall: %v hours\n", space, i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *ownedGamesCmd {
			ownedGames := api.GetOwnedGames(*steamid)
			fmt.Fprintf(os.Stdout, "Total Games:%v\n", ownedGames.RecentGamesResponse.GamesCount)
			output := ""
			for i, game := range ownedGames.RecentGamesResponse.Games {
				output += fmt.Sprintf("%s\n%v\n%s:\nPlaytime 2 Weeks: %v hours\nPlaytime Overall: %v hours\n", space, i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *mostPlayedListCmd {
			results := api.FriendListData(*steamid)
			mostPlayed := make(map[string]float64, 0)
			for i, res := range results.DetailedFriendList {
				var playtime int
				summary := res.FriendSummary
				recent := res.RecentGames

				fmt.Printf("\n\n~~~~\nFriend ID: %v\n%v:\n", summary.PlayerSummaryResponse.Players[0].PersonaName, i)
				// fmt.Println(recent.RecentGamesResponse.Games)
				// fmt.Println(len(recent.RecentGamesResponse.Games))

				if len(recent.Games) == 0 {
					fmt.Println("No games played recently")
					continue
				}
				for i, game := range recent.Games {
					fmt.Printf("-------\nGame ID: %v\n%s\nPast 2 weeks: %v hours\nTotal Playtime: %v hours\n", i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
					playtime += game.Playtime2Week
				}
				mostPlayed[summary.PlayerSummaryResponse.Players[0].PersonaName] = float64(playtime) / 60
				// fmt.Println(summary.PlayerSummaryResponse.Players[0].PersonaName, playtime)

			}
			ss := SortMap(mostPlayed)
			var output string
			for i, kv := range ss {
				output += fmt.Sprintf("%v - %s, %f\n", i, kv.Key, kv.Value)
			}
			fmt.Fprintln(os.Stdout, output)
		}

	} else {
		file, err := os.Create("steamid.txt")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer file.Close()
		_, err = file.WriteString(*steamid)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// Function that takes a string float64 map and returns a sorted map.
func SortMap(unordered map[string]float64) []kv {
	var ss []kv
	for k, v := range unordered {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	return ss
}

// Function which uses a populated friend list to only show a specific status.
// For example, if I wanted to only see who was online, I can input online into
// the status parameter and the function will return the details of the friends
// that are online.
func SpecificFriendStatus(friendList api.DetailedFriendList, condition func(api.DetailedFriend) bool) api.DetailedFriendList {

	var newList api.DetailedFriendList
	for _, friend := range friendList.DetailedFriendList {
		if condition(friend) == true {
			newList.DetailedFriendList = append(newList.DetailedFriendList, friend)
		}
	}
	return newList
}

func isOnline(friend api.DetailedFriend) bool {
	if friend.FriendSummary.PlayerSummaryResponse.Players[0].PersonaState != 0 {
		return true
	} else {
		return false
	}
}
func isOffline(friend api.DetailedFriend) bool {
	if friend.FriendSummary.PlayerSummaryResponse.Players[0].PersonaState == 0 {
		return true
	} else {
		return false
	}
}

func isPlaying(friend api.DetailedFriend) bool {
	if friend.FriendSummary.PlayerSummaryResponse.Players[0].GameExtraInfo != "" {
		return true
	} else {
		return false
	}
}
