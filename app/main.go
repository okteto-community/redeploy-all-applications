package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/okteto-community/redeploy-all-applications/app/api"
	"github.com/okteto-community/redeploy-all-applications/app/model"
)

const redeployAppCommandTemplate = "okteto pipeline deploy -n \"%s\" --name \"%s\" --repository \"%s\" --branch \"%s\" --reuse-params --wait=%t"
const sleepNamespaceCommandTemplate = "okteto namespace sleep \"%s\""

func main() {
	token := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_CONTEXT")
	oktetoThreshold := os.Getenv("OKTETO_THRESHOLD")
	dryRun := os.Getenv("DRY_RUN") == "true"
	ignoreSleeping := os.Getenv("IGNORE_SLEEPING_NAMESPACES") == "true"
	restoreOriginalStatus := os.Getenv("RESTORE_ORIGINAL_NAMESPACE_STATUS") == "true"
	waitForDeploymentToFinish := os.Getenv("WAIT_FOR_DEPLOYMENT") == "true"

	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	if oktetoThreshold == "" {
		oktetoThreshold = "24h"
	}

	logger.Info(fmt.Sprintf("Configuration set: OKTETO_CONTEXT=%s OKTETO_THRESHOLD=%s DRY_RUN=%t IGNORE_SLEEPING_NAMESPACES=%t RESTORE_ORIGINAL_NAMESPACE_STATUS=%t WAIT_FOR_DEPLOYMENT=%t",
		oktetoURL, oktetoThreshold, dryRun, ignoreSleeping, restoreOriginalStatus, waitForDeploymentToFinish))

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

		if ns.Status == model.Sleeping && ignoreSleeping {
			logger.Info(fmt.Sprintf("Skipping namespace '%s' since its sleeping", ns.Name))
			logger.Info("-----------------------------------------------")
			continue
		}

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

			logger.Info(fmt.Sprintf("application '%s' lastupdate %s", app.Name, app.LastUpdated))
			if app.LastUpdated.After(updateThreshold) {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as it was updated recently", app.Name, ns.Name))
				continue
			}

			logger.Info(fmt.Sprintf("Redeploying application '%s' within namespace '%s'", app.Name, ns.Name))

			out, err := redeployApp(ns.Name, app.Name, app.Repository, app.Branch, waitForDeploymentToFinish, dryRun)
			if err != nil {
				logger.Error(fmt.Sprintf("There was an error redeploying the application '%s' within namespace '%s': %s", app.Name, ns.Name, err))
			} else {
				logger.Info(out)
			}
		}

		if restoreOriginalStatus && ns.Status == model.Sleeping {
			out, err := sleepNamespace(ns.Name, dryRun)
			if err != nil {
				logger.Error(fmt.Sprintf("There was an error sleeping namespace '%s': %s", ns.Name, err))
			} else {
				logger.Info(out)
			}
		}
		logger.Info("-----------------------------------------------")
	}
}

// redeployApp executes the Okteto CLI command to redeploy an application
func redeployApp(ns, appName, repo, branch string, wait, dryRun bool) (string, error) {
	cmdStr := fmt.Sprintf(redeployAppCommandTemplate, ns, appName, repo, branch, wait)
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

// sleepNamespace executes the Okteto CLI command to sleep a namespace
func sleepNamespace(ns string, dryRun bool) (string, error) {
	cmdStr := fmt.Sprintf(sleepNamespaceCommandTemplate, ns)
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
