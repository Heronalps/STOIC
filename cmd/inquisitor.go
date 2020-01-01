/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"time"

	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// inquisitorCmd represents the inquisitor command
var (
	interval      int
	inquisitorCmd = &cobra.Command{
		Use:   "inquisitor",
		Short: "Inquisitor keeps probing Nautilus for deployment time of runtimes",
		Long:  `Inquisitor keeps probing Nautilus for deployment tiem of runtimes`,
		Run: func(cmd *cobra.Command, args []string) {
			for {
				client.UpdateDeploymentTimeTable()
				fmt.Println("Waiting for next round ...")
				time.Sleep(time.Second * time.Duration(interval))
			}
		},
	}
)

func init() {
	runCmd.AddCommand(inquisitorCmd)
	inquisitorCmd.Flags().IntVarP(&interval, "interval", "i", 600, "The interval of inquire deployment time on Nautilus")
}
