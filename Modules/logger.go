package modules

import (
	"fmt"
	"time"

	"github.com/gookit/color"
)

func GetTime() string {
	currentTime := time.Now().Format("15:04:05")
	return currentTime
}

func Logger(status bool, content ...interface{}) {
	switch status {
	case true:
		color.Printf("<fg=white></><fg=blue;op=bold>%s</><fg=white> | </><fg=green>%s</><fg=green;op=bold></><fg=green></>", GetTime(), fmt.Sprintln(content...))
	case false:
		color.Printf("<fg=white></><fg=blue;op=bold>%s</><fg=white> | </><fg=red>%s</><fg=red;op=bold></><fg=red></>", GetTime(), fmt.Sprintln(content...))
	default:
		color.Printf("<fg=white></><fg=blue;op=bold>%s</><fg=white> | </><fg=cyan>%s</><fg=red;op=bold></><fg=cyan></>", GetTime(), fmt.Sprintln(content...))
	}
}
