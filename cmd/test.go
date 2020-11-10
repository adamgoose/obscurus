package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:    "test {obscura_url}",
	Short:  "Runs integration tests",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		o, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		// Create a secret
		secret := bytes.NewBufferString("this is a really secret secret")
		o.Path = "/api/v1/secrets"
		r, err := http.Post(o.String(), "text/plain", secret)
		if err != nil {
			return err
		}
		defer r.Body.Close()

		// Get the token
		data := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&data)
		token := data["token"].(string)

		// Lookup the token
		o.Path = "/api/v1/secrets/" + token
		_, err = http.Get(o.String())
		if err != nil {
			return err
		}
		defer r.Body.Close()

		// Read the token
		o.Path = "/api/v1/secrets/" + token + "/value"
		r, err = http.Get(o.String())
		if err != nil {
			return err
		}
		defer r.Body.Close()

		// Get the secret
		b := bytes.NewBuffer(nil)
		io.Copy(b, r.Body)

		// Make sure the secret matches
		if b.String() != "this is a really secret secret" {
			return errors.New("The secret didn't match")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
