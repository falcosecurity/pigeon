package poiana

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
