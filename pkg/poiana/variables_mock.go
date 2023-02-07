package poiana

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"net/http"
)

type MockVariableService struct {
	variables map[string]string
}

func (m MockVariableService) ListRepoVariables(_ context.Context, _, _ string, _ *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
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

func (m MockVariableService) DeleteRepoVariable(_ context.Context, _, _, name string) (*github.Response, error) {
	delete(m.variables, name)
	return nil, nil
}

func (m MockVariableService) CreateRepoVariable(_ context.Context, _, _ string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.variables[variable.Name]; ok {
		// Key is present; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
		}, fmt.Errorf("variable already present")
	}
	// Key not present, add value
	m.variables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (m MockVariableService) UpdateRepoVariable(_ context.Context, _, _ string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.variables[variable.Name]; !ok {
		// There is no key; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
		}, fmt.Errorf("variable not present")
	}
	// Key is present, set value
	m.variables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}

func NewMockVariableService() ActionsVarsService {
	mServ := &MockVariableService{variables: make(map[string]string, 0)}
	// Initial variable set
	_, _ = mServ.CreateRepoVariable(context.Background(), "", "", &github.ActionsVariable{
		Name:  "var0",
		Value: "value0",
	})
	return mServ
}
