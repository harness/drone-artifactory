A plugin to manage artifacts in Jfrog artifactory.

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t plugins/artifactory  -f docker/Dockerfile .
```

# Build Discard CI step
- Build discard step is used to delete the old builds from the Jfrog artifactory.
- Authentication for Jfrog artifactory can be done using Username and Password or Access Token. Refer to below examples.
- Use this step as generic build discard for cases like an upload with a build number or a build name.
- The build discard step will delete the builds from the Jfrog artifactory based on the following parameters:
  - url: The URL of the Jfrog artifactory.
  - username: The username for authentication.
  - password: The password for authentication.
  - access_token: The access token for authentication.
  - build_name: The name of the build.
  - delete_artifacts: The flag to delete the artifacts, if not set will only delete build metadata.
  - exclude_builds: The builds to exclude from deletion.
  - max_builds: The maximum number of builds to keep.
  - max_days: The maximum number of days to keep the builds based on the build timestamp as start time.
  - async: The flag to run the step asynchronously.
 
### Build Discard step example using Username and Access Token:
```yaml
- step:
    type: Plugin
    name: PluginBdi3
    identifier: PluginBdi3
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        url: https://URL.jfrog.io/artifactory/artifactory-test
        username: user
        access_token: <+secrets.getValue("jfrog_access_token")>
        build_name: my-build
        command: build-discard
        max_builds: 3
        max_days: 30
        delete_artifacts: true
        exclude_builds: my-build-1,my-build-2
        async: true
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce.

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
