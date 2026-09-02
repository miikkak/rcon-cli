package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/miikkak/rcon-cli/internal/rcon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "rcon-cli [flags] [RCON command ...]",
	Short: "A CLI for attaching to an RCON enabled game server",
	Example: `
rcon-cli --host mc1 --port 25575
rcon-cli --port 25575 stop
RCON_PORT=25575 rcon-cli stop
`,
	Long: `
rcon-cli is a CLI for attaching to an RCON enabled game server, such as Minecraft.
Without any additional arguments, the CLI will start an interactive session with
the RCON server. Send the keyword "exit" to close the session. The CLI will exit.

If arguments are passed into the CLI, then the arguments are sent
as a single command (joined by spaces), the response is displayed,
and the CLI will exit.
`,
	RunE: func(_ *cobra.Command, args []string) error {
		hostPort := net.JoinHostPort(viper.GetString("host"), strconv.Itoa(viper.GetInt("port")))
		password := viper.GetString("password")

		if len(args) == 0 {
			return rcon.Start(hostPort, password, os.Stdout, os.Stderr)
		}
		return rcon.Execute(hostPort, password, os.Stdout, args...)
	},
}

// Execute adds all child commands to the root command and executes it. This
// is called by main.main(). It only needs to happen once to RootCmd.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rcon-cli.yaml)")
	RootCmd.PersistentFlags().String("host", "localhost", "RCON server's hostname")
	RootCmd.PersistentFlags().String("password", "", "RCON server's password")
	RootCmd.PersistentFlags().Int("port", 25575, "Server's RCON port")
	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".rcon-cli") // name of config file (without extension)
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(home)
		}
		viper.AddConfigPath("/data")
		viper.AddConfigPath("/server")
	}

	// This will allow for env vars like RCON_PORT
	viper.SetEnvPrefix("rcon")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in. A missing config file is fine
	// when we're only searching default locations, but any other error
	// (malformed YAML, an explicit --config pointing at a bad path) should
	// stop the CLI rather than silently falling back to defaults.
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if cfgFile == "" && errors.As(err, &notFound) {
			return
		}
		fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		os.Exit(1)
	}
}
