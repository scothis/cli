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
	"context"
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/streaming/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	"github.com/projectriff/system/pkg/apis"
	streamv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	"github.com/projectriff/system/pkg/refs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPulsarProviderListOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid list",
			Options: &commands.PulsarProviderListOptions{
				ListOptions: rifftesting.InvalidListOptions,
			},
			ExpectFieldErrors: rifftesting.InvalidListOptionsFieldError,
		},
		{
			Name: "valid list",
			Options: &commands.PulsarProviderListOptions{
				ListOptions: rifftesting.ValidListOptions,
			},
			ShouldValidate: true,
		},
	}

	table.Run(t)
}

func TestPulsarProviderListCommand(t *testing.T) {
	pulsarProviderName := "test-pulsar-provider"
	pulsarProviderOtherName := "test-other-pulsar-provider"
	defaultNamespace := "default"
	otherNamespace := "other-namespace"

	table := rifftesting.CommandTable{
		{
			Name: "invalid args",
			Args: []string{},
			Prepare: func(t *testing.T, ctx context.Context, c *cli.Config) (context.Context, error) {
				// disable default namespace
				c.Client.(*rifftesting.FakeClient).Namespace = ""
				return ctx, nil
			},
			ShouldError: true,
		},
		{
			Name: "empty",
			Args: []string{},
			ExpectOutput: `
No pulsar providers found.
`,
		},
		{
			Name: "lists an item",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
NAME                   SERVICE URL   PROVISIONER   STATUS      AGE
test-pulsar-provider   <empty>       <empty>       <unknown>   <unknown>
`,
		},
		{
			Name: "filters by namespace",
			Args: []string{cli.NamespaceFlagName, otherNamespace},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
			},
			ExpectOutput: `
No pulsar providers found.
`,
		},
		{
			Name: "all namespace",
			Args: []string{cli.AllNamespacesFlagName},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderName,
						Namespace: defaultNamespace,
					},
				},
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pulsarProviderOtherName,
						Namespace: otherNamespace,
					},
				},
			},
			ExpectOutput: `
NAMESPACE         NAME                         SERVICE URL   PROVISIONER   STATUS      AGE
default           test-pulsar-provider         <empty>       <empty>       <unknown>   <unknown>
other-namespace   test-other-pulsar-provider   <empty>       <empty>       <unknown>   <unknown>
`,
		},
		{
			Name: "table populates all columns",
			Args: []string{},
			GivenObjects: []runtime.Object{
				&streamv1alpha1.PulsarProvider{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-pulsar",
						Namespace: defaultNamespace,
					},
					Spec: streamv1alpha1.PulsarProviderSpec{
						ServiceURL: "pulsar://localhost:6650",
					},
					Status: streamv1alpha1.PulsarProviderStatus{
						Status: apis.Status{
							Conditions: apis.Conditions{
								{Type: streamv1alpha1.PulsarProviderConditionReady, Status: "True"},
							},
						},
						ProvisionerServiceRef: &refs.TypedLocalObjectReference{Name: "my-pulsar-provisioner", Kind: "Service"},
					},
				},
			},
			ExpectOutput: `
NAME        SERVICE URL               PROVISIONER             STATUS   AGE
my-pulsar   pulsar://localhost:6650   my-pulsar-provisioner   Ready    <unknown>
`,
		},
		{
			Name: "list error",
			Args: []string{},
			WithReactors: []rifftesting.ReactionFunc{
				rifftesting.InduceFailure("list", "pulsarproviders"),
			},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewPulsarProviderListCommand)
}
