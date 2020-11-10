package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adamgoose/obscurus/lib"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "obscurus",
	Short: "Secret sharing with Vault!",
	RunE: func(cmd *cobra.Command, args []string) error {
		h := lib.NewHandler(lib.GetVaultClient())

		port := viper.GetInt("port")
		log.Printf("Listening on http://127.0.0.1:%d\n", port)
		return http.ListenAndServe(fmt.Sprintf(":%d", port), h)
	},
}

// Execute the command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	godotenv.Load()
	viper.AutomaticEnv()

	viper.SetDefault("port", "8000")
}
