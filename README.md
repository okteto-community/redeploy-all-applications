# Redeploy applications in all Okteto namespaces

> This is an experiment and Okteto does not officially support it.

- Create an [Okteto Admin Token](https://www.okteto.com/docs/admin/dashboard/#admin-access-tokens)

- Export the token to a local variable:

```bash
export OKTETO_ADMIN_TOKEN=<<your-token>>
```

- Export the URL of your Okteto instance to a local variable:

```bash
export OKTETO_URL=<<your-okteto-url>>
```

- (Optional) Set the threshold since last update of an application to the OKTETO_THRESHOLD local variable. If an application's has been updated before the threshold, it will be re-deployed. If the application is has been updated more recently than the threshold, it won't be re-deployed. **Default is "24h"**.

```bash
export OKTETO_THRESHOLD=<<your-okteto-threshold>>
```

- (Optional) Set whether the command is a dry run or not to the DRY_RUN local variable. **Default is "false"**.

```bash
export DRY_RUN=<<dry-run mode, true or false>>
```

- Build and run the script:
```bash
cd ./app
go build main.go
./main
```
