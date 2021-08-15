package shared

import "github.com/fr123k/pulumi-wireguard-aws/pkg/model"

func WireguardUserData() (*model.UserData, error) {
    userDataVariables := map[string]string{
        "{{ CLIENT_PUBLICKEY }}":        "CLIENT_PUBLICKEY",
        "{{ CLIENT_IP_ADDRESS }}":       "CLIENT_IP_ADDRESS",
        "{{ MAILJET_API_CREDENTIALS }}": "MAILJET_API_CREDENTIALS",
        "{{ METADATA_URL }}":            "METADATA_URL",
    }

    userData, err := model.NewUserData("cloud-init/wireguard.txt", model.TemplateVariablesEnvironment(userDataVariables))
    if err != nil {
        return nil, err
    }
    return userData, nil
}
