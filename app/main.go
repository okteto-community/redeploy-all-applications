package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/okteto-community/redeploy-all-applications/app/api"
)

const redeployAppCommandTemplate = "okteto pipeline deploy -n \"%s\" --name \"%s\" --repository \"%s\" --branch \"%s\" --reuse-params --wait=false"

func main() {
	token := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_URL")
	oktetoThreshold := os.Getenv("OKTETO_THRESHOLD")
	dryRun := os.Getenv("DRY_RUN") == "true"

	if oktetoThreshold == "" {
		oktetoThreshold = "24h"
	}

	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	if token == "" || oktetoURL == "" {
		logger.Error("OKTETO_TOKEN, OKTETO_URL environment variables are required")
		os.Exit(1)
	}

	if dryRun {
		logger.Info("Dry run mode is enabled")
	}

	u, err := url.Parse(oktetoURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid OKTETO_URL %s", err))
		os.Exit(1)
	}

	threshold, err := time.ParseDuration(oktetoThreshold)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid OKTETO_THRESHOLD %s", err))
		os.Exit(1)
	}

	nsList, err := api.GetNamespaces(u.Host, token, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error requesting the namespaces: %s", err))
		os.Exit(1)
	}

	// We check for applications that were last updated more than 24 hours ago
	updateThreshold := time.Now().Add(-threshold)
	for _, ns := range nsList {
		logger.Info(fmt.Sprintf("Processing namespace '%s'", ns.Name))

		applications, err := api.GetApplicationsWithinNamespace(u.Host, token, ns.Name, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("There was an error requesting the applications within namespace '%s': %s", ns.Name, err))
			logger.Info("-----------------------------------------------")
			continue
		}

		for _, app := range applications {
			if app.Repository == "" {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as does not have a repository", app.Name, ns.Name))
				continue
			}

			if app.LastUpdated.After(updateThreshold) {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as it was updated recently", app.Name, ns.Name))
				continue
			}

			logger.Info(fmt.Sprintf("Redeploying application '%s' within namespace '%s'", app.Name, ns.Name))

			out, err := redeployApp(ns.Name, app.Name, app.Repository, app.Branch, dryRun)
			if err != nil {
				logger.Error(fmt.Sprintf("There was an error redeploying the application '%s' within namespace '%s': %s", app.Name, ns.Name, err))
			} else {
				logger.Info(out)
			}
		}
		logger.Info("-----------------------------------------------")
	}
}

// redeployApp executes the Okteto CLI command to redeploy an application
func redeployApp(ns, appName, repo, branch string, dryRun bool) (string, error) {
	cmdStr := fmt.Sprintf(redeployAppCommandTemplate, ns, appName, repo, branch)
	cmd := exec.Command("bash", "-c", cmdStr)

	if !dryRun {
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}

		return string(out), nil
	}

	return fmt.Sprintf("[DRY MODE] %s", cmdStr), nil
}
