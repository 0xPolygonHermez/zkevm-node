package dependencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	defaultImageAPIServer = "https://hub.docker.com"
	defaultLoginPattern   = "/v2/users/login"
)

// ImagesConfig is the configuration for the images updater.
type ImagesConfig struct {
	Names          []string
	TargetFilePath string
}

type imageUpdater struct {
	fs afero.Fs

	targetFilePath string
	imageAPIServer string

	dockerUsername string
	dockerPassword string

	images []string
}

func newImageUpdater(images []string, targetFilePath string) *imageUpdater {
	return &imageUpdater{
		fs: afero.NewOsFs(),

		targetFilePath: targetFilePath,
		imageAPIServer: defaultImageAPIServer,

		dockerUsername: os.Getenv("DOCKERHUB_USERNAME"),
		dockerPassword: os.Getenv("DOCKERHUB_PASSWORD"),

		images: images,
	}
}

func (iu *imageUpdater) update() error {
	for _, image := range iu.images {
		if err := iu.updateImage(image); err != nil {
			return err
		}
	}
	return nil
}

func (iu *imageUpdater) updateImage(imageName string) error {
	currentDigest, err := iu.readCurrentDigest(imageName)
	if err != nil {
		return err
	}

	remoteDigest, err := iu.readRemoteDigest(imageName)
	if err != nil {
		return err
	}

	if currentDigest != remoteDigest {
		if err := iu.updateDigest(imageName, currentDigest, remoteDigest); err != nil {
			return err
		}
	}
	return nil
}

func (iu *imageUpdater) readCurrentDigest(imageName string) (string, error) {
	data, err := afero.ReadFile(iu.fs, getTargetPath(iu.targetFilePath))
	if err != nil {
		return "", err
	}

	content := struct {
		Version  string
		Services map[string]struct {
			Image string
		}
	}{}
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return "", err
	}
	for _, c := range content.Services {
		if c.Image == "" {
			continue
		}
		items := strings.Split(c.Image, "@")
		const requiredItems = 2
		if len(items) < requiredItems {
			continue
		}
		if strings.HasPrefix(items[0], imageName) {
			log.Infof("Current digest of %q is %q", imageName, items[1])
			return items[1], nil
		}
	}
	return "", fmt.Errorf("image %q not found in %q", imageName, iu.targetFilePath)
}

func (iu *imageUpdater) readRemoteDigest(imageName string) (string, error) {
	token, err := iu.dockerLogin()
	if err != nil {
		return "", err
	}
	return iu.readLatestTag(imageName, token)
}

func (iu *imageUpdater) updateDigest(imageName, currentDigest, remoteDigest string) error {
	targetFilePath := getTargetPath(iu.targetFilePath)
	log.Infof("Updating %q...", targetFilePath)
	oldContent, err := afero.ReadFile(iu.fs, targetFilePath)
	if err != nil {
		return err
	}
	oldImageField := fmt.Sprintf("%s@%s", imageName, currentDigest)
	newImageField := fmt.Sprintf("%s@%s", imageName, remoteDigest)

	newContent := strings.ReplaceAll(string(oldContent), oldImageField, newImageField)

	return afero.WriteFile(iu.fs, targetFilePath, []byte(newContent), 0664)
}

func (iu *imageUpdater) dockerLogin() (string, error) {
	target := fmt.Sprintf("%s%s", iu.imageAPIServer, defaultLoginPattern)
	jsonStr := fmt.Sprintf(`{"username":"%s","password":"%s"}`, iu.dockerUsername, iu.dockerPassword)
	req, err := http.NewRequest(
		"POST", target,
		bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body == nil {
		return "", fmt.Errorf("Empty body returned")
	}

	defer func() {
		err = res.Body.Close()
	}()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Debugf("returned content: %s", content)

	decoder := json.NewDecoder(res.Body)
	r := struct {
		Token string
	}{}
	err = decoder.Decode(&r)
	if err != nil {
		return "", err
	}
	if r.Token == "" {
		return "", fmt.Errorf("Login failed, empty token received")
	}
	return r.Token, nil
}

func (iu *imageUpdater) readLatestTag(imageName, token string) (string, error) {
	target := fmt.Sprintf("%s/v2/repositories/%s/tags/latest", iu.imageAPIServer, imageName)
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body == nil {
		return "", fmt.Errorf("Empty body returned")
	}

	defer func() {
		err = res.Body.Close()
	}()

	decoder := json.NewDecoder(res.Body)
	r := struct {
		Images []struct {
			Digest string
		}
	}{}
	err = decoder.Decode(&r)
	if err != nil {
		return "", err
	}

	if len(r.Images) == 0 {
		return "", fmt.Errorf("No images found for name %q", imageName)
	}
	if r.Images[0].Digest == "" {
		return "", fmt.Errorf("Remote retrieval failed, empty digest received %q", imageName)
	}
	log.Infof("Remote digest of %q is %q", imageName, r.Images[0].Digest)
	return r.Images[0].Digest, nil
}
