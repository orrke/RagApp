package config

import (
	"sync"
)

// Global configuration variables
var (
	Path   string
	Lock   sync.RWMutex
	Config ServerConfig
)

// Languages is a map of all the languages supported by bleve
// The written version is there for a clearer version for the model
var Languages = map[string]string{
	"ar":  "عربي",
	"bg":  "български",
	"ca":  "català",
	"cjk": "日本語",
	"ckb": "کوردی",
	"cs":  "čeština",
	"da":  "dansk",
	"de":  "Deutsch",
	"el":  "ελληνικά",
	"en":  "English",
	"es":  "español",
	"eu":  "euskara",
	"fa":  "فارسی",
	"fi":  "suomi",
	"fr":  "Français",
	"ga":  "Gaeilge",
	"gl":  "galego",
	"hi":  "हिन्दी",
	"hr":  "hrvatski",
	"hu":  "magyar",
	"hy":  "հայերեն",
	"id":  "Bahasa Indonesia",
	"in":  "हिन्दी",
	"it":  "italiano",
	"nl":  "Nederlands",
	"no":  "norsk",
	"pl":  "polski",
	"pt":  "português",
	"ro":  "română",
	"ru":  "русский",
	"sv":  "svenska",
	"tr":  "Türkçe",
}
