package appAPIs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/platform9/appctl/pkg/constants"
)

// if you're feeling brave, try your actual token :P
const dummyToken = "dummyToken"
const dummyAppName = "hello"

func TestListApps(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, constants.APPURL, func(req *http.Request) (*http.Response, error) {
		// TODO: mimic the logic
		return &http.Response{}, nil
	})
	ListApps(dummyToken)
}

func TestCreateApp(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, constants.APPURL, func(req *http.Request) (*http.Response, error) {
		// TODO: mimic the logic
		return &http.Response{}, nil
	})
	createAppCases := []struct {
		name     string
		image    string
		username string
		password string
		env      []string
		port     string
		token    string
	}{
		{name: dummyAppName, image: "private/someimage", username: "testUser", password: "testPassword", env: []string{"TEST_ENV=true", "PRODUCTION=false"}, port: "8888", token: dummyToken},
	}

	for i, test := range createAppCases {
		if err := CreateApp(test.name, test.image, test.username, test.password, test.env, test.port, test.token); err != nil {
			t.Errorf("failed test case %d with error: %s\n", i+1, err.Error())
		}
	}
}

func TestGetAppByName(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/%s", constants.APPURL, dummyAppName), func(req *http.Request) (*http.Response, error) {
		// TODO: mimic the logic
		var body = map[string]interface{}{}
		return httpmock.NewJsonResponse(200, body)
	})
	getAppByNameCases := []struct {
		name  string
		token string
	}{
		{name: dummyAppName, token: dummyToken},
	}
	for i, test := range getAppByNameCases {
		appInfo, err := GetAppByName(test.name, test.token)
		if err != nil {
			t.Errorf("failed test case %d with error: %s\n", i+1, err.Error())
		}
		t.Log(appInfo)
	}
}

func TestDeleteAppByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("%s/%s", constants.APPURL, dummyAppName), func(req *http.Request) (*http.Response, error) {
		// TODO: mimic the logic
		return httpmock.NewStringResponse(200, ""), nil
	})
}
