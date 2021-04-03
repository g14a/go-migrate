/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this gen except in compliance with the License.
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
	"github.com/fatih/color"
	"github.com/g14a/go-migrate/pkg/gen"
	"log"
	"os"
	"os/exec"
	//"path/filepath"

	//"github.com/g14a/go-migrate/pkg/gen"
	"github.com/itchyny/gojq"

	//"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a migrations directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.MkdirAll("migrations/scripts", 0755)
		wd, _ := os.Getwd()

		goModInfo, err := exec.Command("go", "mod", "edit", "-json").Output()

		query, err := gojq.Parse(".Module.Path | ..")
		if err != nil {
			log.Fatal(err)
		}

		goModDetails := make(map[string]interface{})

		errJson := json.Unmarshal(goModInfo, &goModDetails)
		if errJson != nil {
			log.Fatal(errJson)
		}

		iter := query.Run(goModDetails)

		var goModPath string
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				log.Fatal(err)
			}
			goModPath = v.(string)
		}

		err = gen.CreateInitConfig(goModPath)
		if err != nil {
			log.Fatal(err)
		}

		color.Green(" ✓ Created " + wd + "/migrations/main.go")
		color.Green(" ✓ Created " + wd + "/migrations/store.go")
		color.Green(" ✓ Created " + wd + "/migrations/migrate.json")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
