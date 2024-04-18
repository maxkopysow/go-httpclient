package examples

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/maxkopysow/go-httpclient/gohttp_mock"
)

func TestMain(m *testing.M) {

	fmt.Println("About to start test cases for package 'exapmles' ")

	gohttp_mock.MockupServer.Start()

	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	//INIT

	t.Run("TestErrorFetchingGithub", func(t *testing.T) {

		gohttp_mock.MockupServer.DeleteMocks()

		gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
			Method: http.MethodGet,
			Url:    "https://api.github.com",
			Error:  errors.New("timeout getting github endpoints"),
		})

		//EXC
		endpoints, err := GetEndpoints()

		if endpoints != nil {
			t.Error("no endpoints expected")
		}
		if err == nil {
			t.Error("error was expected")
		}

		if err.Error() != "timeout getting github endpoints" {
			t.Error("invalid error message recieved")
		}
	})

	t.Run("TestErrorUnmarshallResponseBody", func(t *testing.T) {
		gohttp_mock.MockupServer.DeleteMocks()

		gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
			Method:             http.MethodGet,
			Url:                "https://api.github.com",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       `{"current_user_url" : 123}`,
		})

		//EXC
		endpoints, err := GetEndpoints()

		if endpoints != nil {
			t.Error("no endpoints expected")
		}
		if err == nil {
			t.Error("error was expected")
		}

		if !strings.Contains(err.Error(), "cannot unmarshal number into Go struct field") {
			t.Error("invalid error message recieved")
			t.Error(err.Error())

		}
	})
	t.Run("TestNoError", func(t *testing.T) {
		gohttp_mock.MockupServer.DeleteMocks()

		gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
			Method:             http.MethodGet,
			Url:                "https://api.github.com",
			ResponseStatusCode: http.StatusOK,
			ResponseBody:       `{"current_user_url" : "https://api.github.com/user"}`,
		})
		//EXC
		endpoints, err := GetEndpoints()

		if err != nil {
			t.Errorf("no error expected we got %s", err.Error())
		}

		if endpoints == nil {
			t.Error("endpoints were expected and we got nil")
		}

		// if !strings.Contains(err.Error(), "cannot unmarshall number into Go struct field") {
		// 	t.Error("invalid error message recieved")
		// }

	})

}
