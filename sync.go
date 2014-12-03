package gcssync

import (
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"github.com/dustin/go-humanize"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	fileWalkerBufferSize = 20
	fileWalkerEstimate   = 500
	uploadGoroutinesNum  = 10
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

	// this is unnecessary, but left there in case filters for files will be implemented
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
	var localFiles []fileInfo
	var listsFetched sync.WaitGroup

	listsFetched.Add(1)
	go func() {
		localFiles = getLocalFiles(from)
		listsFetched.Done()
	}()

	var objects *storage.Objects
	var err error

	listsFetched.Add(1)
	go func() {
		objects, err = c.ListFiles()
		if err != nil {
			fmt.Println(err)
		}
		listsFetched.Done()
	}()

	listsFetched.Wait()

	fmt.Printf("Found %d local files\n", len(localFiles))
	fmt.Printf("Found %d remote files\n", len(objects.Items))

	remoteCache := make(map[string]time.Time, len(objects.Items))

	for _, object := range objects.Items {
		time, err := time.Parse(time.RFC3339, object.Updated)
		if err != nil {
			continue
		}
		remoteCache[object.Name] = time
	}

	already := 0
	created := 0
	update := 0

	toDo := make([]fileInfo, 0, len(localFiles))

	for _, file := range localFiles {
		remoteTime, exists := remoteCache[file.path]
		if exists {
			already++
			if file.info.ModTime().After(remoteTime) {
				toDo = append(toDo, file)
				update++
			}
		} else {
			created++
			toDo = append(toDo, file)
		}
	}

	var uploadDone sync.WaitGroup
	uploadFile := make(chan fileInfo, uploadGoroutinesNum)
	uploadDone.Add(1)
	go func() {
		for _, file := range toDo {
			uploadFile <- file
		}
		close(uploadFile)
		uploadDone.Done()
	}()

	for i := 0; i < uploadGoroutinesNum; i++ {
		uploadDone.Add(1)
		go func() {
			for {
				file, more := <-uploadFile
				if more {
					success, object, err := c.UploadFile(filepath.Join(from, file.path), file.path)
					if err != nil || success == false {
						fmt.Println(err)
					} else {
						fmt.Printf("Uploaded %s %s\n", object.Name, humanize.Bytes(object.Size))
					}
				} else {
					uploadDone.Done()
					return
				}
			}
		}()
	}

	uploadDone.Wait()

	fmt.Println(len(toDo))
	fmt.Println(already, created, update)
}
