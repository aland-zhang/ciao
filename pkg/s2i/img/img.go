package img

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caicloud/ciao/pkg/types"
)

const (
	prefix     = "kubeflow-kernel-code."
	codeFile   = "code.py"
	dockerFile = "Dockerfile"
	imageOwner = "gaocegege"
)

// Client is the type for using img.
type Client struct {
}

// New creates a new Client.
func New() *Client {
	return &Client{}
}

// SourceToImage converts the code to the image.
func (c Client) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	// This is a hack to let kubernetes do not pull from docker registry.
	imageName := fmt.Sprintf("%s:v1", filepath.Join(imageOwner, parameter.GenerateName))

	dir, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(dir, codeFile), []byte(code), 0666)
	if err != nil {
		return "", err
	}

	if err = c.writeDockerfile(dir, parameter); err != nil {
		return "", err
	}

	cmd := exec.Command("img", "build", "-t", imageName, dir)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("[kubeflow] Failed to build the image: %s", string(output))
		return "", err
	}

	fmt.Printf("[kubeflow] Pushing the image...\n")
	cmd = exec.Command("img", "push", imageName)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("[kubeflow] Failed to push the image: %s", string(output))
		return "", err
	}
	return imageName, nil
}

func (c Client) writeDockerfile(dir string, parameter *types.Parameter) error {
	var template string
	switch parameter.Framework {
	case types.FrameworkTypeTensorFlow:
		template = tensorflowTemplate
	case types.FrameworkTypePyTorch:
		template = pytorchTemplate
	}
	return ioutil.WriteFile(filepath.Join(dir, dockerFile), []byte(template), 0666)
}