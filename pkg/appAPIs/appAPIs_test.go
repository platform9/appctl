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
		responseCode         int
		responseBody         map[string]string
		expectedResponseBody map[string]string
	}{
		"TestPass": {
			responseCode:         200,
			responseBody:         map[string]string{"Message": "Success"},
			expectedResponseBody: map[string]string{"Message": "Success"},
		},
		"TestBadRequest": {
			responseCode:         400,
			responseBody:         map[string]string{"Message": constants.BadRequest},
			expectedResponseBody: map[string]string{"Message": constants.BadRequest},
		},
		"TestAccessForbidden": {
			responseCode:         403,
			responseBody:         map[string]string{"Message": constants.AccessForbidden},
			expectedResponseBody: map[string]string{"Message": constants.AccessForbidden},
		},
		"TestMaxAppDeployLimit": {
			responseCode:         429,
			responseBody:         map[string]string{"Message": constants.MaxAppDeployLimit},
			expectedResponseBody: map[string]string{"Message": constants.MaxAppDeployLimit},
		},
		"TestInternalServerError": {
			responseCode:         500,
			responseBody:         map[string]string{"Message": constants.InternalServerError},
			expectedResponseBody: map[string]string{"Message": constants.InternalServerError},
		},
	}
	for testName, test := range listAppsCases {
		httpmock.Activate()
		httpmock.RegisterResponder(http.MethodGet, constants.APPURL, func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(test.responseCode, test.responseBody)
		})
		response, err := ListApps(dummyToken)
		if err != nil {
			if response != nil || err.Error() != test.responseBody["Message"] {
				logAPIFailure(t, err)
			}
		} else {
			if response["Message"].(string) != test.expectedResponseBody["Message"] {
				t.Errorf("test case: %s\t\tserver response: %v\n expected to be equal", testName, response)
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
			responseCode:      http.StatusOK,
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
				errMessage := fmt.Errorf("failed test case %s with error: %s\n", testName, err.Error())
				logAPIFailure(t, errMessage)
			}
		}
		httpmock.DeactivateAndReset()
	}
}

func TestGetAppByName(t *testing.T) {
	getAppByNameCases := map[string]struct {
		appName        string
		token          string
		expectedExists bool
	}{
		"Exists":    {appName: "app1", token: dummyToken, expectedExists: true},
		"NotExists": {appName: "app2", token: dummyToken, expectedExists: false},
	}
	for testName, test := range getAppByNameCases {
		httpmock.Activate()
		if test.expectedExists {
			httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/%s", constants.APPURL, test.appName), func(req *http.Request) (*http.Response, error) {
				body := map[string]interface{}{
					"Message": test.appName,
				}
				return httpmock.NewJsonResponse(200, body)
			})
		} else {
			httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s/%s", constants.APPURL, test.appName), func(req *http.Request) (*http.Response, error) {
				body := map[string]interface{}{}
				return httpmock.NewJsonResponse(400, body)
			})
		}
		appInfo, err := GetAppByName(test.appName, test.token)
		// if err != nil and the app is expected to exist, we fail the test
		if err != nil && test.expectedExists {
			errMessage := fmt.Errorf("failed test case %s with error: %s\n", testName, err.Error())
			logAPIFailure(t, errMessage)
		}
		t.Log(appInfo)
		httpmock.DeactivateAndReset()
	}
}

func TestDeleteAppByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("%s/%s", constants.APPURL, dummyAppName), func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, ""), nil
	})
	if err := DeleteAppByName(dummyAppName, dummyToken); err != nil {
		t.Log(err)
	}
	deleteAppByNameCases := map[string]struct {
		appName        string
		token          string
		expectedExists bool
	}{
		"Exists":    {appName: "app1", token: dummyToken, expectedExists: true},
		"NotExists": {appName: "app2", token: dummyToken, expectedExists: false},
	}

	for testName, test := range deleteAppByNameCases {
		httpmock.Activate()
		if test.expectedExists {
			httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("%s/%s", constants.APPURL, test.appName), func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, ""), nil
			})
		} else {
			httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("%s/%s", constants.APPURL, test.appName), func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(400, ""), nil
			})
		}
		err := DeleteAppByName(test.appName, test.token)
		if err != nil && test.expectedExists {
			errMessage := fmt.Errorf("failed test case %s with error: %s\n", testName, err.Error())
			logAPIFailure(t, errMessage)
		}
		httpmock.DeactivateAndReset()
	}
}
