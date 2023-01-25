package poiana

import (
	"context"
	"fmt"
	"net/http"
)

type ActionsVarsService interface {
	ListRepoVariables(ctx context.Context, owner, repo string) (*Variables, error)
	DeleteRepoVariable(ctx context.Context, owner, repo, name string) error
	CreateOrUpdateRepoVariable(ctx context.Context, owner, repo string, variable *Variable) error
}

// ListRepoVariables lists all variables available in a repository
// without revealing their encrypted values.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#list-repository-variables
func (s *actionsService) ListRepoVariables(ctx context.Context, owner, repo string) (*Variables, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/variables", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	variables := new(Variables)
	_, err = s.client.Do(ctx, req, &variables)
	if err != nil {
		return nil, err
	}
	return variables, nil
}

// GetRepoVariable gets a single repository variable without revealing its encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#get-a-repository-variable
func (s *actionsService) GetRepoVariable(ctx context.Context, owner, repo, name string) (*Variable, error) {
	u := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, name)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	variable := new(Variable)
	_, err = s.client.Do(ctx, req, variable)
	if err != nil {
		return nil, err
	}
	return variable, nil
}

// CreateOrUpdateRepoVariable creates or updates a repository variable.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#create-or-update-a-repository-variable
func (s *actionsService) CreateOrUpdateRepoVariable(ctx context.Context, owner, repo string, variable *Variable) error {
	u := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, variable.Name)
	req, err := s.client.NewRequest("PATCH", u, variable)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if resp.StatusCode != http.StatusNoContent {
			u = fmt.Sprintf("repos/%v/%v/actions/variables", owner, repo)
			req, err = s.client.NewRequest("POST", u, variable)
			if err == nil {
				_, err = s.client.Do(ctx, req, nil)
			}
		}
	}
	return err
}

// DeleteRepoVariable deletes a variable in a repository using the variable name.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#delete-a-repository-variable
func (s *actionsService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error {
	u := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, name)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err == nil {
		_, err = s.client.Do(ctx, req, nil)
	}
	return err
}
