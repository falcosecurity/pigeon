package poiana

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-github/v49/github"
	"github.com/google/go-querystring/query"
)

// Secret represents a repository action variable.
type Variable struct {
	Name      string           `json:"name"`
	Value     string           `json:"value"`
	CreatedAt github.Timestamp `json:"created_at"`
	UpdatedAt github.Timestamp `json:"updated_at"`
}

type Variables struct {
	TotalCount int         `json:"total_count"`
	Variables  []*Variable `json:"variables"`
}

type actionsService struct {
	*github.ActionsService
	client *github.Client
}

type Client struct {
	*github.Client
	Actions *actionsService
}

func NewClient(c *github.Client) *Client {
	return &Client{
		Client: c,
		Actions: &actionsService{
			ActionsService: c.Actions,
			client:         c,
		},
	}
}

func httpAddOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}
	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// ListRepoVariables lists all variables available in a repository
// without revealing their encrypted values.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#list-repository-variables
func (s *actionsService) ListRepoVariables(ctx context.Context, owner, repo string, opts *github.ListOptions) (*Variables, *github.Response, error) {
	url := fmt.Sprintf("repos/%v/%v/actions/variables", owner, repo)
	u, err := httpAddOptions(url, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	variables := new(Variables)
	resp, err := s.client.Do(ctx, req, &variables)
	if err != nil {
		return nil, resp, err
	}
	return variables, resp, nil
}

// GetRepoVariable gets a single repository variable without revealing its encrypted value.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#get-a-repository-variable
func (s *actionsService) GetRepoVariable(ctx context.Context, owner, repo, name string) (*Variable, *github.Response, error) {
	url := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, name)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	variable := new(Variable)
	resp, err := s.client.Do(ctx, req, variable)
	if err != nil {
		return nil, resp, err
	}
	return variable, resp, nil
}

// CreateOrUpdateRepoVariable creates or updates a repository variable.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#create-or-update-a-repository-variable
func (s *actionsService) CreateOrUpdateRepoVariable(ctx context.Context, owner, repo string, variable *Variable) (*github.Response, error) {
	url := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, variable.Name)
	req, err := s.client.NewRequest("PATCH", url, variable)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if resp.StatusCode != http.StatusNoContent {
			url = fmt.Sprintf("repos/%v/%v/actions/variables", owner, repo)
			req, err = s.client.NewRequest("POST", url, variable)
			if err != nil {
				return nil, err
			}
			return s.client.Do(ctx, req, nil)
		}
	}
	return resp, err
}

// DeleteRepoVariable deletes a variable in a repository using the variable name.
//
// GitHub API docs: https://docs.github.com/en/rest/actions/variables#delete-a-repository-variable
func (s *actionsService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) (*github.Response, error) {
	url := fmt.Sprintf("repos/%v/%v/actions/variables/%v", owner, repo, name)
	req, err := s.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(ctx, req, nil)
}
