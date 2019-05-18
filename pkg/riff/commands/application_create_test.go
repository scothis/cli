/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package commands_test

import (
	"github.com/projectriff/riff/pkg/cli"
	"github.com/projectriff/riff/pkg/riff/commands"
	"github.com/projectriff/riff/pkg/testing"
	buildv1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestApplicationCreateOptions(t *testing.T) {
	table := testing.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.InvalidResourceOptions,
			},
			ExpectFieldError: testing.InvalidResourceOptionsFieldError.Also(
				cli.ErrMissingField(cli.ImageFlagName),
				cli.ErrMissingOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
			),
		},
		{
			Name: "git source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
			},
			ShouldValidate: true,
		},
		{
			Name: "local source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
			},
			ShouldValidate: true,
		},
		{
			Name: "no source",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
			},
			ExpectFieldError: cli.ErrMissingOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
		},
		{
			Name: "multiple sources",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				LocalPath:       ".",
			},
			ExpectFieldError: cli.ErrMultipleOneOf(cli.GitRepoFlagName, cli.LocalPathFlagName),
		},
		{
			Name: "git source with cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				CacheSize:       "8Gi",
			},
			ShouldValidate: true,
		},
		{
			Name: "local source with cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
				CacheSize:       "8Gi",
			},
			ExpectFieldError: cli.ErrDisallowedFields(cli.CacheSizeFlagName),
		},
		{
			Name: "invalid cache",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				CacheSize:       "X",
			},
			ExpectFieldError: cli.ErrInvalidValue("X", cli.CacheSizeFlagName),
		},
		{
			Name: "with git subpath",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "master",
				SubPath:         "some/directory",
			},
			ShouldValidate: true,
		},
		{
			Name: "with local subpath",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				LocalPath:       ".",
				SubPath:         "some/directory",
			},
			ExpectFieldError: cli.ErrDisallowedFields(cli.SubPathFlagName),
		},
		{
			Name: "missing git revision",
			Options: &commands.ApplicationCreateOptions{
				ResourceOptions: testing.ValidResourceOptions,
				Image:           "example.com/repo:tag",
				GitRepo:         "https://example.com/repo.git",
				GitRevision:     "",
			},
			ExpectFieldError: cli.ErrMissingField(cli.GitRevisionFlagName),
		},
	}

	table.Run(t)
}

func TestApplicationCreateCommand(t *testing.T) {
	defaultNamespace := "default"
	applicationName := "my-application"
	imageTag := "registry.example.com/repo:tag"
	imageDigest := "registry.example.com/repo@sha256:deadbeefdeadbeefdeadbeefdeadbeef"
	gitRepo := "https://example.com/repo.git"
	gitMaster := "master"
	gitSha := "deadbeefdeadbeefdeadbeefdeadbeef"
	subPath := "some/path"
	cacheSize := "8Gi"
	cacheSizeQuantity := resource.MustParse(cacheSize)
	localPath := "."

	table := testing.CommandTable{
		{
			Name:        "invalid args",
			Args:        []string{},
			ShouldError: true,
		},
		{
			Name: "git repo",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with revision",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.GitRevisionFlagName, gitSha},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitSha,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with subpath",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.SubPathFlagName, subPath},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitMaster,
							},
							SubPath: subPath,
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			Name: "git repo with cache",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo, cli.CacheSizeFlagName, cacheSize},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image:     imageTag,
						CacheSize: &cacheSizeQuantity,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ExpectOutput: `
Created application "my-application"
`,
		},
		{
			// TODO impelement
			Skip: true,
			Name: "local path",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.LocalPathFlagName, localPath},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
					},
					Status: buildv1alpha1.ApplicationStatus{
						BuildStatus: buildv1alpha1.BuildStatus{
							LatestImage: imageDigest,
						},
					},
				},
			},
		},
		{
			Name: "error existing application",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			GivenObjects: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
				},
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
		{
			Name: "error during create",
			Args: []string{applicationName, cli.ImageFlagName, imageTag, cli.GitRepoFlagName, gitRepo},
			WithReactors: []testing.ReactionFunc{
				testing.InduceFailure("create", "applications"),
			},
			ExpectCreates: []runtime.Object{
				&buildv1alpha1.Application{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: defaultNamespace,
						Name:      applicationName,
					},
					Spec: buildv1alpha1.ApplicationSpec{
						Image: imageTag,
						Source: buildv1alpha1.Source{
							Git: &buildv1alpha1.GitSource{
								URL:      gitRepo,
								Revision: gitMaster,
							},
						},
					},
				},
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewApplicationCreateCommand)
}