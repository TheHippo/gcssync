package gcssync

import (
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
)

func (c *Client) ListFiles() (*storage.Objects, error) {

	if res, err := c.service.Objects.List(c.bucketName).Do(); err == nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Could not fetch object list: %s", err.Error())
	}

}
