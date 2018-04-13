package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Client is a wrapper type around the moby client.
type Client struct {
	*client.Client
}

// New returns a pointer to a new instance of client.
func New() (*Client, error) {
	mobyClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Client{mobyClient}, nil
}

func (c *Client) ContainerRemoveByName(name string) error {
	containers, err := c.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			Limit: 1,
			Filters: filters.NewArgs(filters.KeyValuePair{
				Key:   "name",
				Value: name,
			}),
		},
	)
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return nil
	}
	return c.ContainerRemove(
		context.Background(),
		containers[0].ID,
		types.ContainerRemoveOptions{Force: true},
	)
}

// ImagePublish publishes a tagged image.
// ex. client.ImagePublish(true, "alpine", "localhost:5000/alpine")
func (c *Client) ImagePublish(pull bool, source, target string) (string, error) {
	rs := []io.Reader{}

	if pull {
		rc, err := c.ImagePull(
			context.Background(),
			source,
			types.ImagePullOptions{},
		)
		if err != nil {
			return "", err
		}
		rs = append(rs, rc)
	}

	err := c.ImageTag(
		context.Background(),
		source,
		target,
	)
	if err != nil {
		return "", err
	}

	rc, err := c.ImagePush(
		context.Background(),
		target,
		types.ImagePushOptions{
			RegistryAuth: "{}",
		},
	)
	if err != nil {
		return "", err
	}
	rs = append(rs, rc)

	b, err := ioutil.ReadAll(io.MultiReader(rs...))
	if err != nil {
		return "", err
	}

	out := strings.TrimSpace(string(b))

	// if last line is an error, it should be returned
	lines := strings.Split(out, "\n")
	if len(lines) > 0 {
		last := lines[len(lines)-1]
		if strings.Index(last, "error") != -1 {
			return out, fmt.Errorf(last)
		}
	}

	return out, nil
}
