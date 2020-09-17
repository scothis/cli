/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package k8s_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/projectriff/cli/pkg/k8s"
	rifftesting "github.com/projectriff/cli/pkg/testing"
	streamingv1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	"github.com/vmware-labs/reconciler-runtime/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	cachetesting "k8s.io/client-go/tools/cache/testing"
)

func TestWaitUntilReady(t *testing.T) {
	processor := &streamingv1alpha1.Processor{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Processor",
			APIVersion: "streaming.projectriff.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "my-processor",
			UID:       "c6acbbab-87dd-11e9-807c-42010a80011d",
		},
		Status: streamingv1alpha1.ProcessorStatus{
			Status: apis.Status{
				Conditions: apis.Conditions{
					{
						Type:   apis.ConditionReady,
						Status: corev1.ConditionUnknown,
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		resource *streamingv1alpha1.Processor
		events   []watch.Event
		err      error
	}{{
		name:     "transitions true",
		resource: processor.DeepCopy(),
		events: []watch.Event{
			updateReady(processor, corev1.ConditionTrue, ""),
		},
	}, {
		name:     "transitions false",
		resource: processor.DeepCopy(),
		events: []watch.Event{
			updateReady(processor, corev1.ConditionFalse, "test not ready"),
		},
		err: fmt.Errorf("failed to become ready: %s", "test not ready"),
	}, {
		name:     "ignore other resources",
		resource: processor.DeepCopy(),
		events: []watch.Event{
			updateReadyOther(processor, corev1.ConditionFalse, "not my app"),
			updateReady(processor, corev1.ConditionTrue, ""),
		},
	}, {
		name:     "bail on delete",
		resource: processor.DeepCopy(),
		events: []watch.Event{
			updateReady(processor, corev1.ConditionUnknown, ""),
			watch.Event{Type: watch.Deleted, Object: processor.DeepCopy()},
		},
		err: fmt.Errorf("%s %q deleted", "processor", "my-processor"),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lw := cachetesting.NewFakeControllerSource()
			defer lw.Shutdown()
			ctx := k8s.WithListerWatcher(context.Background(), lw)

			client := rifftesting.NewClient(processor)
			done := make(chan error, 1)
			defer close(done)
			go func() {
				done <- k8s.WaitUntilReady(ctx, client.StreamingRuntime().RESTClient(), "processors", processor)
			}()

			time.Sleep(5 * time.Millisecond)
			for _, event := range test.events {
				lw.Change(event, 1)
			}

			err := <-done
			if expected, actual := fmt.Sprintf("%s", test.err), fmt.Sprintf("%s", err); expected != actual {
				t.Errorf("expected error %v, actually %v", expected, actual)
			}
		})
	}
}

func updateReady(processor *streamingv1alpha1.Processor, status corev1.ConditionStatus, message string) watch.Event {
	processor = processor.DeepCopy()
	processor.Status.Conditions[0].Status = status
	processor.Status.Conditions[0].Message = message
	return watch.Event{Type: watch.Modified, Object: processor}
}

func updateReadyOther(processor *streamingv1alpha1.Processor, status corev1.ConditionStatus, message string) watch.Event {
	processor = processor.DeepCopy()
	processor.UID = "not-a-uid"
	processor.Status.Conditions[0].Status = status
	processor.Status.Conditions[0].Message = message
	return watch.Event{Type: watch.Modified, Object: processor}
}
