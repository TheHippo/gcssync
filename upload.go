package gcssync

import (
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"github.com/dustin/go-humanize"
	"os"
)

func (c *Client) UploadFile(localName, targetName string) {
	fmt.Println("ok")
	if _, err := os.Stat(localName); err != nil {
		fmt.Println("Local file not accessible")
		return
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
		fmt.Println("Could not open", localName)
		return
	}

	if res, err := c.service.Objects.Insert(c.bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Uploaded: %s %s %s", res.Name, humanize.Bytes(res.Size), res.SelfLink)
	} else {
		fmt.Println("Could not upload file: ", err.Error())
		return
	}
}
