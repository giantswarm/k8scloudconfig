package token

import (
	"fmt"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/the-anna-project/id"
)

const (
	IDLength     = 3
	SecretLength = 8
)

func GenerateToken() (string, error) {
	tokenIDServiceConfig := id.DefaultServiceConfig()
	tokenIDServiceConfig.Length = IDLength

	tokenIDService, err := id.NewService(tokenIDServiceConfig)
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	tokenID, err := tokenIDService.New()
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	tokenSecretServiceConfig := id.DefaultServiceConfig()
	tokenSecretServiceConfig.Length = SecretLength

	tokenSecretService, err := id.NewService(tokenSecretServiceConfig)
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	tokenSecret, err := tokenSecretService.New()
	if err != nil {
		return "", microerror.MaskAny(err)
	}

	return fmt.Sprintf("%s.%s", tokenID, tokenSecret), nil
}
