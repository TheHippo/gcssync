package gcssync

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func (c *Client) ListFiles() {

	if res, err := c.service.Objects.List(c.bucketName).Do(); err == nil {
		fmt.Printf("Objects in bucket %s (%s):\n", c.bucketName, c.projectId)
		// fmt.Println(res.NextPageToken)
		for _, object := range res.Items {
			fmt.Printf("%s %s %s\n", object.Name, humanize.Bytes(object.Size), object.SelfLink)
		}
	} else {
		fmt.Printf("Objects.List failed: %s\n", err.Error())
	}

}
