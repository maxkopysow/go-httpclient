package example

import (
	"fmt"

	"github.com/maxkopysow/go-httpclient.git/gohttp"
)

var (
	githubClient = getGithubClient()
)

func getGithubClient() gohttp.Client {

	client := gohttp.NewBuilder().
		DisableTimeouts(true).
		SetMaxIdleConnections(5).
		Build()
	return client
}

type User struct {
	FirstName string `json:"first_name"`
}

func main() {
	getUrl()
	getUrl()
	getUrl()
	createUser(User{FirstName: "First"})

}

func getUrl() {

	response, err := githubClient.Get("https://api.github.com", nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode())

}

func createUser(user User) {
	response, err := githubClient.Post("https://api.github.com", nil, user)

	if err != nil {
		panic(err)
	}

	fmt.Println(response.StatusCode())

}
