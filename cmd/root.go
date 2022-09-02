/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"log"
	"time"

	"os"

	"log_extraction/tool"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "batteryLogTool",
	Short: "Battery log extraction",
	Long:  `Battery log extraction`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()

		f, _ := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		defer f.Close()
		log.SetOutput(f)

		if mongoStr != "" {
			saveToMongo = true
		}
		if outFile != "" {
			saveToFile = true
		}

		file, err := os.Open(inFile)
		if err != nil {
			log.Printf("Can not open file: %s, err: [%v]", inFile, err)
		}
		defer file.Close()

		client, err := tool.NewMongoDbPool(mongoStr)
		if err != nil {
			log.Printf("Connect to mongodb err: [%v]", err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			text, is := tool.ReText(line)
			if is {
				if saveToFile {
					tool.WriteFile(outFile, text)
				}
				if saveToMongo {
					jsonStr := tool.ReJsonText(text)
					client.InsertToDb(jsonStr)
				}

			}
		}
		slapse := time.Since(start)
		log.Println("elapsed time is : ", slapse)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var inFile string
var outFile string
var logFile string
var mongoStr string
var saveToFile bool
var saveToMongo bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.log_extraction.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&inFile, "in", "i", "", "file input")
	rootCmd.Flags().StringVarP(&outFile, "out", "o", "", "file out")
	rootCmd.Flags().StringVarP(&logFile, "log", "l", "extraction.log", "log file")
	rootCmd.Flags().StringVarP(&mongoStr, "mongostr", "m", "", `mongodb connect str, scheme must be "mongodb" or "mongodb+srv"`)
	rootCmd.MarkFlagRequired("in")

}
