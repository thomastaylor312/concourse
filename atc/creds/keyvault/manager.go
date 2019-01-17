package keyvault

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/concourse/concourse/atc/creds"
)

// KeyVaultManager is an Azure Key Vault implementation of the creds.Manager
// interface
type KeyVaultManager struct {
	ServicePrincipalID  string `long:"service-principal-id" description:"Azure service principal ID. It should have read and list access to the Key Vault you are trying to access"`
	ServicePrincipalKey string `long:"service-principal-key" description:"Azure service principle key or password"`
	TenantID            string `long:"tenant-id" description:"The ID of the Azure AD tenant your service principal is part of"`
	KeyVaultURL         string `long:"key-vault-url" description:"The URL of the Key Vault you wish to use"`
	KeyPrefix           string `long:"key-prefix" default:"concourse" description:"Value under which to prefix key names."`
	Environment         string `long:"environment" default:"AzurePublicCloud" description:"The Azure environment to use. If you need to change this from the default, you'll know it"`

	Reader SecretReader
}

// MarshalJSON is the custom JSON marshalling function for this manager
func (manager *KeyVaultManager) MarshalJSON() ([]byte, error) {
	health, err := manager.Health()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&map[string]interface{}{
		"service_principal_id": manager.ServicePrincipalID,
		"tenant_id":            manager.TenantID,
		"key_vault_url":        manager.KeyVaultURL,
		"key_prefix":           manager.KeyPrefix,
		"environment":          manager.Environment,
		"health":               health,
	})
}

// Init creates and configures the proper Key Vault client
func (manager *KeyVaultManager) Init(log lager.Logger) error {
	// Create the keyvault client with the proper credentials
	conf := auth.NewClientCredentialsConfig(manager.ServicePrincipalID, manager.ServicePrincipalKey, manager.TenantID)

	// Grab the proper endpoint configs and set them
	env, err := azure.EnvironmentFromName(manager.Environment)
	if err != nil {
		return err
	}
	conf.AADEndpoint = env.ActiveDirectoryEndpoint
	// The default endpoints sometimes have a trailing slash, which messes
	// things up, so remove it
	conf.Resource = strings.TrimSuffix(env.KeyVaultEndpoint, "/")

	kv := keyvault.New()
	authz, err := conf.Authorizer()
	if err != nil {
		return fmt.Errorf("unable to authorize with Azure API: %s", err)
	}
	kv.Authorizer = authz

	manager.Reader = NewKeyVaultReader(kv, manager.KeyVaultURL)

	return nil
}

// Health checks if the manager can properly access the Key Vault
func (manager *KeyVaultManager) Health() (*creds.HealthResponse, error) {
	health := &creds.HealthResponse{
		Method: "/health",
	}

	// Try to fetch a non-existent secret. It should not return an error for a
	// non-existent secret, so if it does, we know something is up
	_, _, err := manager.Reader.Get("i_should_never_exist")
	if err != nil {
		health.Error = err.Error()
		return health, nil
	}

	health.Response = map[string]string{
		"status": "UP",
	}

	return health, nil
}

// IsConfigured returns a boolean indicating if the manager has the proper
// configuration. This is just a basic check, and is mostly done to make sure
// there isn't an empty configuration. More in depth checking is done in the
// Validate function
func (manager *KeyVaultManager) IsConfigured() bool {
	return manager.ServicePrincipalID != ""
}

// Validate returns an error if all of the proper configuration is not in place
func (manager *KeyVaultManager) Validate() error {
	if manager.ServicePrincipalID == "" || manager.ServicePrincipalKey == "" {
		return fmt.Errorf("must provide a service principal ID and key")
	}

	if manager.KeyVaultURL == "" {
		return fmt.Errorf("must provide the key vault URL")
	}

	if manager.TenantID == "" {
		return fmt.Errorf("must provide the tenant ID")
	}

	return nil
}

// NewVariablesFactory implements the manager interface and returns a
// VariablesFactory implementation for Azure Key Vault
func (manager *KeyVaultManager) NewVariablesFactory(log lager.Logger) (creds.VariablesFactory, error) {
	return NewKeyVaultFactory(log, manager.Reader, manager.KeyPrefix), nil
}