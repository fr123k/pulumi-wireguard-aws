package shared

import "github.com/fr123k/pulumi-wireguard-aws/pkg/model"

// prebakedUserData is a shared helper that loads a minimal cloud-init template
// for a pre-baked image and substitutes the given environment variables.
// This is used by all pre-baked image variants (temporal, wireguard, etc.) to
// avoid duplicating the NewUserData + TemplateVariablesEnvironment boilerplate.
func prebakedUserData(templatePath string, vars map[string]string) (*model.UserData, error) {
	return model.NewUserData(templatePath, model.TemplateVariablesEnvironment(vars))
}
