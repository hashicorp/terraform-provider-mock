package provider

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func ProviderFactories(resources string) map[string]func() (tfprotov6.ProviderServer, error) {
	provider := NewForTesting("test", resources)()
	return map[string]func() (tfprotov6.ProviderServer, error){
		"tfcoremock": providerserver.NewProtocol6WithError(provider),
	}
}

func LoadFile(t *testing.T, file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("could not read file %s: %v", file, err.Error())
	}

	return string(data)
}

func CleanupTestingDirectories(t *testing.T) func() {
	return func() {
		files, err := os.ReadDir("terraform.resource")
		if err != nil {
			if os.IsNotExist(err) {
				return // Then it's fine.
			}

			t.Fatalf("could not read the resource directory for cleanup: %v", err)
		}
		defer os.Remove("terraform.resource")

		if len(files) != 0 {
			t.Fatalf("failed to tidy up after test")
		}
	}
}

func SaveResourceId(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		module := state.RootModule()
		rs, ok := module.Resources[name]
		if !ok {
			return errors.New("missing resource " + name)
		}

		*id = rs.Primary.Attributes["id"]
		return nil
	}
}

func CheckResourceIdChanged(name string, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		module := state.RootModule()
		rs, ok := module.Resources[name]
		if !ok {
			return errors.New("missing resource " + name)
		}

		if *id == rs.Primary.Attributes["id"] {
			return errors.New("id value for " + name + " has not changed")
		}
		return nil
	}
}
