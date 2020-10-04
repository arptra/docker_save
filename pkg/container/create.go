package container

import (
	"compress/gzip"
	"crypto/sha256"
	"docker_save/pkg/auth"
	"docker_save/pkg/parsemanifest"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	BaseLayerManifest  = "{\"id\":\"\",\"created\":\"2020-08-03T15:17:02.223842633Z\",\"container_config\":{\"Hostname\":\"\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":null,\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"Onbuild\":null,\"Labels\":null}}"
	OtherLayerManifest = "{\"id\":\"\",\"parent\":\"\",\"created\":\"2020-08-03T15:17:02.223842633Z\",\"container_config\":{\"Hostname\":\"\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":null,\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"Onbuild\":null,\"Labels\":null}}"
)

func AddId(hash string, template string) string {
	suffix := template[0:strings.Index(template, "id")+5] + hash
	res := suffix + template[strings.Index(template, "id")+5:]
	return res
}

func AddParent(hash string, template string) string {
	suffix := template[0:strings.Index(template, "parent")+9] + hash
	res := suffix + template[strings.Index(template, "parent")+9:]
	return res
}

func GetLayer(manifest parsemanifest.Manifest, prop *auth.ConProp, id int, DirName string) string {
	suffix := prop.Uri[0:strings.Index(prop.Uri, "manifest")]
	configUrl := suffix + "blobs/" + manifest.Layers.Layers[id].Digest
	resp, _, _ := prop.AuthReq.Get(configUrl).
		Set("Authorization", "Bearer "+prop.Token["token"]).
		End()
	gReader, _ := gzip.NewReader(resp.Body)
	out, _ := os.Create(DirName + "/layer.tar")
	io.Copy(out, gReader)
	out.Close()
	hash := sha256.New()
	input, _ := os.Open(DirName + "/layer.tar")
	defer input.Close()
	io.Copy(hash, input)
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	return sum
}

func CraeteManifestLayerFiles(manifestList []string, nameDir string, manifest parsemanifest.Manifest, prop *auth.ConProp) {
	last := len(manifestList)
	for index, value := range manifestList {
		file, _ := os.Create(nameDir + "/" + value + "/json")
		fileVer, _ := os.Create(nameDir + "/" + value + "/VERSION")
		io.Copy(fileVer, strings.NewReader("1.0"))
		fileVer.Close()

		if index == 0 {
			str := AddId(value, BaseLayerManifest)
			io.Copy(file, strings.NewReader(str))
		} else {
			layer := AddId(value, OtherLayerManifest)
			parentHash := manifestList[index-1]
			str := AddParent(parentHash, layer)
			if index == last-1 {
				str = str[0:strings.Index(str, "container_config")]
				str = str + parsemanifest.GetContainerConfig(parsemanifest.GetConf(prop, manifest)) + "}}," + "\"architecture\":\"amd64\"," + "\"os\":\"linux\"}"
			}
			io.Copy(file, strings.NewReader(str))
		}
		file.Close()
	}
}

func MainManifest(nameDir string, manifestList []string, manifest parsemanifest.Manifest) {
	last := len(manifestList)
	str := "[{\"Config\":\"" + manifest.Config.Config.Digest[strings.Index(manifest.Config.Config.Digest, ":")+1:] + ".json\","
	str = str + "\"RepoTags\":" + "[\"repotag:ver.test\"]," + "\"Layers\":["
	for index, value := range manifestList {
		str = str + "\""
		str = str + value + "/layer.tar"
		str = str + "\""
		if index != last-1 {
			str = str + ","
		}
	}
	str = str + "]}]"
	file, _ := os.Create(nameDir + "/manifest.json")
	io.Copy(file, strings.NewReader(str))
	file.Close()
}

func CreateRepositories(manifest parsemanifest.Manifest, nameDir string) {
	str := "{\"repotag:\"{\"ver.test\":\"" + manifest.Config.Config.Digest[strings.Index(manifest.Config.Config.Digest, ":")+1:] + "\"}}"
	file, _ := os.Create(nameDir + "/repositories")
	io.Copy(file, strings.NewReader(str))
	file.Close()
}

func Create(name string, manifest parsemanifest.Manifest, prop *auth.ConProp) {
	var manifestFile []string

	os.Mkdir(name, 0777)
	for id, _ := range manifest.Layers.Layers {
		hash := GetLayer(manifest, prop, id, name)
		manifestFile = append(manifestFile, hash)
		os.Mkdir(name+"/"+hash, 0777)
		os.Rename(name+"/layer.tar", name+"/"+hash+"/layer.tar")
	}
	CraeteManifestLayerFiles(manifestFile, name, manifest, prop)
	parsemanifest.MakeConfFile(parsemanifest.GetConf(prop, manifest), manifest, name)
	MainManifest(name, manifestFile, manifest)
	CreateRepositories(manifest, name)
}
