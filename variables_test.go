package main

// import (
// 	"context"
// 	"github.com/FedeDP/GhEnvSet/pkg/poiana"
// 	"github.com/google/go-github/v49/github"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )

// type MockVariableService struct {
// 	variables map[string]string
// }

// func (m MockVariableService) ListRepoVariables(ctx context.Context, owner, repo string) (*poiana.Variables, error) {
// 	vars := make([]*poiana.Variable, 0)
// 	for key, val := range m.variables {
// 		vars = append(vars, &poiana.Variable{
// 			Name:  key,
// 			Value: val,
// 		})
// 	}

// 	return &poiana.Variables{
// 		TotalCount: len(m.variables),
// 		Variables:  vars,
// 	}, nil
// }

// func (m MockVariableService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error {
// 	delete(m.variables, name)
// 	return nil
// }

// func (m MockVariableService) CreateOrUpdateRepoVariable(ctx context.Context, owner, repo string, variable *poiana.Variable) error {
// 	m.variables[variable.Name] = variable.Value
// 	return nil
// }

// func newMockVariableService() poiana.ActionsVarsService {
// 	mServ := &MockVariableService{variables: make(map[string]string, 0)}
// 	// Initial variable set
// 	_ = mServ.CreateOrUpdateRepoVariable(context.Background(), "", "", &poiana.Variable{
// 		Name:      "test0",
// 		Value:     "value0",
// 		CreatedAt: github.Timestamp{},
// 		UpdatedAt: github.Timestamp{},
// 	})
// 	return mServ
// }

// func TestSyncVariables(t *testing.T) {
// 	ctx := context.Background()
// 	variables := map[string]string{
// 		"test1": "value1",
// 		"test2": "value2",
// 	}

// 	mockServ := newMockVariableService()
// 	err := syncVariables(ctx, mockServ, "", "", variables)
// 	assert.NoError(t, err)

// 	vars, err := mockServ.ListRepoVariables(ctx, "", "")
// 	assert.NoError(t, err)

// 	for _, v := range vars.Variables {
// 		assert.Equal(t, variables[v.Name], v.Value)
// 	}
// }
