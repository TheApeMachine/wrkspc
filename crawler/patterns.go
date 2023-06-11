package crawler

import (
	"regexp"

	"github.com/wrk-grp/errnie"
)

/*
Define a slice of string containing the country codes we want to ignore.
*/
var ignoredCountryCodes = []string{
	"ad", "ae", "af", "ag", "ai", "al", "am", "ao", "aq", "ar", "as", "at",
	"au", "aw", "ax", "az", "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi",
	"bj", "bl", "bm", "bn", "bo", "bq", "br", "bs", "bt", "bv", "bw", "by",
	"bz", "ca", "cc", "cd", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn",
	"co", "cr", "cu", "cv", "cw", "cx", "cy", "cz", "de", "dj", "dk", "dm",
	"do", "dz", "ec", "ee", "eg", "eh", "er", "es", "et", "eu", "fi", "fj",
	"fk", "fm", "fo", "fr", "ga", "gb", "gd", "ge", "gf", "gg", "gh", "gi",
	"gl", "gm", "gn", "gp", "gq", "gr", "gs", "gt", "gu", "gw", "gy", "hk",
	"hm", "hn", "hr", "ht", "hu", "id", "ie", "il", "im", "in", "io", "iq",
	"ir", "is", "it", "je", "jm", "jo", "jp", "ke", "kg", "kh", "ki", "km",
	"kn", "kp", "kr", "kw", "ky", "kz", "la", "lb", "lc", "li", "lk", "lr",
	"ls", "lt", "lu", "lv", "ly", "ma", "mc", "md", "me", "mf", "mg", "mh",
	"mk", "ml", "mm", "mn", "mo", "mp", "mq", "mr", "ms", "mt", "mu", "mv",
	"mw", "mx", "my", "mz", "na", "nc", "ne", "nf", "ng", "ni", "nl", "no",
	"np", "nr", "nu", "nz", "om", "pa", "pe", "pf", "pg", "ph", "pk", "pl",
	"pm", "pn", "pr", "ps", "pt", "pw", "py", "qa", "re", "ro", "rs", "ru",
	"rw", "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sj", "sk", "sl",
	"sm", "sn", "so", "sr", "ss", "st", "sv", "sx", "sy", "sz", "tc", "td",
	"tf", "tg", "th", "tj", "tk", "tl", "tm", "tn", "to", "tr", "tt", "tv",
	"tw", "tz", "ua", "ug", "uk", "us", "uy", "uz", "va", "vc", "ve", "vg",
	"vi", "vn", "vu", "wf", "ws", "ye", "yt", "za", "zm", "zw",
}

/*
normalizeLink normalizes a link by removing the trailing slash and adding the
base URL if it is not present.
*/
func (s *Spider) normalizeLink(link string) (string, bool) {
	errnie.Trace()

	// Check if the link is empty and return if so.
	if link == "" || link == "/" {
		return link, false
	}

	// Check if the link is directing to a alternate language and return if so.
	// First we split the link by the slash character, as well as the period.
	// Then we can check if any of the parts are equal to the language codes.
	// If so, return the link.
	parts := regexp.MustCompile(`[/.]`).Split(link, -1)
	for _, part := range parts {
		if s.contains(ignoredCountryCodes, part) {
			return link, false
		}
	}

	// Check if the link is a mailto link and return if so.
	if len(link) < 7 || link[:7] == "mailto:" {
		return link, false
	}

	// Check if the link is a tel link and return if so.
	if len(link) < 4 || link[:4] == "tel:" {
		return link, false
	}

	// Check if the link is a javascript link and return if so.
	if len(link) < 11 || link[:11] == "javascript:" {
		return link, false
	}

	// Check if the link is a hash link and return if so.
	if len(link) < 2 || link[:2] == "//" {
		return link, false
	}

	// Check if the last character is a slash and remove it if so.
	if link[len(link)-1] == '/' {
		link = link[:len(link)-1]
	}

	// Check if the link contains www. and remove it if so.
	if link[:4] == "www." {
		link = link[4:]
	}

	// Check for http or https in the link and add the base URL if not present.
	if link[:4] != "http" {
		link = s.url + link
	}

	return link, true
}

func (s *Spider) contains(arr []string, str string) bool {
	errnie.Trace()

	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}
