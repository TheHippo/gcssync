package gcssync

import (
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
)

const (
	objectCapacity = 1000
)

// ListFiles returns all files in the connected bucket
func (c *Client) ListFiles() ([]*storage.Object, error) {
	result := make([]*storage.Object, 0, objectCapacity)
	next := ""
	var res *storage.Objects
	var err error
	for {
		if next != "" {
			res, err = c.service.Objects.List(c.bucketName).PageToken(next).Do()
		} else {
			res, err = c.service.Objects.List(c.bucketName).Do()
		}
		if err != nil {
			return nil, fmt.Errorf("Could not fetch object list: %s", err.Error())
		}
		next = res.NextPageToken
		result = append(result, res.Items...)
		if next == "" {
			break
		}
	}
	return result, nil

}
