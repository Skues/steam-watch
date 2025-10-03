package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"steam/code/api"
)

type kv struct {
	Key   string
	Value float64
}

func main() {
	steamid := flag.String("id", "", "Enter user's SteamID")
	functionCmd := flag.NewFlagSet("function", flag.ExitOnError)
	friendListCmd := functionCmd.Bool("FL", false, "Friend list function")
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
		fmt.Println("SteamID: ", *steamid)
	} else {
		fileText, err := os.ReadFile("steamid.txt")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		*steamid = string(fileText)
		fmt.Println(*steamid)

	}
	if *steamid == "" {
		fmt.Fprintln(os.Stderr, "No SteamID found from flag or local save.")
		os.Exit(1)
	}
	if os.Args[1] == "function" {
		functionCmd.Parse(os.Args[2:])
		space := "~~~~~~~~~~~~~~~"
		if *friendListCmd {
			friendList := api.GetFriendList(*steamid)
			for i, friend := range friendList.FriendListResponse.FriendList {
				summary := api.GetPlayerSummary(friend.FriendSteamID)
				fmt.Fprintf(os.Stdout, "%s\nNumber: %v\nName: %s\nFriend since: %s\nCurrently: %s\nRelationship: %s\n%s\n", space, i, summary.PlayerSummaryResponse.Players[0].PersonaName, api.UnixToTime(friend.FriendSince), api.PersonaStateStr(summary.PlayerSummaryResponse.Players[0].PersonaState), friend.Relationship, space)
			}
		}
		if *playerSummaryCmd {
			playerSummary := api.GetPlayerSummary(*steamid)
			state := api.CommunityVisibilityState(playerSummary.PlayerSummaryResponse.Players[0].CommunityVisibilityState)

			player := playerSummary.PlayerSummaryResponse.Players[0]
			output := fmt.Sprintf("%s:\n%s\nLast Online: %s", player.PersonaName, api.PersonaStateStr(player.PersonaState), api.UnixToTime(player.LastLogoff))
			if state == "Public" {
				output += fmt.Sprintf("\nTime Created: %s\nCurrently playing: %s\nLocation: %s", api.UnixToTime(player.TimeCreated), player.GameExtraInfo, player.LocCountryCode)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *recentlyPlayedCmd {
			recentlyPlayed := api.GetRecentlyPlayed(*steamid)
			fmt.Fprintf(os.Stdout, "Total Games:%v", recentlyPlayed.RecentGamesResponse.TotalCount)
			output := ""
			for i, game := range recentlyPlayed.RecentGamesResponse.Games {
				output += fmt.Sprintf("%v\n%s:\nPlaytime 2 Weeks: %v\nPlaytime Overall:%v\n", i, game.Name, game.Playtime2Week, game.PlaytimeForever)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *ownedGamesCmd {
			ownedGames := api.GetOwnedGames(*steamid)
			fmt.Fprintf(os.Stdout, "Total Games:%v", ownedGames.RecentGamesResponse.TotalCount)
			output := ""
			for i, game := range ownedGames.RecentGamesResponse.Games {
				output += fmt.Sprintf("%v\n%s:\nPlaytime 2 Weeks: %v\nPlaytime Overall:%v\n", i, game.Name, game.Playtime2Week, game.PlaytimeForever)
			}
			fmt.Fprintln(os.Stdout, output)
		}
		if *mostPlayedListCmd {
			results := api.FriendListPlaytime(*steamid)
			mostPlayed := make(map[string]float64, 0)
			for i, res := range results {
				var playtime int
				summary := res.Summary
				recent := res.Recent

				fmt.Printf("\n\n~~~~\nFriend ID: %v\n%v:\n", summary.PlayerSummaryResponse.Players[0].PersonaName, i)
				// fmt.Println(recent.RecentGamesResponse.Games)
				// fmt.Println(len(recent.RecentGamesResponse.Games))

				if len(recent.RecentGamesResponse.Games) == 0 {
					fmt.Println("No games played recently")
					continue
				}
				for i, game := range recent.RecentGamesResponse.Games {
					fmt.Printf("-------\nGame ID: %v\n%s\nPast 2 weeks: %v hours\nTotal Playtime: %v hours\n", i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
					playtime += game.Playtime2Week
				}
				mostPlayed[summary.PlayerSummaryResponse.Players[0].PersonaName] = float64(playtime) / 60
				// fmt.Println(summary.PlayerSummaryResponse.Players[0].PersonaName, playtime)

			}
			ss := SortMap(mostPlayed)
			var output string
			for _, kv := range ss {
				output += fmt.Sprintf("%s, %f\n", kv.Key, kv.Value)
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
