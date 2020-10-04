package flags

import (
	"docker_save/pkg/auth"
	"flag"
)

var Configs *auth.ConfigParams

func OptionParse() {
	user := flag.String("user", "user", "string")
	pass := flag.String("pass", "pass", "string")
	url := flag.String("url", "url", "string")

	flag.Parse()
	Configs = &auth.ConfigParams{
		*user,
		*pass,
		*url,
	}
}
