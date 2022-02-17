package constants

import (
	"fmt"
	"os"
	"regexp"
)

// API Variables.
var (
	APPURL               = "http://fastpath.appctl.net:6112/v1/apps"
	TABLEFORMAT          = "NAME | URL | IMAGE | READY | AGE | REASON"
	DOMAIN               = "appctl.us.auth0.com"
	DEVICECODEURL        = "https://" + DOMAIN + "/oauth/device/code"
	CLIENTID             = "37cBgJP3K2yiYq2gamzEv1sH7vN2x6Z1"
	DEVICEREQUESTPAYLOAD = "client_id=" + CLIENTID + "&scope=" + GetAllScope()
	// Grant type is urlencoded
	GrantType = "grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code"
)

type ListAppInfo struct {
	Name         string
	URL          string
	Image        string
	Port         string
	ReadyStatus  string
	CreationTime string
	Reason       string
}

// Available SCOPES for auth0 access.
var AllScopes = []string{
	"openid",
	"profile",
	"email",
	"offline_access", // To get a refresh token.
	"create:clients", "delete:clients", "read:clients", "update:clients",
	"create:resource_servers", "delete:resource_servers", "read:resource_servers", "update:resource_servers",
	"create:roles", "delete:roles", "read:roles", "update:roles",
	"create:rules", "delete:rules", "read:rules", "update:rules",
	"create:users", "delete:users", "read:users", "update:users",
	"read:branding", "update:branding",
	"read:email_templates", "update:email_templates",
	"read:connections", "update:connections",
	"read:client_keys", "read:logs", "read:tenant_settings",
	"read:custom_domains", "create:custom_domains", "update:custom_domains", "delete:custom_domains",
	"read:anomaly_blocks", "delete:anomaly_blocks",
	"create:log_streams", "delete:log_streams", "read:log_streams", "update:log_streams",
	"create:actions", "delete:actions", "read:actions", "update:actions",
	"create:organizations", "delete:organizations", "read:organizations", "update:organizations",
	"read:prompts", "update:prompts",
}

// Only required scopes for IDToken generation.
var RequiredScopes = AllScopes[:4]

const (
	// Time to wait to get app deployed.
	APPDEPLOYINTERVAL = 5

	// Token poll interval
	TOKENPOLLINTERVAL = 5

	// Fetch secure app endpoint.
	SECUREENDPOINT = 2

	// Maximum app deployed status code.
	MaxAppDeployStatusCode = "429"

	// HTTPS string
	HTTPS = "https"
)

//Configfile
var (
	HOMEDIR, _ = os.UserHomeDir()
	//Config Dir for pf9
	CONFIGDIR      = HOMEDIR + "/.config/pf9"
	CONFIGFILE     = "config.json"
	CONFIGFILEPATH = CONFIGDIR + "/" + CONFIGFILE
)

func GetAllScope() string {
	var scope string
	for _, val := range RequiredScopes {
		scope += val + " "
	}
	return (scope)
}

var (
	// Valid App Name to deploy.
	ValidAppNameRegex = fmt.Sprintf(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
)

// Validate a regex.
func RegexValidate(name string, regex string) bool {
	Regex := regexp.MustCompile(regex)
	return Regex.MatchString(name)
}

// Error Messages.
var (
	ConnectionRefused      = "connection refused"
	NetworkUnReachable     = "Network is unreachable"
	InternetConnectivity   = "Please check your internet connectivity and try again."
	BackendServerDown      = "Backend server is down. Please try later!!"
	AccessForbidden        = "Access Forbidden."
	MaxAppDeployLimit      = "Maximum apps deploy limit reached!!"
	InvalidImage           = "Unable to fetch image"
	MaxAppDeployLimitError = "Maximum apps deploy limit reached!!"
	FailedToParseImage     = "Failed to parse image"
)

const (
	CLIVersion          = "appctl version: v1.2"
	UTCClusterTimeStamp = "2006-01-02T15:04:05Z"
)
