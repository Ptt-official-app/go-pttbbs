package ansi

const (
	ESC_STR        = "\x1b"
	ESC_CHR        = '\x1b'
	ANSI_RESET_STR = ESC_STR + "[m"
)

var ANSI_RESET = []byte{ESC_CHR, '[', 'm'}

func ANSIColor(color string) string {
	colorStr := ESC_STR + "[" + color + "m"
	return colorStr
}

func ANSIReset() string {
	return ESC_STR + "[m"
}
