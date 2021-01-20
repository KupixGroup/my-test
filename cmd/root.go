/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Buys struct {
	Product     string
	Productsena int
	User        int
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-test",
	Short: "A brief description of your application",
	Long:  `Hello`,
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
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(buysCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var pushCmd = &cobra.Command{
	Use:   "send",
	Short: "Example: 5(userid) pilus24@mail.ru 10(productid) Samsung 100000(sena)",
	Run: func(cmd *cobra.Command, args []string) {

		iduser, _ := strconv.Atoi(args[0])
		tou := args[1]
		idproduct, _ := strconv.Atoi(args[2])
		nameproduct := args[3]
		senaproduct, _ := strconv.Atoi(args[4])

		if iduser > 0 && idproduct > 0 {
			if isEmailValid(tou) {
				//для электронной почты
				/*
					from := "from@gmail.com"
					password := "<Email Password>"

					to := []string{
						tou,
					}

					// smtp server configuration.
					smtpHost := "smtp.gmail.com"
					smtpPort := "587"

					// Message.
					message := []byte("Ваша покупка была успешно завершена.")

					// Authentication.
					auth := smtp.PlainAuth("", from, password, smtpHost)

					// Sending email.
					err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
					if err != nil {
						fmt.Println(err)
						return
					}
				*/
				fmt.Println("Email sent to " + tou)

			} else {
				// Yoki sms api
				//Или смс api
				fmt.Println("SMS sent phone to " + tou)
			}

			file, err := ioutil.ReadFile("buy.json")
			if err != nil {
				fmt.Printf("error open file")
			}

			data := []Buys{}

			// Here the magic happens!
			json.Unmarshal(file, &data)

			newStruct := &Buys{
				User:        iduser,
				Product:     nameproduct,
				Productsena: senaproduct,
			}

			data = append(data, *newStruct)

			// Preparing the data to be marshalled and written.
			dataBytes, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("error json")
			}

			err = ioutil.WriteFile("buy.json", dataBytes, 0644)
			if err != nil {
				fmt.Printf("error save cache")
			}

		} else {
			fmt.Printf("error id")
		}

	},
}

var buysCmd = &cobra.Command{
	Use:   "buys",
	Short: "All buying list",
	Run: func(cmd *cobra.Command, args []string) {

		file, err1 := ioutil.ReadFile("buy.json")
		if err1 != nil {
			fmt.Printf("// error while reading file %s\n")
			fmt.Printf("File error: %v\n", err1)
			os.Exit(1)
		}
		var rep []Buys
		err2 := json.Unmarshal(file, &rep)
		if err2 != nil {
			fmt.Println("error:", err2)
		}
		fmt.Printf("Buying\n")
		for k := range rep {
			fmt.Printf("Producr: %v, Product price:  %v, Userid: %v\n", rep[k].Product, rep[k].Productsena, rep[k].User)
		}

	},
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

		// Search config in home directory with name ".my-test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".my-test")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
