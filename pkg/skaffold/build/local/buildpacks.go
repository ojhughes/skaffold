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

package local

import (
	"context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build/buildpacks"
	//"github.com/buildpack/lifecycle/image"
	//"github.com/buildpack/pack"
	//"github.com/buildpack/pack/config"
	//"github.com/buildpack/pack/docker"
	"io"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/pkg/errors"
)

func (b *Builder) buildBuildpacks(ctx context.Context, out io.Writer, artifact *latest.BuildpacksArtifact, tag string) (string, error) {
	cfg, err := config.NewDefault()
	if err != nil {
		return "", errors.Wrap(err, "initializing pack client")
	}
	factory, err := image.NewFactory()
	if err != nil {
		return "", errors.Wrap(err, "initializing pack client")
	}

	dockerClient, err := docker.New()
	if err != nil {
		return "", errors.Wrap(err, "initializing pack client")
	}
	client := *pack.NewClient(&cfg, &pack.ImageFetcher{
		Factory: factory,
		Docker:  dockerClient,
	})

	customArtifactBuilder := buildpacks.NewArtifactBuilder(b.pushImages, b.localDocker.ExtraEnv())

	if err := customArtifactBuilder.Build(ctx, out, artifact, tag); err != nil {
		return "", errors.Wrap(err, "building custom artifact")
	}

	if b.pushImages {
		return docker.RemoteDigest(tag, b.insecureRegistries)
	}

	return b.localDocker.ImageID(ctx, tag)
}
