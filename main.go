package main

import (
	"adventbot/bot"
	"adventbot/cli"
	"adventbot/config"
	"adventbot/utils"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"sync"
)

func main() {
	if len(os.Args) == 2 {
		if _, err := os.Stat("roles.json"); err == nil {
			if bytes, err := ioutil.ReadFile("roles.json"); err == nil {
				if err := json.Unmarshal(bytes, &config.LangToRole); err != nil {
					utils.SaveFile("roles.json", config.LangToRole)
				}
			} else {
				utils.SaveFile("roles.json", config.LangToRole)
			}
		} else {
			utils.SaveFile("roles.json", config.LangToRole)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		//TODO: Discord bot.
		if discord, err := discordgo.New("Bot " + os.Args[1]); err == nil {
			go bot.Main(&wg, discord)
		} else {
			println("Couldn't initialize discord session.")
			println(err)
			os.Exit(0)
		}

		//TODO: CLI.
		go cli.Main(&wg)

		wg.Wait()
	} else {
		println("adventbot [API Key]")
		println("This bot is made specifically for the adventofcode discord server.")
		println("You'll need to modify the source yourself if you want it to work in your server.")
		println("Everything actually *important* is under the config/config.go file.")
	}
}
