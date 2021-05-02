// Package cmd /*
package cmd

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Run the downward migration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		downUntil, _ := cmd.Flags().GetString("until")
		migrationsBuild := exec.Command("go", "build")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		migrationsBuild.Dir = wd + "/migrations"

		errBuild := migrationsBuild.Start()
		if errBuild != nil {
			log.Fatal(errBuild)
		}

		errWait := migrationsBuild.Wait()
		if errWait != nil {
			log.Fatal(errWait)
		}

		migrationArgs := []string{"down"}
		if downUntil != "" {
			migrationArgs = append(migrationArgs, "--until", downUntil)
		}

		migrationsRun := exec.Command("./migrations", migrationArgs...)
		migrationsRun.Dir = wd + "/migrations"
		var outBuf, errBuf bytes.Buffer
		migrationsRun.Stdout = &outBuf
		migrationsRun.Stderr = &errBuf

		errRun := migrationsRun.Run()
		if errRun != nil {
			log.Fatal(errRun)
		}

		if errBuf.Len() > 0 {
			log.Fatal(errBuf.String())
		}

		if outBuf.Len() > 0 {
			fmt.Printf("%v\n", outBuf.String())
		}

		color.Green("  >>> migration : complete")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().StringP("until", "u", "", "Migrate down until a specific point\n")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
