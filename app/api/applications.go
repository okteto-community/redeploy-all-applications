package api

import (
	"fmt"
	"log/slog"

	"github.com/okteto-community/redeploy-all-applications/app/model"
)

// GetApplicationsWithinNamespace retrieves all the applications deployed within a namespace
func GetApplicationsWithinNamespace(baseURL, token, namespace string, logger *slog.Logger) ([]model.Application, error) {
	applicationsURLs := fmt.Sprintf("https://%s/%s", baseURL, fmt.Sprintf(applicationsAPIPathTemplate, namespace))
	var applications []model.Application
	if err := sendRequest(applicationsURLs, token, &applications, logger); err != nil {
		return nil, err
	}
	return applications, nil
}
