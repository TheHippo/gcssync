package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"github.com/TheHippo/gcssync"
	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"os"
	"path/filepath"
)

const (
	_                = iota
	errorAuthInfo    = iota
	errorProjectInfo = iota
	errorClientInit  = iota
	errorListFiles   = iota
	errorUploadFiles = iota
	errorSyncFiles   = iota
)

const (
	version = "0.1.1"
)

const (
	scope       = storage.DevstorageFull_controlScope
	authURL     = "https://accounts.google.com/o/oauth2/auth"
	tokenURL    = "https://accounts.google.com/o/oauth2/token"
	entityName  = "allUsers"
	redirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

func main() {
	app := cli.NewApp()
	app.Name = "gcssync"
	app.Usage = "Sync files with Google Cloud Storage"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "cachefile",
			Value:  "cache.json",
			Usage:  "Cache file for caching auth tokens",
			EnvVar: "AUTH_CACHE_FILE",
		},
		cli.StringFlag{
			Name:   "bucketname, b",
			Value:  "",
			Usage:  "Name of bucket",
			EnvVar: "BUCKET_NAME",
		},
		cli.StringFlag{
			Name:   "projectid, p",
			Value:  "",
			Usage:  "Google project",
			EnvVar: "PROJECT_ID",
		},
		cli.StringFlag{
			Name:   "clientid, c",
			Value:  "",
			Usage:  "Auth client id",
			EnvVar: "AUTH_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "clientsecret, s",
			Value:  "",
			Usage:  "Client secrect",
			EnvVar: "AUTH_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "code",
			Value:  "",
			Usage:  "Authorization Code",
			EnvVar: "AUTH_CODE",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List remote files",
			Action:    listFiles,
		},
		{
			Name:      "upload",
			ShortName: "u",
			Usage:     "Upload a single file",
			Action:    uploadFile,
		},
		{
			Name:      "sync",
			ShortName: "s",
			Usage:     "Syncs a folder to a Google Cloudstorage bucket",
			Action:    syncFolder,
		},
	}
	app.Run(os.Args)
}

func generateOAuthConfig(c *cli.Context) (*oauth.Config, error) {
	clientID := c.GlobalString("clientid")
	if clientID == "" {
		return &oauth.Config{}, fmt.Errorf("Could not find Client ID")
	}
	clientSecret := c.GlobalString("clientsecret")
	if clientSecret == "" {
		return &oauth.Config{}, fmt.Errorf("Could not find Client Secret")
	}

	return &oauth.Config{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Scope:        scope,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		TokenCache:   oauth.CacheFile(c.GlobalString("cachefile")),
		RedirectURL:  redirectURL,
	}, nil
}

func generateServiceConfig(c *cli.Context) (*gcssync.ServiceConfig, error) {
	projectID := c.GlobalString("projectid")
	if projectID == "" {
		return &gcssync.ServiceConfig{}, fmt.Errorf("Could not find project id")
	}
	bucketName := c.GlobalString("bucketname")
	if bucketName == "" {
		return &gcssync.ServiceConfig{}, fmt.Errorf("Cloud not find bucket name")
	}
	return &gcssync.ServiceConfig{
		ProjectID:  projectID,
		BucketName: bucketName,
	}, nil
}

func getClient(c *cli.Context) *gcssync.Client {
	oauthConfig, err := generateOAuthConfig(c)
	if err != nil {
		fmt.Println("Missing auth informations", err.Error())
		os.Exit(errorAuthInfo)
	}
	serviceConfig, err := generateServiceConfig(c)
	if err != nil {
		fmt.Println("Missing project config", err.Error())
		os.Exit(errorProjectInfo)
	}

	client, err := gcssync.NewClient(oauthConfig, c.GlobalString("code"), serviceConfig)
	if err != nil {
		fmt.Println("Error initilizing client: ", err.Error())
		os.Exit(errorClientInit)
	}

	return client
}

func listFiles(c *cli.Context) {
	client := getClient(c)
	files, err := client.ListFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(errorListFiles)
		return
	}
	for _, object := range files {
		fmt.Printf("%s %s\n", object.Name, humanize.Bytes(object.Size))
	}
	fmt.Printf("Objects in %s - %d\n", client.GetBucketname(), len(files))
}

func uploadFile(c *cli.Context) {
	client := getClient(c)
	if len(c.Args()) != 2 {
		fmt.Println("Need local and remote name!")
		os.Exit(errorUploadFiles)
	}

	success, object, err := client.UploadFile(c.Args().Get(0), c.Args().Get(1))
	if !success {
		fmt.Println(err.Error())
		os.Exit(errorUploadFiles)
		return
	}

	fmt.Printf("Uploaded file to %s\n", client.GetBucketname())
	fmt.Printf("%s %s\n", object.Name, humanize.Bytes(object.Size))

}

func syncFolder(c *cli.Context) {
	client := getClient(c)
	var local, remote string
	switch len(c.Args()) {
	case 0:
		local = ""
		remote = ""
	case 1:
		local = c.Args().Get(0)
		remote = ""
	case 2:
		local = c.Args().Get(0)
		remote = c.Args().Get(1)
	default:
		fmt.Println("To many arguments")
		os.Exit(errorSyncFiles)
	}
	local, err := filepath.Abs(local)
	if err != nil {
		fmt.Println("Could not get absolute path")
		os.Exit(errorSyncFiles)
	}
	client.SyncFolder(local, remote)
}
