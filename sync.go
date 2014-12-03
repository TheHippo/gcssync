package gcssync

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	fileWalkerBufferSize = 20
	fileWalkerEstimate   = 500
)

type fileInfo struct {
	info os.FileInfo
	path string
}

func getLocalFiles(dirname string) []fileInfo {
	rawFiles := make(chan fileInfo, fileWalkerBufferSize)
	var done sync.WaitGroup

	done.Add(1)
	go func() {
		filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
			if err == nil && !f.IsDir() {
				rel, _ := filepath.Rel(dirname, path)
				rawFiles <- fileInfo{
					info: f,
					path: rel,
				}
			}
			return nil
		})
		close(rawFiles)
		done.Done()
	}()

	result := make([]fileInfo, 0, fileWalkerEstimate)
	done.Add(1)
	go func() {
		for {
			f, more := <-rawFiles
			if more {
				result = append(result, f)
			} else {
				done.Done()
				return
			}
		}
	}()
	done.Wait()
	return result
}

func (c *Client) SyncFolder(from, to string) {
	fmt.Println(from, to)
	localfiles := getLocalFiles(from)
	fmt.Printf("Found %d local files\n", len(localfiles))
}
