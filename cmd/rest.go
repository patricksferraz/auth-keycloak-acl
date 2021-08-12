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
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/c-4u/auth-service/application/rest"
	"github.com/c-4u/auth-service/infrastructure/external"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// NewRestCmd represents the rest command
func NewRestCmd() *cobra.Command {
	var restPort int

	restCmd := &cobra.Command{
		Use:   "rest",
		Short: "Run rest Service",

		Run: func(cmd *cobra.Command, args []string) {
			service := external.NewKeycloak(
				os.Getenv("KEYCLOAK_BASE_PATH"),
				os.Getenv("KEYCLOAK_REALM"),
				os.Getenv("KEYCLOAK_CLIENT_ID"),
				os.Getenv("KEYCLOAK_CLIENT_SECRET"),
				os.Getenv("KEYCLOAK_AUDIENCE"),
			)

			deliveryChan := make(chan ckafka.Event)
			kafka, err := external.NewKafka(
				os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
				deliveryChan,
			)
			if err != nil {
				log.Fatal(err)
			}

			employeeServiceAddr := os.Getenv("EMPLOYEE_SERVICE_ADDR")

			go kafka.DeliveryReport()
			rest.StartRestServer(service, kafka, employeeServiceAddr, restPort)
		},
	}

	restCmd.Flags().IntVarP(&restPort, "port", "p", 8080, "rest server port")

	return restCmd
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	if os.Getenv("ENV") == "dev" {
		err := godotenv.Load(basepath + "/../.env")
		if err != nil {
			log.Printf("Error loading .env files")
		}
	}

	rootCmd.AddCommand(NewRestCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
