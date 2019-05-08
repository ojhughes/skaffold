/*
Copyright 2019 The Skaffold Authors

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

package buildpacks

import (
	"context"
	"sort"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/docker"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/pkg/errors"
)

// GetDependencies returns dependencies listed for a buildpacks artifact
func GetDependencies(ctx context.Context, workspace string, a *latest.BuildpacksArtifact) ([]string, error) {
	path := make([]string, 1)
	if a.AppPath != "" {
		path = append(path, a.AppPath)
	} else {
		path = append(path, ".")
	}
	files, err := docker.WalkWorkspace(workspace, a.Ignore, path)
	if err != nil {
		return nil, errors.Wrapf(err, "walking workspace %s", workspace)
	}
	var dependencies []string
	for file := range files {
		dependencies = append(dependencies, file)
	}
	sort.Strings(dependencies)
	return dependencies, nil
}
