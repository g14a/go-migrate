package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func RunSetConfig(cmd *cobra.Command, FS afero.Fs, wd string) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	store, err := cmd.Flags().GetString("store")
	if err != nil {
		return err
	}

	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return err
	}

	mc, err := config.GetMetanaConfig(FS, wd)
	if os.IsNotExist(err) {
		_, err = os.Create(".metana.yml")
		if err != nil {
			return err
		}
	}

	if env != "" {
		if len(mc.Environments) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), color.YellowString("No environment configured yet.\nTry initializing one with `metana init --env "+env+"`\n"))
			return nil
		}
		err = config.SetEnvironmentMetanaConfig(mc, env, store, FS, wd)
		if err != nil {
			return err
		}
		return nil
	}

	if dir != "" {
		mc.Dir = dir
	}

	if store != "" {
		mc.StoreConn = store
	}

	err = config.SetMetanaConfig(mc, FS, wd)
	if err != nil {
		return err
	}

	return nil
}
