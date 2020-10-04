package main

import (
	"docker_save/pkg/auth"
	"docker_save/pkg/container"
	"docker_save/pkg/flags"
	"docker_save/pkg/parsemanifest"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

func main() {
	//configure
	flags.OptionParse()
	authconf := flags.Configs
	//fmt.Println(authconf.User, authconf.Password, authconf.Url)

	//init
	token := auth.GetToken(authconf)
	authenticatedRequest := gorequest.New()
	_, body, _ := authenticatedRequest.Get(authconf.Url).
		Set("Authorization", "Bearer "+token["token"]).
		End()

	conProp := &auth.ConProp{
		authenticatedRequest,
		authconf.Url,
		token,
	}

	//driver program
	manifest := parsemanifest.Parse(body)
	fmt.Println(manifest)
	container.Create("Container", manifest, conProp)

}
