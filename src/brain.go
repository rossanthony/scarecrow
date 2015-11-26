package scarecrow

import (
	"fmt"
	rivescript "github.com/aichaos/rivescript-go"
	"os"
	"regexp"
	"time"
)

var re_nasties = regexp.MustCompile(`[^A-Za-z0-9_\-@ ]`)

// InitBrain loads the RiveScript brain.
func (self *Scarecrow) InitBrain() {
	self.Brain = rivescript.New()
	self.Brain.LoadDirectory(self.BotsConfig.Personality.Brain.Replies)
	self.Brain.SortReplies()
}

// GetReply actually gets a response for a user.
func (self *Scarecrow) GetReply(botUsername, username, message string) string {
	// Path to the user's persisted profile data.
	safeUsername := re_nasties.ReplaceAllString(username, "_")
	profile := fmt.Sprintf("./users/%s.json", safeUsername)

	// Load their user variables.
	self.LoadUservars(profile)

	// Get a reply.
	reply := self.Brain.Reply(username, message)

	// Save the user variables.
	self.SaveUservars(username, profile)

	// Log it.
	self.LogTransaction(username, message, botUsername, reply)
	return reply
}

// LogTransaction logs a full transaction between a user and the bot.
func (self *Scarecrow) LogTransaction(username, message, bot, reply string) {
	// Don't log if the bot has no username (e.g., is a console bot)
	if bot == "" {
		return
	}

	// Print to the console.
	payload := fmt.Sprintf("<%s>\n[%s] %s\n[%s] %s\n\n",
		time.Now().Format(time.RFC850), username, message, bot, reply)
	fmt.Printf(payload)

	// And make logs on disk.
	safeUsername := re_nasties.ReplaceAllString(username, "_")
	safeBot := re_nasties.ReplaceAllString(bot, "_")
	logDir := fmt.Sprintf("./logs/%s", safeBot)
	logFile := fmt.Sprintf("%s/%s.log", logDir, safeUsername)
	MakeDirectory(logDir)

	fh, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		self.Error("Couldn't write log file: %s", err)
		return
	}
	defer fh.Close()

	if _, err = fh.WriteString(payload); err != nil {
		self.Error("Couldn't write log file: %s", err)
	}
}
