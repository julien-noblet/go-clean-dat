/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/xml"
	"log"
	"os"

	mamedattools "github.com/julien-noblet/go-clean-dat/mameDatTools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean a mame dat file",
	Long: `Clean a mame dat file.
	Can use other dat files to remove games from the input file.
	Can remove clones, mechanical, devices, pinball, systems, preliminary, imperfect and sourcefile from a mame dat file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("input") == "" {
			err := rootCmd.Help()
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal("Error: input flag is required!")
		}
		mame := mamedattools.OpenAndUnmarshal(viper.GetString("input"))
		startGames := len(mame.Games)
		mame = mamedattools.CleanMame(&mame)
		for _, other := range viper.GetStringSlice("others") {
			mameB := mamedattools.OpenAndUnmarshal(other)
			mameB = mamedattools.CleanMame(&mameB)
			log.Println("Comparing with", other)
			mame = mamedattools.NotInB(&mame, &mameB, other)
		}
		log.Printf("file %s has %d games, %d games removed", viper.GetString("input"), startGames, startGames-len(mame.Games))
		// Ouvre le fichier new.dat et ecrit la struct Mame dedans
		newFile, err := os.Create(viper.GetString("output"))
		if err != nil {
			log.Fatal(err)
		}
		defer newFile.Close()

		// ajoute le header
		_, err = newFile.Write([]byte(xml.Header))
		if err != nil {
			log.Fatal(err)
		}
		// Marshal la struct Mame dans le fichier xml
		enc, err := xml.MarshalIndent(mame, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		_, err = newFile.Write(enc)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
