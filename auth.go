package gcssync

import (
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/storage/v1"
	"fmt"
	"net/http"
)

type ServiceConfig struct {
	ProjectID  string
	BucketName string
}

type Client struct {
	service    *storage.Service
	projectId  string
	bucketName string
}

func NewClient(oauthConfig *oauth.Config, authCode string, serviceConfig *ServiceConfig) (*Client, error) {
	transport := &oauth.Transport{
		Config:    oauthConfig,
		Transport: http.DefaultTransport,
	}

	token, err := oauthConfig.TokenCache.Token()
	if err != nil {
		if authCode == "" {
			url := oauthConfig.AuthCodeURL("")
			fmt.Println("Visit URL to get a code then run again with -code=YOUR_CODE")
			fmt.Println(url)
			return &Client{}, fmt.Errorf("Could not get auth code")
		}
		token, err = transport.Exchange(authCode)
		if err != nil {
			return &Client{}, fmt.Errorf("Could not exchange token: %s", err.Error())
		}
		fmt.Printf("Token cached %s\n", oauthConfig.TokenCache)
	}

	transport.Token = token

	httpClient := transport.Client()

	service, err := storage.New(httpClient)

	if err != nil {
		return &Client{}, fmt.Errorf("Could not init client: %s", err.Error())
	}

	return &Client{
		service:    service,
		projectId:  serviceConfig.ProjectID,
		bucketName: serviceConfig.BucketName,
	}, nil
}