package core

type Runtime struct {
	Name      string
	Extension string
}

var LanguageMatrix = map[string]*Runtime{
	"gnuc++11": {
		Name:      "GNU C++ 11",
		Extension: "cpp",
	},
	"gnuc++14": {
		Name:      "GNU C++ 14",
		Extension: "cpp",
	},
	"gnuc++17": {
		Name:      "GNU C++ 17",
		Extension: "cpp",
	},
	"gnuc++20": {
		Name:      "GNU C++ 20",
		Extension: "cpp",
	},
	"python3": {
		Name:      "Python 3",
		Extension: "py",
	},
	"go": {
		Name:      "Go",
		Extension: "go",
	},
}
