# Steam Watch
Steam Watch is a Go application that uses the Steam Web API to analyse a user’s friend list and playtime information.
Functionality of this program fully depends on the user having a public profile.

Has a TUI interface created using [bubbletea](https://github.com/charmbracelet/bubbletea) and CLI commands.

## Features
* Gather information of a users entire friend list, including SteamID
* Use SteamID to get a player summary, which includes:
  *  Profile Name
  *  Profile State
  *  Real Name
  *  Location Country Code
  *  Time Created
  *  and more !
* Identifies most played games:
  * Within the last 2 weeks
  * Overall   
* Generates insights from a friends list by combining all the data

## Getting Started
Make sure you find your SteamID before you use this application.
There are multiple methods to find your SteamID but I recommended this website: [SteamID](https://steamidcheck.com/profile) and follow the onscreen prompts to get your SteamID.

Having Go installed is the only prerequisite for running this project.

### Installation

( need to finish then it can be easily installed )

## Usage

This application has TUI and CLI built-in.

### CLI
list cli commands

### TUI

`(app name) --tui` will start the application with a TUI, the controls will be onscreen once inside the TUI. 


## Credits
[API Docs](https://steamcommunity.com/dev)
