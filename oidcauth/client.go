package oidcauth

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Setup ensures the user has a valid id token for OIDC requests
func Setup() error {
	var (
		idToken      = viper.GetString("id_token")
		accessToken  = viper.GetString("access_token")
		refreshToken = viper.GetString("refresh_token")
		clientID     = viper.GetString("client_id")
		clientSecret = viper.GetString("client_secret")
		err          error
	)

	manager, err := NewManager(context.Background(), clientID, clientSecret)
	if err != nil {
		return err
	}
	idToken, accessToken, refreshToken, err = EnsureValidTokens(
		manager,
		idToken,
		accessToken,
		refreshToken,
	)
	if err != nil {
		return err
	}
	viper.Set("id_token", idToken)
	viper.Set("access_token", accessToken)
	viper.Set("refresh_token", refreshToken)

	// Write config
	{
		data, err := yaml.Marshal(viper.AllSettings())
		if err != nil {
			return fmt.Errorf("Error marshaling viper settings: %q", err)
		}

		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("Error getting home dir: %q", err)
		}

		filename := filepath.Join(home, ".keyman.yaml")

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return fmt.Errorf("Error writing config file: %q", err)
		}
	}

	return nil
}
