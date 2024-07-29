/*
Copyright © 2024 Nicholas Hadley <contact@nicholashadley.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/hadley31/cs2pm/util"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs plugins from a registry file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dest := cmd.Flag("dir").Value.String()

		plugins, err := util.ReadYamlFile("cs2pm.yaml")

		if err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}

		for _, config := range plugins.Plugins {
			wg.Add(1)
			installPlugin(&config, dest, wg)
		}

		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("dir", "d", "", "Directory to install the plugin to")
}

func installPlugin(config *util.PluginConfig, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Installing plugin %s\n", config.Name)

	extractDir := filepath.Join(dest, config.ExtractPrefix)

	tempFile, err := util.DownloadPlugin(config.DownloadUrl)

	if err != nil {
		panic(err)
	}

	util.UnzipPlugin(tempFile, extractDir)
}
