package auth

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/parnurzeal/gorequest"
)

type ConfigParams struct {
	User     string
	Password string
	Url      string
}

type ConProp struct {
	AuthReq *gorequest.SuperAgent
	Uri     string
	Token   map[string]string
}

// This function parses the Www-Authenticate header provided in the challenge
// It has the following format
// Bearer realm="https://gitlab.com/jwt/auth",service="container_registry",scope="repository:andrew18/container-test:pull"
func parseBearer(bearer []string) map[string]string {
	out := make(map[string]string)
	for _, b := range bearer {
		for _, s := range strings.Split(b, " ") {
			if s == "Bearer" {
				continue
			}
			for _, params := range strings.Split(s, ",") {
				fields := strings.Split(params, "=")
				key := fields[0]
				val := strings.Replace(fields[1], "\"", "", -1)
				out[key] = val
			}
		}
	}
	return out
}

func GetToken(conf *ConfigParams) map[string]string {

	// Based on
	// http://www.cakesolutions.net/teamblogs/docker-registry-api-calls-as-an-authenticated-user

	request := gorequest.New()

	//url := "https://index.docker.io/v2/odewahn/myalpine/tags/list"
	//url := "https://index.docker.io/v2/odewahn/myalpine/manifests/latest"

	// First step is to get the endpoint where we'll be authenticating
	resp, _, _ := request.Get(conf.Url).End()

	// This has the various things we'll need to parse and use in the request
	params := parseBearer(resp.Header["Www-Authenticate"])
	paramsJSON, _ := json.Marshal(&params)
	log.Println(string(paramsJSON))

	// Get the token
	challenge := gorequest.New()
	resp, body, _ := challenge.Get(params["realm"]).
		SetBasicAuth(conf.User, conf.Password).
		Query(string(paramsJSON)).
		End()

	token := make(map[string]string)
	json.Unmarshal([]byte(body), &token)

	// Now reissue the challenge with the toekn in the Header
	// curl -IL https://index.docker.io/v2/odewahn/image/tags/list
	return token
}
