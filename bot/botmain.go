package bot

import (
	"adventbot/config"
	"github.com/bwmarrin/discordgo"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var Messages = make(chan config.Message)
var langtab = 0

func Main(wg *sync.WaitGroup, discord *discordgo.Session) {
	discord.AddHandler(EmoteAddHandler)
	discord.AddHandler(EmoteRemoveHandler)
	if err := discord.Open(); err != nil {
		println("Couldn't connect to discord.")
		println(err)
		os.Exit(0)
	}

	//go SetupEmotes(discord) //Comment me out after first run!
	go DayHandler(discord)

	//discord.MessageReactionAdd(config.RoleChannelID, config.LangMessageID[2], "Brainfuck:517837754957037568")

	//Got the CLI->Bot communication line down.
	stop := false

	for !stop {
		message := <-Messages
		switch message.Type {
		case config.M_NEWLANG:
			_ = discord.MessageReactionAdd(config.RoleChannelID, config.LangMessageID[langtab%len(config.LangMessageID)], message.Data)
			langtab++
		case config.M_NEWPARADIGM:
			_ = discord.MessageReactionAdd(config.RoleChannelID, config.ParadigmMessageID, message.Data)
		case config.M_RMLANG:
			for _, v := range config.LangMessageID {
				_ = discord.MessageReactionRemove(config.RoleChannelID, v, message.Data, config.BotUserID)
			}
		case config.M_RMPARADIGM:
			_ = discord.MessageReactionRemove(config.RoleChannelID, config.ParadigmMessageID, message.Data, config.BotUserID)
		case config.M_QUIT:
			stop = true
		}
	}

	if err := discord.Close(); err != nil {
		println("Troubles closing the connection.")
		println(err)
	}

	wg.Done()
}

func SetupEmotes(discord *discordgo.Session) {
	for k := range config.SetupLangToRole {
		_ = discord.MessageReactionAdd(config.RoleChannelID, config.LangMessageID[langtab%len(config.LangMessageID)], k) //It really doesn't matter if this has an error.
		langtab++
	}

	for k := range config.SetupParadigmToRole {
		_ = discord.MessageReactionAdd(config.RoleChannelID, config.ParadigmMessageID, k) //This is expected to error out on subsequent runs anyway.
	}
}

func DayHandler(discord *discordgo.Session) {
	for {
		loc, _ := time.LoadLocation("EST")
		start := time.Until(time.Date(time.Now().Year(), 12, 1, 0, 0, 0, 0, loc))
		if time.Now().Month() == 12 {
			start = time.Until(time.Date(time.Now().Year(), 12, time.Now().Day()+1, 0, 0, 0, 0, loc))
		}
		<-time.After(start)
		daily := time.Tick(24 * time.Hour)

		Greetings := []string{
			"Morning! It is December ",
			"If you're seeing this, you've either just woke up several hours ago, or you're up way too late. Either way, it is December ",
			"How are you this fine morning? It's December ",
			"What have you been wanting for Christmas? It's December ",
			"*The elves have made you a breakfast. How nice.* It's December ",
			"Rise and shine! It's the ",
			"This is plain masochism, but I won't stop ya! It's the ",
			"If I knew encouraging someone to be productive was THIS easy, I would've done it ages ago! Hope you slept well, it is December the ",
		}

		Postings := []string{
			" Today's problem has been posted. ",
			" It seems Christmas needs your help again. ",
			" Need a brain teaser to wake up? ",
			" MMM MMM! I love the smell of **NEW CODING PROBLEMS** in the morning. ",
			" Long story short, shit's fucked again. ",
			" Can you smell what Eric is cooking? ",
		}

		Harassments := []string{
			"Chop-chop!",
			"These problems won't solve themselves!",
			"Have fun!",
			"At this rate, we'll have to postpone Christmas. Please help!",
			"Send lawyers, guns, and money.",
			"I'm sure you don't get paid nearly enough for this, but my checkbook's too small to cover it!",
		}

		NaturalNumbers := []string{
			"st",
			"nd",
			"rd",
			"th",
		}

		FinalGreeting := "Merry Christmas! Hope these problems didn't pile up on you~ Get to it, I'm sure it won't take long. "
		//I love NiiNya <3
		for {
			_ = <-daily
			if time.Now().Day() < 25 {
				morning := Greetings[rand.Int()%len(Greetings)] + strconv.FormatInt(int64(time.Now().Day()), 10) + NaturalNumbers[int(math.Min(float64(time.Now().Day()-1), 3))] + Postings[rand.Int()%len(Postings)] + Harassments[rand.Int()%len(Harassments)]
				_, _ = discord.ChannelMessageSend(config.SpoilerID, morning)
				_, _ = discord.ChannelMessageSend(config.NoSpoilerID, morning)
			} else if time.Now().Day() == 25 {
				morning := FinalGreeting + Postings[rand.Int()%len(Postings)] + Harassments[rand.Int()%len(Harassments)]
				_, _ = discord.ChannelMessageSend(config.SpoilerID, morning)
				_, _ = discord.ChannelMessageSend(config.NoSpoilerID, morning)
			} else {
				break
			}
		}
	}
}

func EmoteAddHandler(discord *discordgo.Session, emote *discordgo.MessageReactionAdd) {
	langmess := false
	for _,v := range config.LangMessageID {
		if emote.MessageID == v {
			langmess = true
		}
	}

	if langmess {
		//Add a language role.
		if err := discord.GuildMemberRoleAdd(emote.GuildID, emote.UserID, config.LangToRole[emote.Emoji.ID]); err != nil {
			println("Failed to add role", config.LangToRole[emote.Emoji.ID], "to user", emote.UserID)
			println(err.Error())
		}
	} else if emote.MessageID == config.ParadigmMessageID {
		//Add a paradigm role.
		if err := discord.GuildMemberRoleAdd(emote.GuildID, emote.UserID, config.ParadigmToRole[emote.Emoji.ID]); err != nil {
			println("Failed to add role", config.ParadigmToRole[emote.Emoji.ID], "to user", emote.UserID)
			println(err.Error())
		}
	}
}

func EmoteRemoveHandler(discord *discordgo.Session, emote *discordgo.MessageReactionRemove) {
	langmess := false
	for _,v := range config.LangMessageID {
		if emote.MessageID == v {
			langmess = true
		}
	}

	if langmess {
		//Add a language role.
		if err := discord.GuildMemberRoleRemove(emote.GuildID, emote.UserID, config.LangToRole[emote.Emoji.ID]); err != nil {
			println("Failed to remove role", config.LangToRole[emote.Emoji.ID], "from user", emote.UserID)
			println(err.Error())
		}
	} else if emote.MessageID == config.ParadigmMessageID {
		//Add a paradigm role.
		if err := discord.GuildMemberRoleRemove(emote.GuildID, emote.UserID, config.ParadigmToRole[emote.Emoji.ID]); err != nil {
			println("Failed to remove role", config.ParadigmToRole[emote.Emoji.ID], "from user", emote.UserID)
			println(err.Error())
		}
	}
}
