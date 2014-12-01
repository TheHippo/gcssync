package main

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/consulted/gcssync"
	"os"
)

const (
	_                = iota
	errorAuthInfo    = iota
	errorProjectInfo = iota
	errorClientInit  = iota
	errorUploadFiles = iota
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
	}
	app.Run(os.Args)
}

func generateOAuthConfig(c *cli.Context) (*oauth.Config, error) {
	clientId := c.GlobalString("clientid")
	if clientId == "" {
		return &oauth.Config{}, fmt.Errorf("Could not find Client ID")
	}
	clientSecret := c.GlobalString("clientsecret")
	if clientSecret == "" {
		return &oauth.Config{}, fmt.Errorf("Could not find Client Secret")
	}

	return &oauth.Config{
		ClientId:     clientId,
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
	client.ListFiles()
}

func uploadFile(c *cli.Context) {
	client := getClient(c)
	if len(c.Args()) != 2 {
		fmt.Println("Need local and remote name!")
		os.Exit(errorUploadFiles)
	}

	client.UploadFile(c.Args().Get(0), c.Args().Get(1))
}
