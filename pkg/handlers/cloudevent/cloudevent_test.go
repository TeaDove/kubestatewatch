/*
Copyright 2018 Bitnami

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

package cloudevent

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/marvasgit/kubestatewatch/config"
)

func TestCloudEventInit(t *testing.T) {
	s := &CloudEvent{}
	expectedError := fmt.Errorf(cloudEventErrMsg, "Missing cloudevent url")

	var Tests = []struct {
		cloudevent config.CloudEvent
		err        error
	}{
		{config.CloudEvent{Url: "foo"}, nil},
		{config.CloudEvent{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.CloudEvent = tt.cloudevent
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}
