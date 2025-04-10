A plugin to upload files to Jfrog artifactory.

Run the following script to install git-leaks support to this repo.
```
chmod +x ./git-hooks/install.sh
./git-hooks/install.sh
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t plugins/artifactory  -f docker/Dockerfile .
```

# Gradle Build and Publish
- Gradle build step is used to build the Gradle project and create artifacts.
- Publish step is used to publish the Gradle project artifacts to the artifactory repositories.
- Authentication for Jfrog artifactory can be done using Username and Password or Access Token. Refer to below examples.
- Additional build discard with below parameters can be done.
    - delete_artifacts: The flag to delete the artifacts, if not set will only delete build metadata.
    - exclude_builds: The builds to exclude from deletion.
    - max_builds: The maximum number of builds to keep.
    - max_days: The maximum number of days to keep the builds based on the build timestamp as start time.
    - async: The flag to run the step asynchronously.

### Gradle Build step example using Username and Password:
```yaml
- step:
  type: Plugin
  name: Plugin_gradle_run
  identifier: Plugin_gradle_run
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: gradle
      username: user
      password: <+secrets.getValue("jfrog_user")>
      url: https://URL.jfrog.io/artifactory/artifactory-test/
      repo_resolve: repo_resolve_gradle
      repo_deploy: repo_deploy_gradle
      tasks: clean build
      build_name: t2
      build_number: t4
      threads: "3"
      project_key: new_dev_test
```

### Gradle Publish step example using Username and Password:
```yaml
- step:
  type: Plugin
  name: Plugin_gradle_publish
  identifier: Plugin_gradle_publish
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: gradle
      command: publish
      url: https://URL.jfrog.io/artifactory/artifactory-test/
      username: user
      password: <+secrets.getValue("jfrog_user")>
      build_name: t2
      build_number: t4
      repo_resolve: repo_resolve_gradle_02
      repo_deploy: repo_deploy_gradle_02
      deployer_id: gradle-deployer
```

### Gradle Build step example using Access Token
The config is same as the "Gradle Build step example using Username and Password" for gradle
"username" should be set as a valid username using which the access token was created
"password" should be set as the access token value, access token will be a very long string

### Gradle Publish step example using Access Token
The config is same as the "Gradle Publish step example using Username and Password" for gradle
"username" should be set as a valid username using which the access token was created
"password" should be set as the access token value, access token will be a very long string


### Gradle Publish step with build discard additional parameters
```yaml
- step:
  type: Plugin
  name: Plugin_gradle_publish
  identifier: Plugin_gradle_publish
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: gradle
      command: publish
      url: https://URL.jfrog.io
      username: user
      password: <+secrets.getValue("jfrog_password")>
      build_name: gradle02
      build_number: 3
      repo_resolve: repo_resolve_gradle
      repo_deploy: repo_deploy_gradle_02
      deployer_id: gradle-deployer-01
      max_builds: 3
      max_days: 30
      delete_artifacts: true
      exclude_builds: gradle-01,gradle-02
      async: true
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
