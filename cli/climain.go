package cli

import (
	"adventbot/bot"
	"adventbot/config"
	"adventbot/utils"
	"bufio"
	"os"
	"strings"
	"sync"
)

func Main(wg *sync.WaitGroup) {
	var exit = false
	reader := bufio.NewReader(os.Stdin)

	for !exit {
		read := strings.Split(ReadLine(reader), " ")

		switch read[0] {
		case "lang":
			if len(read) > 3 {
				switch read[1] {
				case "add":
					if len(read) > 4 {
						config.LangToRole[read[2]] = read[3]
						utils.SaveFile("roles.json", config.LangToRole)
						bot.Messages <- config.Message{
							Type: config.M_NEWLANG,
							Data: read[2],
						}
					} else {
						println("Usage: lang [add/remove] [emote ID] (role ID, only needed on add)")
					}
				case "remove":
					delete(config.LangToRole, read[2])
					utils.SaveFile("roles.json", config.LangToRole)
					bot.Messages <- config.Message{
						Type: config.M_RMLANG,
						Data: read[2],
					}
				default:
					println("Usage: lang [add/remove] [emote ID] (role ID, only needed on add)")
				}
			} else {
				println("Usage: lang [add/remove] [emote ID] (role ID, only needed on add)")
			}
		case "paradigm":
			if len(read) > 3 {
				switch read[1] {
				case "add":
					if len(read) > 4 {
						config.ParadigmToRole[read[2]] = read[3]
						utils.SaveFile("paradigm.json", config.ParadigmToRole)
						bot.Messages <- config.Message{
							Type: config.M_NEWPARADIGM,
							Data: read[2],
						}
					} else {
						println("Usage: paradigm [add/remove] [emote ID] (role ID, only needed on add)")
					}
				case "remove":
					delete(config.ParadigmToRole, read[2])
					utils.SaveFile("paradigm.json", config.ParadigmToRole)
					bot.Messages <- config.Message{
						Type: config.M_RMPARADIGM,
						Data: read[2],
					}
				default:
					println("Usage: paradigm [add/remove] [emote ID] (role ID, only needed on add)")
				}
			} else {
				println("Usage: paradigm [add/remove] [emote ID] (role ID, only needed on add)")
			}
		case "exit":
			exit = true
			bot.Messages <- config.Message{
				Type: config.M_QUIT,
				Data: "",
			}
			println("Shutting the bot down.")
		default:
			println("lang [add/remove] [emote ID] (role ID) -- Adds or removes roles on the fly.")
			println("paradigm [add/remove] [emote ID] (role ID) -- Adds or removes roles on the fly.")
			println("exit -- Quits the bot.")
		}
	}

	wg.Done()
}

func ReadLine(reader *bufio.Reader) string {
	if data, err := reader.ReadString('\n'); err == nil {
		return strings.Split(data, "\n")[0]
	}

	println("Can't read from stdio. Exiting!")

	return "exit"
}
