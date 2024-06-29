/*
Copyright © 2024 Julien Noblet
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-clean-dat",
	Short: "A tool to clean a mame dat file",
	Long: `A tool to clean a mame dat file.
	This tool can remove clones, mechanical, devices, pinball, systems, preliminary, imperfect and sourcefile from a mame dat file.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func binder(name string) {
	err := viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
	if err != nil {
		panic("can't bind " + name + " flag")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".go-clean-dat.yaml", "config file")
	rootCmd.PersistentFlags().StringP("input", "i", "", "input file (required)")
	rootCmd.PersistentFlags().StringP("output", "o", "out.dat", "output file")
	rootCmd.PersistentFlags().StringSlice("others", []string{}, "others input dat file")
	rootCmd.PersistentFlags().Bool("no-clones", false, "no clones")
	rootCmd.PersistentFlags().Bool("no-mechanical", false, "no mechanical")
	rootCmd.PersistentFlags().Bool("no-devices", false, "no devices")
	rootCmd.PersistentFlags().Bool("no-pinball", false, "no pinball")
	rootCmd.PersistentFlags().StringSlice("no-system", []string{}, "no system")
	rootCmd.PersistentFlags().Bool("no-preliminary", false, "no preliminary")
	rootCmd.PersistentFlags().Bool("no-imperfect", false, "no imperfect")
	rootCmd.PersistentFlags().StringSlice("no-sourcefile", []string{}, "no sourcefile")

	// bind to viper Pflags
	binder("input")
	binder("output")
	binder("others")
	binder("no-clones")
	binder("no-mechanical")
	binder("no-devices")
	binder("no-pinball")
	binder("no-system")
	binder("no-preliminary")
	binder("no-imperfect")
	binder("no-sourcefile")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current directory
		pwd := os.Getenv("PWD")

		// Search config in home directory with name ".go-clean-dat" (without extension).
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-clean-dat")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
