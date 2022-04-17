package berrt

var iconMap = map[string]string{
	"ok":        "\xF0\x9F\x91\x8D",
	"bad":       "\xF0\x9F\x91\x8E",
	"bug":       "\xF0\x9F\x90\x9E",
	"idea":      "\xF0\x9F\x8E\xA1",
	"explode":   "\xF0\x9F\x92\xA5",
	"trace":     "\xF0\x9F\x94\x8D",
	"info":      "\xF0\x9F\x94\xB0",
	"internal":  "\xE2\x99\xBB",
	"warning":   "\xE2\x9A\xA0",
	"lightning": "\xE2\x9A\xA1",
}

func NewIcon(name string) string {
	return iconMap[name]
}
