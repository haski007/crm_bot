package botlogs

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Logger struct to log updates of telegram bot
type Logger struct {
	path string
}

// NewLogger creates logger object,
// dirPath - logs directory.
// If dirPath is not directory - creates default logs directory in root path
func NewLogger(dirPath string) *Logger {
	if isDir(dirPath) {
		if []rune(dirPath)[utf8.RuneCountInString(dirPath)-1] != rune('/') {
			dirPath += "/"
		}

		return &Logger{path: dirPath}
	}
	os.Mkdir("logs", os.ModePerm)
	return &Logger{path: "logs/"}
}

// MessageLog build log to file with name logger.dirPath + [date] + .log
func (l *Logger) MessageLog(from *tgbotapi.User, text string) error {

	// ---> Open or create file
	date := time.Now().Format("02-01-2006")

	f, err := os.OpenFile(l.path + date + ".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	log := fmt.Sprintf(
`--------------------
...Message...
From: %v [%v]
User id: %v
Time: %v
Text: "%v"
Language: %v
`,
		from.FirstName + from.LastName, "@" + from.UserName,
		from.ID,
		time.Now().Format("15:04:05"),
		text,
		from.LanguageCode)

	
	_, err = f.WriteString(log)
	if err != nil {
		return err
	}

	return nil
}
