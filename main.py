import requests
import json

url = "https://overfast-api.tekrop.fr/"

# allHeroes = requests.get(url + "heroes")
# print(allHeroes.content)
specificURL = "https://overfast-api.tekrop.fr/players/CFR-21760/stats/career?gamemode=competitive&platform=pc"


def getPlayer(player: str):
    data = requests.get(url + f"players/{player}/stats/summary")
    code = data.status_code
    print(code)
    return data.json()


def playerCareer(player: str, gamemode: str, platform: str):
    data = requests.get(
        url + f"players/{player}/stats/career?gamemode={gamemode}&platform={platform}"
    )
    return data.json()


with open("data.json", "w") as f:
    json.dump(getPlayer("CFR-21760"), f)
    print("written")

# print(getPlayer("CFR-21760"))
# print(allHeroes.json())
