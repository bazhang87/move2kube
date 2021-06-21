/*
Copyright IBM Corporation 2021

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

package generators

import (
	"path/filepath"

	"github.com/konveyor/move2kube/environment"
	"github.com/konveyor/move2kube/internal/apiresource"
	"github.com/konveyor/move2kube/internal/common"
	irtypes "github.com/konveyor/move2kube/types/ir"
	plantypes "github.com/konveyor/move2kube/types/plan"
	transformertypes "github.com/konveyor/move2kube/types/transformer"
	"github.com/sirupsen/logrus"
)

const (
	KnativeArtifacts = "KnativeYamls"
)

// Knative implements Transformer interface
type Knative struct {
	Config transformertypes.Transformer
	Env    environment.Environment
}

type KnativeConfig struct {
}

func (t *Knative) Init(tc transformertypes.Transformer, env environment.Environment) error {
	t.Config = tc
	t.Env = env
	return nil
}

func (t *Knative) GetConfig() (transformertypes.Transformer, environment.Environment) {
	return t.Config, t.Env
}

func (t *Knative) BaseDirectoryDetect(dir string) (namedServices map[string]plantypes.Service, unnamedServices []plantypes.Transformer, err error) {
	return nil, nil, nil
}

func (t *Knative) DirectoryDetect(dir string) (namedServices map[string]plantypes.Service, unnamedServices []plantypes.Transformer, err error) {
	return nil, nil, nil
}

func (t *Knative) Transform(newArtifacts []transformertypes.Artifact, oldArtifacts []transformertypes.Artifact) (pathMappings []transformertypes.PathMapping, createdArtifacts []transformertypes.Artifact, err error) {
	logrus.Debugf("Translating IR using Kubernetes transformer")
	pathMappings = []transformertypes.PathMapping{}
	for _, a := range newArtifacts {
		if a.Artifact != transformertypes.IRArtifactType {
			continue
		}
		var ir irtypes.IR
		err := a.GetConfig(transformertypes.IRConfigType, &ir)
		if err != nil {
			logrus.Errorf("unable to load config for Transformer into %T : %s", ir, err)
			continue
		}
		var pC transformertypes.PlanConfig
		err = a.GetConfig(transformertypes.PlanConfigType, &pC)
		if err != nil {
			logrus.Errorf("unable to load config for Transformer into %T : %s", pC, err)
			continue
		}
		tempDest := filepath.Join(t.Env.TempPath, common.DeployDir, "knative")
		logrus.Debugf("Starting Kubernetes transform")
		logrus.Debugf("Total services to be transformed : %d", len(ir.Services))
		apis := []apiresource.IAPIResource{&apiresource.KnativeService{}}
		if files, err := apiresource.TransformAndPersist(irtypes.NewEnhancedIRFromIR(ir), tempDest, apis, pC.TargetCluster); err == nil {
			for _, f := range files {
				if destPath, err := filepath.Rel(t.Env.TempPath, f); err != nil {
					logrus.Errorf("Invalid yaml path : %s", destPath)
				} else {
					pathMappings = append(pathMappings, transformertypes.PathMapping{
						Type:     transformertypes.DefaultPathMappingType,
						SrcPath:  f,
						DestPath: destPath,
					})
				}
			}
			logrus.Debugf("Total transformed objects : %d", len(files))
		} else {
			logrus.Errorf("Unable to transform and persist IR : %s", err)
			return nil, nil, err
		}
	}
	return pathMappings, nil, nil
}