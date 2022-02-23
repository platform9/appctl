package appAPIs

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/platform9/appctl/pkg/constants"
)

// if you're feeling brave, try your actual token :P
const dummyToken = "dummyToken"
const dummyAppName = "hello"

func logAPIFailure(t *testing.T, err error) {
	t.Errorf("failed with error: %s\n", err.Error())
}

func TestListApps(t *testing.T) {
	// These test cases are keyed according to whether they're expected to pass or result in an error
	// The `ListApps` function returns err, nil in case of an error
	listAppsCases := map[string]struct {
		responseCode int
		responseBody map[string]string
	}{
		"TestPass": {responseCode: 200,
			responseBody: map[string]string{"Message": "Success"},
		},
		"TestBadRequest": {
			responseCode: 400,
			responseBody: map[string]string{"Message": constants.BadRequest},
		},
		"TestAccessForbidden": {
			responseCode: 403,
			responseBody: map[string]string{"Message": constants.AccessForbidden},
		},
		"TestMaxAppDeployLimit": {
			responseCode: 429,
			responseBody: map[string]string{"Message": constants.MaxAppDeployLimit},
		},
		"TestInternalServerError": {
			responseCode: 500,
			responseBody: map[string]string{"Message": constants.InternalServerError},
		},
	}
	for testName, test := range listAppsCases {
		httpmock.Activate()
		httpmock.RegisterResponder(http.MethodGet, constants.APPURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(test.responseCode, test.responseBody)
		})
		response, err := ListApps(dummyToken)
		if err != nil {
			if response != nil {
				t.Errorf("failed case %s\n", testName)
			}
			if err.Error() != test.responseBody["Message"] {
				logAPIFailure(t, err)
			}
		} else {
			if response["Message"].(string) != test.responseBody["Message"] {
				t.Errorf("test case: %s\t\tserver response: %v\n", testName, response)
			}
		}
		httpmock.DeactivateAndReset()
	}
}

func TestCreateApp(t *testing.T) {
	createAppCases := map[string]struct {
		name              string
		image             string
		username          string
		password          string
		env               []string
		port              string
		token             string
		responseCode      int
		expectedErrPrefix string
	}{
		"TestPrivateImage": {
			name:              "privatedummyapp",
			image:             "private/someimage",
			username:          "testUser",
			password:          "testPassword",
			env:               []string{"TEST_ENV=true", "PRODUCTION=false"},
			port:              "8888",
			token:             dummyToken,
			responseCode:      http.StatusOK,
			expectedErrPrefix: "",
		},
		"TestPublicImage": {
			name:              "noAuth",
			image:             "public/someimage",
			username:          "",
			password:          "",
			env:               []string{"TEST_ENV=true", "PRODUCTION=false"},
			port:              "8888",
			token:             dummyToken,
			responseCode:      http.StatusOK,
			expectedErrPrefix: "",
		},
		"TestNoEnv": {
			name:              "noEnv",
			image:             "public/someimage",
			username:          "",
			password:          "",
			env:               []string{},
			port:              "8888",
			token:             dummyToken,
			expectedErrPrefix: "",
		},
		"TestFailBadRequest": {
			name:              "noEnvFail",
			image:             "public/someimage",
			username:          "",
			password:          "",
			env:               []string{},
			port:              "8888",
			token:             dummyToken,
			responseCode:      400,
			expectedErrPrefix: constants.BadRequest,
		},
	}

	for testName, test := range createAppCases {
		httpmock.Activate()
		httpmock.RegisterResponder(http.MethodPost, constants.APPURL, func(req *http.Request) (*http.Response, error) {

			return httpmock.NewJsonResponse(test.responseCode, map[string]string{
				"Message": testName,
			})
		})
		err := CreateApp(test.name, test.image, test.username, test.password, test.env, test.port, test.token)
		if err != nil {
			if !strings.HasPrefix(err.Error(), test.expectedErrPrefix) {
				t.Errorf("failed test case %s with error: %s\n", testName, err.Error())
			}
		}
		httpmock.DeactivateAndReset()
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
