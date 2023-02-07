package poiana

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"net/http"
)

type MockVariableService struct {
	repoVariables map[string]string
	orgVariables  map[string]string
}

func (m MockVariableService) ListRepoVariables(_ context.Context, _, _ string, _ *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
	vars := make([]*github.ActionsVariable, 0)
	for key, val := range m.repoVariables {
		vars = append(vars, &github.ActionsVariable{
			Name:  key,
			Value: val,
		})
	}

	return &github.ActionsVariables{
		TotalCount: len(m.repoVariables),
		Variables:  vars,
	}, nil, nil
}

func (m MockVariableService) DeleteRepoVariable(_ context.Context, _, _, name string) (*github.Response, error) {
	delete(m.repoVariables, name)
	return nil, nil
}

func (m MockVariableService) CreateRepoVariable(_ context.Context, _, _ string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.repoVariables[variable.Name]; ok {
		// Key is present; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
		}, fmt.Errorf("variable already present")
	}
	// Key not present, add value
	m.repoVariables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (m MockVariableService) UpdateRepoVariable(_ context.Context, _, _ string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.repoVariables[variable.Name]; !ok {
		// There is no key; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
		}, fmt.Errorf("variable not present")
	}
	// Key is present, set value
	m.repoVariables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}

func (m MockVariableService) ListOrgVariables(ctx context.Context, owner string, opts *github.ListOptions) (*github.ActionsVariables, *github.Response, error) {
	vars := make([]*github.ActionsVariable, 0)
	for key, val := range m.orgVariables {
		vars = append(vars, &github.ActionsVariable{
			Name:  key,
			Value: val,
		})
	}

	return &github.ActionsVariables{
		TotalCount: len(m.orgVariables),
		Variables:  vars,
	}, nil, nil
}

func (m MockVariableService) DeleteOrgVariable(ctx context.Context, owner, name string) (*github.Response, error) {
	delete(m.orgVariables, name)
	return nil, nil
}

func (m MockVariableService) CreateOrgVariable(ctx context.Context, owner string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.orgVariables[variable.Name]; ok {
		// Key is present; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusForbidden,
			},
		}, fmt.Errorf("variable already present")
	}
	// Key not present, add value
	m.orgVariables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusCreated,
		},
	}, nil
}

func (m MockVariableService) UpdateOrgVariable(ctx context.Context, owner string, variable *github.ActionsVariable) (*github.Response, error) {
	if _, ok := m.orgVariables[variable.Name]; !ok {
		// There is no key; return error
		return &github.Response{
			Response: &http.Response{
				StatusCode: http.StatusNotFound,
			},
		}, fmt.Errorf("variable not present")
	}
	// Key is present, set value
	m.orgVariables[variable.Name] = variable.Value
	return &github.Response{
		Response: &http.Response{
			StatusCode: http.StatusNoContent,
		},
	}, nil
}

func NewMockVariableService() ActionsVarsService {
	mServ := &MockVariableService{
		repoVariables: make(map[string]string, 0),
		orgVariables:  make(map[string]string, 0),
	}
	// Initial variable set
	_, _ = mServ.CreateRepoVariable(context.Background(), "", "", &github.ActionsVariable{
		Name:  "repoVar0",
		Value: "repoValue0",
	})
	_, _ = mServ.CreateOrgVariable(context.Background(), "", &github.ActionsVariable{
		Name:  "orgVar0",
		Value: "orgValue0",
	})
	return mServ
}
