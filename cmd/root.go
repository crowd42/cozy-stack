package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cozy/cozy-stack/client"
	"github.com/cozy/cozy-stack/client/request"
	"github.com/cozy/cozy-stack/pkg/config"
	"github.com/cozy/cozy-stack/pkg/permissions"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DefaultStorageDir is the default directory name in which data
// is stored relatively to the cozy-stack binary.
const DefaultStorageDir = "storage"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cozy-stack",
	Short: "cozy-stack is the main command",
	Long: `Cozy is a platform that brings all your web services in the same private space.
With it, your web apps and your devices can share data easily, providing you
with a new experience. You can install Cozy on your own hardware where no one
profiles you.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return Configure()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Display the usage/help by default
		return cmd.Help()
	},
	// Do not display usage on error
	SilenceUsage: true,
	// We have our own way to display error messages
	SilenceErrors: true,
}

var cfgFile string

func newClient(domain, scope string) *client.Client {
	// For the CLI client, we rely on the admin APIs to generate a CLI token.
	// We may want in the future rely on OAuth to handle the permissions with
	// more granularity.
	c := newAdminClient()
	token, err := c.GetToken(&client.TokenOptions{
		Domain:   domain,
		Subject:  "CLI",
		Audience: permissions.CLIAudience,
		Scope:    []string{scope},
	})
	if err != nil {
		log.Errorf("Could not generate access to domain %s", domain)
		log.Error(err)
		os.Exit(1)
	}
	return &client.Client{
		Domain:     domain,
		Authorizer: &request.BearerAuthorizer{Token: token},
	}
}

func newAdminClient() *client.Client {
	var pass []byte
	if !config.IsDevRelease() {
		pass = []byte(os.Getenv("COZY_ADMIN_PASSWORD"))
		if len(pass) == 0 {
			var err error
			fmt.Printf("Password:")
			pass, err = gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
		}
	}
	return &client.Client{
		Domain:     config.AdminServerAddr(),
		Authorizer: &request.BasicAuthorizer{Password: string(pass)},
	}
}

func init() {
	flags := RootCmd.PersistentFlags()
	flags.StringVarP(&cfgFile, "config", "c", "", "configuration file (default \"$HOME/.cozy.yaml\")")

	flags.String("host", "localhost", "server host")
	checkNoErr(viper.BindPFlag("host", flags.Lookup("host")))

	flags.IntP("port", "p", 8080, "server port")
	checkNoErr(viper.BindPFlag("port", flags.Lookup("port")))

	flags.String("admin-host", "localhost", "administration server host")
	checkNoErr(viper.BindPFlag("admin.host", flags.Lookup("admin-host")))

	flags.Int("admin-port", 6060, "administration server port")
	checkNoErr(viper.BindPFlag("admin.port", flags.Lookup("admin-port")))

	flags.String("log-level", "info", "define the log level")
	checkNoErr(viper.BindPFlag("log.level", flags.Lookup("log-level")))
}

// Configure Viper to read the environment and the optional config file
func Configure() error {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("cozy")
	viper.AutomaticEnv()

	if cfgFile != "" {
		// Read given config file and skip other paths
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(config.Filename)
		for _, cfgPath := range config.Paths {
			viper.AddConfigPath(cfgPath)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, isParseErr := err.(viper.ConfigParseError); isParseErr {
			log.Errorf("Failed to read cozy-stack configurations from %s", viper.ConfigFileUsed())
			return err
		}

		if cfgFile != "" {
			return fmt.Errorf("Could not locate config file: %s", cfgFile)
		}
	}

	if viper.ConfigFileUsed() != "" {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	if err := config.UseViper(viper.GetViper()); err != nil {
		return err
	}

	return nil
}

func checkNoErr(err error) {
	if err != nil {
		panic(err)
	}
}
