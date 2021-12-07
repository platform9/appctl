package constants

var (
	APPURL               = "http://127.0.0.1:6112/v1/apps"
	TABLEFORMAT          = "NAME | URL | IMAGE | NAMESPACE | CREATIONTIME"
	DOMAIN               = "platform9.us.auth0.com"
	DEVICECODEURL        = "https://" + DOMAIN + "/oauth/device/code"
	CLIENTID             = "HEVMcEBvvQ1wnRmzOxlShZXvjp07bnMz"
	DEVICEREQUESTPAYLOAD = "client_id=" + CLIENTID + "&scope=" + GetAllScope()

	// SCOPES for auth0 access token.
	SCOPEIDS    = "openid offline_access "
	SCOPECLIENT = "create:clients delete:clients read:clients update:clients delete:resource_servers delete:resource_servers delete:resource_servers "
	SCOPESERVER = "create:resource_servers delete:resource_servers read:resource_servers update:resource_servers "
	SCOPEROLES  = "create:roles delete:roles read:roles update:roles "
	SCOPERULES  = "create:rules delete:rules read:rules update:rules "
	SCOPEUSERS  = "create:users delete:users read:users update:users "
	// Scope for branding, email_template, connections.
	SCOPEBRANDTEMPCON = "read:branding update:branding read:email_templates update:email_templates read:connections update:connections "
	SCOPEDOMAIN       = "read:custom_domains create:custom_domains update:custom_domains delete:custom_domains "
	SCOPEBLOCKSTREAM  = "read:anomaly_blocks delete:anomaly_blocks create:log_streams delete:log_streams read:log_streams update:log_streams "
	SCOPEACTION       = "create:actions delete:actions read:actions update:actions "
	SCOPEORG          = "create:organizations delete:organizations read:organizations update:organizations "
	SCOPEOTHER        = "read:client_keys read:logs read:tenant_settings read:prompts update:prompts "
)

const (
	// Time to wait to get app deployed.
	APPDEPLOYINTERVAL = 5

	// Token poll interval
	TOKENPOLLINTERVAL = 5
)

func GetAllScope() string {
	return (SCOPEIDS + SCOPECLIENT + SCOPESERVER + SCOPEROLES + SCOPERULES + SCOPEBRANDTEMPCON + SCOPEDOMAIN + SCOPEBLOCKSTREAM + SCOPEACTION + SCOPEORG + SCOPEOTHER)
}
