/*
Copyright IBM Corporation 2020

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

package main

import (
	"path/filepath"

	"github.com/konveyor/move2kube/internal/move2kube"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	artifactsPath = "artifacts"
)

type validateFlags struct {
	artifactspath string
}

func validateHandler(flags validateFlags) {
	artifactspath, err := filepath.Abs(flags.artifactspath)
	if err != nil {
		log.Fatalf("Failed to make the directory path %q absolute. Error: %q", artifactspath, err)
	}
	move2kube.PrintValidate(artifactspath)
}

func getValidateCommand() *cobra.Command {
	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	viper.AutomaticEnv()

	flags := validateFlags{}
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Prints all next steps in generated artifacts",
		Long:  "Next step actions are distributed among the artifacts generated by Move2Kube. This command aggregates next steps from the artifacts and gives a comprehensive view.",
		Run:   func(*cobra.Command, []string) { validateHandler(flags) },
	}

	validateCmd.Flags().StringVarP(&flags.artifactspath, artifactsPath, "a", ".", "Specify directory containing the artifacts generated by Move2Kube.")

	must(validateCmd.MarkFlagRequired(artifactsPath))

	return validateCmd
}
