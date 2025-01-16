package git

import (
	"strings"

	giturls "github.com/chainguard-dev/git-urls"
)

// AreSameRepository returns if 2 git urls are the same repository or not
func AreSameRepository(repoA, repoB string) bool {
	parsedRepoA, _ := giturls.Parse(repoA)
	parsedRepoB, _ := giturls.Parse(repoB)

	if parsedRepoA.Hostname() != parsedRepoB.Hostname() {
		return false
	}

	// In short SSH URLs like git@github.com:okteto/movies.git, path doesn't start with '/', so we need to remove it
	// in case it exists. It also happens with '.git' suffix. You don't have to specify it, so we remove in both cases
	repoPathA := strings.TrimSuffix(strings.TrimPrefix(parsedRepoA.Path, "/"), ".git")
	repoPathB := strings.TrimSuffix(strings.TrimPrefix(parsedRepoB.Path, "/"), ".git")

	return repoPathA == repoPathB
}
