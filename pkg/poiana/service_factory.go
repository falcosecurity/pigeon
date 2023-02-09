package poiana

type ServiceFactory interface {
	NewOrgVariableService(org string) ActionsVarsService
	NewOrgSecretService(org string) ActionsSecretsService
	NewRepoVariableService(org, repo string) ActionsVarsService
	NewRepoSecretService(org, repo string) ActionsSecretsService
}
