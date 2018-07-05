package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	arrowRight = "➞"
	printFlags = true
)

var (
	cfgFile     string
	prettyPrint = true
)
var rootCmd = &cobra.Command{
	Use:   "iredmail-cli",
	Short: "A command line inteface to manage iRedMail server",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.iredmail-cli.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&prettyPrint, "pretty", "", true, "")
	rootCmd.PersistentFlags().MarkHidden("pretty")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".iredmail-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".iredmail-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("No config file found (default ~/.iredmail-cli.yml")
	}
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fatal("%v\n", err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func success(format string, a ...interface{}) {
	c := color.New(color.FgGreen, color.Bold)
	c.Printf(format, a...)
}

func warning(format string, a ...interface{}) {
	c := color.New(color.FgRed, color.Bold)
	c.Printf(format, a...)
}

func info(format string, a ...interface{}) {
	c := color.New()
	c.Printf(format, a...)
}

func fatal(format string, a ...interface{}) {
	c := color.New(color.FgRed, color.Bold)
	c.Printf(format, a...)
	os.Exit(1)
}

func usageTemplate(useLine string, args ...bool) string {
	var printFlags bool
	var flags string

	if len(args) > 0 {
		printFlags = args[0]
	}
	if printFlags {
		flags = " [flags]"
	}

	return `Usage:{{if .Runnable}}
    iredmail-cli ` + useLine + flags + `{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
  
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}
  
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
  
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
