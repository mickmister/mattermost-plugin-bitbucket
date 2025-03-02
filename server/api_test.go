package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"

	"github.com/mattermost/mattermost-plugin-bitbucket/server/testutils"
)

func TestPlugin_ServeHTTP(t *testing.T) {
	httpTestJSON := testutils.HTTPTest{
		T:       t,
		Encoder: testutils.EncodeJSON,
	}

	httpTestString := testutils.HTTPTest{
		T:       t,
		Encoder: testutils.EncodeString,
	}

	for name, test := range map[string]struct {
		httpTest         testutils.HTTPTest
		request          testutils.Request
		expectedResponse testutils.ExpectedResponse
		userID           string
	}{
		"unauthorized test json": {
			httpTest: httpTestJSON,
			request: testutils.Request{
				Method: http.MethodPost,
				URL:    "/api/v1/todo",
				Body:   nil,
			},
			expectedResponse: testutils.ExpectedResponse{
				StatusCode:   http.StatusUnauthorized,
				ResponseType: testutils.ContentTypeJSON,
				Body:         APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized},
			},
			userID: "",
		}, "unauthorized test http": {
			httpTest: httpTestString,
			request: testutils.Request{
				Method: http.MethodGet,
				URL:    "/api/v1/reviews",
				Body:   nil,
			},
			expectedResponse: testutils.ExpectedResponse{
				StatusCode:   http.StatusUnauthorized,
				ResponseType: testutils.ContentTypePlain,
				Body:         "Not authorized\n",
			},
			userID: "",
		},
	} {
		t.Run(name, func(t *testing.T) {
			p := NewPlugin()
			p.setConfiguration(
				&Configuration{
					BitbucketOrg:               "mockOrg",
					BitbucketOAuthClientID:     "mockID",
					BitbucketOAuthClientSecret: "mockSecret",
					WebhookSecret:              "mockSecret",
					EncryptionKey:              "mockKey",
				})
			p.initializeAPI()
			p.SetAPI(&plugintest.API{})

			req := test.httpTest.CreateHTTPRequest(test.request)
			req.Header.Add("Mattermost-User-ID", test.userID)
			rr := httptest.NewRecorder()
			p.ServeHTTP(&plugin.Context{}, rr, req)
			test.httpTest.CompareHTTPResponse(rr, test.expectedResponse)
		})
	}
}

func TestGetToken(t *testing.T) {
	httpTestString := testutils.HTTPTest{
		T:       t,
		Encoder: testutils.EncodeString,
	}

	for name, test := range map[string]struct {
		httpTest         testutils.HTTPTest
		request          testutils.Request
		context          *plugin.Context
		expectedResponse testutils.ExpectedResponse
	}{
		"not authorized": {
			httpTest: httpTestString,
			request: testutils.Request{
				Method: http.MethodGet,
				URL:    "/api/v1/token",
				Body:   nil,
			},
			context: &plugin.Context{},
			expectedResponse: testutils.ExpectedResponse{
				StatusCode:   http.StatusUnauthorized,
				ResponseType: testutils.ContentTypePlain,
				Body:         "Not authorized\n",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			p := NewPlugin()
			p.setConfiguration(
				&Configuration{
					BitbucketOrg:               "mockOrg",
					BitbucketOAuthClientID:     "mockID",
					BitbucketOAuthClientSecret: "mockSecret",
					WebhookSecret:              "mockSecret",
					EncryptionKey:              "mockKey",
				})
			p.initializeAPI()

			p.SetAPI(&plugintest.API{})

			req := test.httpTest.CreateHTTPRequest(test.request)
			rr := httptest.NewRecorder()

			p.ServeHTTP(test.context, rr, req)

			test.httpTest.CompareHTTPResponse(rr, test.expectedResponse)
		})
	}
}
