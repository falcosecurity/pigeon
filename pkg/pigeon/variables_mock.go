// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
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

package pigeon

import (
	"context"

	"github.com/google/go-github/v50/github"
)

type MockVariableService struct {
	variables map[string]string
}

func (m MockVariableService) ListVariables(context.Context, *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
	vars := make([]*github.ActionsVariable, 0)
	for key, val := range m.variables {
		vars = append(vars, &github.ActionsVariable{
			Name:  key,
			Value: val,
		})
	}

	return &github.ActionsVariables{
		TotalCount: len(m.variables),
		Variables:  vars,
	}, nil, nil
}

func (m MockVariableService) DeleteVariable(_ context.Context, name string) (*github.Response, error) {
	delete(m.variables, name)
	return nil, nil
}

func (m MockVariableService) CreateOrUpdateVariable(_ context.Context, variable *github.ActionsVariable) (*github.Response, error) {
	m.variables[variable.Name] = variable.Value
	return nil, nil
}

func NewMockVariableService() ActionsVarsService {
	return &MockVariableService{
		variables: make(map[string]string, 0),
	}
}
