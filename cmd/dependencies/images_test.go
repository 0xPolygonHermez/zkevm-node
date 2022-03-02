package dependencies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hermeznetwork/hermez-core/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_image_readCurrentDigest(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description      string
		input            string
		expectedOutput   string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description: "single container with matching image and existing digest",

			input: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
`,
			expectedOutput: "sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955",
		},
		{
			description: "multiple containers with matching image and existing digest",
			input: `version: '3'
services:
    hez-core:
        image: hezcore

    hez-network:
        image: imageorg/imagerepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
    hez-prover:
        image: hezprover
`,
			expectedOutput: "sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955",
		},
		{
			description: "single container with matching image and non-existing digest",

			input: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo:latest
`,
			expectedError:    true,
			expectedErrorMsg: "image \"imageorg/imagerepo\" not found",
		},
		{
			description: "single container with non-matching image",

			input: `version: '3'
services:
    hez-core:
        image: imageNeworg/imageNewrepo@sha256:f7bc95017b64a6dee98dd2e3a98cbea8c715de137d0f599b1b16f683c2dae955
`,
			expectedError:    true,
			expectedErrorMsg: "image \"imageorg/imagerepo\" not found",
		},
		{
			description:      "invalid yaml",
			input:            "not valid yaml",
			expectedError:    true,
			expectedErrorMsg: "yaml: unmarshal errors",
		},
	}

	const (
		defaultPath = "/a/b/dockerCompose.yml"
		imageName   = "imageorg/imagerepo"
	)

	subject := &imageUpdater{
		fs:             appFs,
		targetFilePath: defaultPath,
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, createFile(appFs, defaultPath, tc.input))
			defer func() {
				require.NoError(t, appFs.Remove(defaultPath))
			}()

			actualOutput, err := subject.readCurrentDigest(imageName)
			require.NoError(t, checkError(err, tc.expectedError, tc.expectedErrorMsg))
			require.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func Test_image_readRemoteDigest(t *testing.T) {
	const (
		imageName             = "imageorg/imagerepo"
		defaultDockerUsername = "user"
		defaultDockerPassword = "pass"
		defaultJWTToken       = "expectedToken"
		defaultDigest         = "expectedDigest"
	)

	tcs := []struct {
		description      string
		loginResponse    string
		getTagResponse   string
		expectedOutput   string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description:    "happy path",
			loginResponse:  defaultJWTToken,
			getTagResponse: defaultDigest,
			expectedOutput: defaultDigest,
		},
		{
			description:      "empty token received",
			loginResponse:    "",
			expectedError:    true,
			expectedErrorMsg: "Login failed, empty token received",
		},
		{
			description:      "empty digest received",
			loginResponse:    defaultJWTToken,
			getTagResponse:   "",
			expectedError:    true,
			expectedErrorMsg: "Remote retrieval failed, empty digest received",
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc(defaultLoginPattern, func(res http.ResponseWriter, req *http.Request) {
				require.Equal(t, http.MethodPost, req.Method)
				require.Equal(t, "application/json", req.Header.Get("Content-type"))

				decoder := json.NewDecoder(req.Body)
				var data struct {
					Username, Password string
				}
				require.NoError(t, decoder.Decode(&data))
				require.Equal(t, defaultDockerUsername, data.Username)
				require.Equal(t, defaultDockerPassword, data.Password)

				_, err := res.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, tc.loginResponse)))
				require.NoError(t, err)
			})
			mux.HandleFunc("/v2/repositories/imageorg/imagerepo/tags/latest", func(res http.ResponseWriter, req *http.Request) {
				if req.Header.Get("Authorization") != fmt.Sprintf("JWT %s", tc.loginResponse) {
					http.Error(res, "Invalid auth", http.StatusForbidden)
				}
				_, err := res.Write([]byte(fmt.Sprintf(`{"images":[{"digest":"%s"}]}`, tc.getTagResponse)))
				require.NoError(t, err)
			})
			ts := httptest.NewServer(mux)
			defer ts.Close()

			subject := &imageUpdater{
				imageAPIServer: ts.URL,
				dockerUsername: defaultDockerUsername,
				dockerPassword: defaultDockerPassword,
			}

			actualOutput, err := subject.readRemoteDigest(imageName)
			require.NoError(t, checkError(err, tc.expectedError, tc.expectedErrorMsg))
			require.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func Test_image_updateDigest(t *testing.T) {
	var appFs = afero.NewMemMapFs()

	tcs := []struct {
		description               string
		initialFileContents       string
		oldDigest                 string
		newDigest                 string
		expectedFinalFileContents string
		expectedError             bool
		expectedErrorMsg          string
	}{
		{
			description: "single container with matching image and existing digest",
			initialFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo@sha256:oldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo@sha256:newDigest`,
		},
		{
			description: "single container, not matching image",
			initialFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/anotherimagerepo@sha256:oldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/anotherimagerepo@sha256:oldDigest`,
		},
		{
			description: "single container with matching image, non-existing digest",
			initialFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo@sha256:veryOldDigest`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    hez-core:
        image: imageorg/imagerepo@sha256:veryOldDigest`,
		},
		{
			description: "multiple container with matching image and existing digest",
			initialFileContents: `version: '3'
services:
    hez-network:
        image: imageorg/networkImagerepo@sha256:oldDigest
    hez-core:
        image: imageorg/imagerepo@sha256:oldDigest
    hez-prover:
        image: imageorg/proverImagerepo@sha256:oldDigest
`,
			oldDigest: "sha256:oldDigest",
			newDigest: "sha256:newDigest",
			expectedFinalFileContents: `version: '3'
services:
    hez-network:
        image: imageorg/networkImagerepo@sha256:oldDigest
    hez-core:
        image: imageorg/imagerepo@sha256:newDigest
    hez-prover:
        image: imageorg/proverImagerepo@sha256:oldDigest
`,
		},
		{
			description:               "invalid yaml",
			initialFileContents:       "not valid yaml",
			oldDigest:                 "sha256:oldDigest",
			newDigest:                 "sha256:newDigest",
			expectedFinalFileContents: "not valid yaml",
		},
	}

	const (
		defaultPath = "/a/b/dockerCompose.yml"
		imageName   = "imageorg/imagerepo"
	)

	subject := &imageUpdater{
		fs:             appFs,
		targetFilePath: defaultPath,
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, createFile(appFs, defaultPath, tc.initialFileContents))
			defer func() {
				require.NoError(t, appFs.Remove(defaultPath))
			}()

			err := subject.updateDigest(imageName, tc.oldDigest, tc.newDigest)
			require.NoError(t, checkError(err, tc.expectedError, tc.expectedErrorMsg))
			actualFileContents, err := afero.ReadFile(appFs, defaultPath)
			require.NoError(t, err)
			require.Equal(t, tc.expectedFinalFileContents, string(actualFileContents))
		})
	}
}

func createFile(appFs afero.Fs, path, content string) error {
	f, err := appFs.Create(path)

	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Errorf("Could not close %q: %v", path, err)
		}
	}()

	_, err = f.WriteString(content)

	return err
}

func checkError(err error, expected bool, msg string) error {
	if !expected && err != nil {
		return fmt.Errorf("Unexpected error %v", err)
	}
	if expected {
		if err == nil {
			return fmt.Errorf("Expected error didn't happen")
		}
		if msg == "" {
			return fmt.Errorf("Expected error message not defined")
		}
		if !strings.HasPrefix(err.Error(), msg) {
			return fmt.Errorf("Wrong error, expected %q, got %q", msg, err.Error())
		}
	}
	return nil
}
