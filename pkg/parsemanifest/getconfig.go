package parsemanifest

import (
	"docker_save/pkg/auth"
	"fmt"
	"os"
	"strings"
)

type ContainerConfig struct {
	ContainerConf string `json:"container_config"`
}

func GetConf(prop *auth.ConProp, manifest Manifest) string {
	suffix := prop.Uri[0:strings.Index(prop.Uri, "manifest")]
	configUrl := suffix + "blobs/" + manifest.Config.Config.Digest
	_, body, _ := prop.AuthReq.Get(configUrl).
		Set("Authorization", "Bearer "+prop.Token["token"]).
		End()
	return body
}

func GetContainerConfig(data string) string {
	str := data[strings.Index(data, "container_config"):strings.Index(data, "}},\"created")]
	return str
}

func MakeConfFile(data string, manifest Manifest, nameDir string) {
	name := manifest.Config.Config.Digest[strings.Index(manifest.Config.Config.Digest, ":")+1:]
	file, err := os.Create(nameDir + "/" + name + ".json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(data)
}
