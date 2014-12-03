package gcssync

import (
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"os"
)

func (c *Client) UploadFile(localName, targetName string) (bool, *storage.Object, error) {
	if _, err := os.Stat(localName); err != nil {
		return false, nil, fmt.Errorf("Local file %s not accessible", localName)
	}
	object := &storage.Object{
		Name: targetName,
		Acl: []*storage.ObjectAccessControl{
			&storage.ObjectAccessControl{
				Entity: "allUsers",
				Bucket: c.bucketName,
				Object: targetName,
				Role:   "READER",
			},
		},
	}
	file, err := os.Open(localName)
	if err != nil {
		return false, nil, fmt.Errorf("Could not open %s", localName)
	}

	if res, err := c.service.Objects.Insert(c.bucketName, object).Media(file).Do(); err == nil {
		return true, res, nil
	}
	return false, nil, fmt.Errorf("Error while uploading file: %s", err.Error())
}
