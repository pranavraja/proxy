package proxy

import (
	"fmt"
	"net/url"
)

const Reset = "\033[0m"
const ColorFormat = "\033[%dm"

type Color int

var (
	Bold      Color = 1
	Dark      Color = 2
	Reverse   Color = 7
	Underline Color = 4

	Blue    Color = 34
	Cyan    Color = 36
	Green   Color = 32
	Grey    Color = 30
	Magenta Color = 35
	Red     Color = 31
	White   Color = 37
	Yellow  Color = 33

	BgBlue    Color = 44
	BgCyan    Color = 46
	BgGreen   Color = 42
	BgGrey    Color = 40
	BgMagenta Color = 45
	BgRed     Color = 41
	BgWhite   Color = 47
	BgYellow  Color = 43
)

func Colored(str string, colors ...Color) string {
	coloredStr := ""
	for _, color := range colors {
		coloredStr += fmt.Sprintf(ColorFormat, color)
	}
	return coloredStr + str + Reset
}

func ColoredMethod(method string) string {
	color := Blue
	switch method {
	case "DELETE":
		color = Magenta
	case "PUT", "POST":
		color = Cyan
	case "GET", "HEAD":
		color = Yellow
	}
	return Colored(method, Bold, color)
}

func ColoredStatusLine(statusCode int, statusText string) string {
	color := Magenta
	switch statusCode {
	case 301, 302:
		color = Cyan
	case 200, 304:
		color = Green
	case 404, 500:
		color = Red
	}
	return Colored(fmt.Sprintf("%d %s", statusCode, statusText), Bold, color)
}

func ColoredURL(u *url.URL) string {
	return fmt.Sprintf("%s://%s%s", Colored(u.Scheme, White, Bold), Colored(u.Host, White, Bold), u.RequestURI())
}
