package examples

import (
	"errors"
	"net/http"
	"testing"

	"github.com/maxkopysow/go-httpclient/gohttp_mock"
)

func TestCreateRepo(t *testing.T) {
	t.Run("timeoutFromGithub", func(t *testing.T) {

		gohttp_mock.DeleteMocks()
		gohttp_mock.AddMock(gohttp_mock.Mock{
			Method:      http.MethodPost,
			Url:         "https://api.github.com/user/repos",
			RequestBody: `{"name":"test-repo","private":true}`,

			Error: errors.New("timeout from github"),
		})

		repository := Repository{
			Name:    "test-repo",
			Private: true,
		}

		repo, err := CreateRepo(repository)

		if repo != nil {
			t.Error("no repo expected when get a timeout from github")
		}

		if err == nil {
			t.Error("an error is expected when get a timeout from github")
		}

		if err.Error() != "timeout from github" {
			t.Error("invalid error message ")
		}
	})

	t.Run("NoError", func(t *testing.T) {

		gohttp_mock.DeleteMocks()
		gohttp_mock.AddMock(gohttp_mock.Mock{
			Method:      http.MethodPost,
			Url:         "https://api.github.com/user/repos",
			RequestBody: `{"name":"test-repo","private":true}`,

			ResponseStatusCode: http.StatusCreated,
			ResponseBody:       `{"id":123,"name":"test-repo"}`,
		})

		repository := Repository{
			Name:    "test-repo",
			Private: true,
		}

		repo, err := CreateRepo(repository)

		if err != nil {
			t.Error("no error expected when get valid response from github")
		}

		if repo == nil {
			t.Error("an valid repo was expected")
		}

		if repo.Name != repository.Name {
			t.Error("invalid repository name obtained from github")
		}
	})
}
