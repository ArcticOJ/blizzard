package rejson

type Command string

const (
	SET       Command = "JSON.SET"
	GET       Command = "JSON.GET"
	DEL       Command = "JSON.DEL"
	MGET      Command = "JSON.MGET"
	TYPE      Command = "JSON.TYPE"
	NUMINCRBY Command = "JSON.NUMINCRBY"
	NUMMULTBY Command = "JSON.NUMMULTBY"
	STRAPPEND Command = "JSON.STRAPPEND"
	STRLEN    Command = "JSON.STRLEN"
	ARRAPPEND Command = "JSON.ARRAPPEND"
	ARRLEN    Command = "JSON.ARRLEN"
	ARRPOP    Command = "JSON.ARRPOP"
	ARRINDEX  Command = "JSON.ARRINDEX"
	ARRTRIM   Command = "JSON.ARRTRIM"
	ARRINSERT Command = "JSON.ARRINSERT"
	OBJKEYS   Command = "JSON.OBJKEYS"
	OBJLEN    Command = "JSON.OBJLEN"
	DEBUG     Command = "JSON.DEBUG"
	FORGET    Command = "JSON.FORGET"
	RESP      Command = "JSON.RESP"
)
