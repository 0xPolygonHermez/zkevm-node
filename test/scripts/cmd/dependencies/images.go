package dependencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// ImagesConfig is the configuration for the images updater.
type ImagesConfig struct {
	Names          []string
	TargetFilePath string
}

type imageUpdater struct {
	fs afero.Fs

	targetFilePath string

	dockerUsername string
	dockerPassword string

	images []string
}

func newImageUpdater(images []string, targetFilePath string) *imageUpdater {
	return &imageUpdater{
		fs: afero.NewOsFs(),

		targetFilePath: targetFilePath,

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
	err := iu.dockerLogin()
	if err != nil {
		return "", err
	}
	return iu.readLatestTag(imageName)
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

	return afero.WriteFile(iu.fs, targetFilePath, []byte(newContent), 0664) //nolint:gomnd
}

func (iu *imageUpdater) dockerLogin() error {
	c := exec.Command("docker", "login", "--username", iu.dockerUsername, "--password-stdin") // #nosec G204
	stdin, err := c.StdinPipe()
	if err != nil {
		return err
	}
	passReader := strings.NewReader(iu.dockerPassword)
	_, err = io.Copy(stdin, passReader)
	if err != nil {
		return err
	}
	err = stdin.Close()
	if err != nil {
		return err
	}

	return c.Run()
}

type result struct {
	RepoDigests []string
}

func (iu *imageUpdater) readLatestTag(imageName string) (string, error) {
	c := exec.Command("docker", "pull", imageName) // #nosec G204
	err := c.Run()
	if err != nil {
		return "", err
	}

	c = exec.Command("docker", "inspect", imageName) // #nosec G204
	output, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(output)
	decoder := json.NewDecoder(reader)

	r := struct {
		Results []result
	}{}
	err = decoder.Decode(&r.Results)
	if err != nil {
		return "", err
	}
	items := strings.Split(r.Results[0].RepoDigests[0], "@")
	const requiredItems = 2
	if len(items) < requiredItems {
		return "", fmt.Errorf("Returned image does not include digest %q", r.Results[0].RepoDigests[0])
	}
	log.Infof("Remote digest of %q is %q", imageName, items[1])
	return items[1], nil
}
