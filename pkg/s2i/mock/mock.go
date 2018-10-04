// Copyright 2018 Caicloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mock

import "github.com/caicloud/ciao/pkg/types"

// Mocker is the type for S2I mocker.
type Mocker struct{}

// New returns a new Mocker.
func New() *Mocker {
	return &Mocker{}
}

// SourceToImage uses the default image without building one.
func (m Mocker) SourceToImage(code string, parameter *types.Parameter) (string, error) {
	return "kubeflow/tf-dist-mnist-test:1.0", nil
}
