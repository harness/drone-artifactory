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

#  Adds dependencies from the local file-system to the build info
This step is used to add dependencies from the local file-system to the build info.
The dependencies are added to the build info in the Artifactory server. 

###  Add dependencies with a dependency pattern:
```yaml
- step:
    type: Plugin
    name: AddBuildDependencyStep
    identifier: AddBuildDependencyStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: add-build-dependencies
        url: https://URL.jfrog.io
        username: user
        password: <+secrets.getValue("jfrog_user")>
        dependency: /harness/new_build/\*\*/\*.jar
        exclusions: /harness/add_deps/exclude_test/\*\*
        build_name: gol-01
        build_number: 0.03.01
        module: test_module
```

###  Add dependencies with a spec file path:
```yaml
- step:
    type: Plugin
    name: AddBuildDependencyStep
    identifier: AddBuildDependencyStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: add-build-dependencies
        url: https://URL.jfrog.io
        username: user
        password: <+secrets.getValue("jfrog_user")>        
        spec_path: /harness/build_add_deps_spec01.json
        build_name: gol-01
        build_number: 0.03.01
        module: test_module
```

###  Add dependencies of Artifactory files to build info:
```yaml
- step:
    type: Plugin
    name: AddBuildDependencyStep
    identifier: AddBuildDependencyStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: add-build-dependencies
        url: https://URL.jfrog.io
        username: user
        password: <+secrets.getValue("jfrog_user")>
        from_rt: true
        dependency: mvn_repo_resolve_releases_01/com/\*\*/\*.jar        
        build_name: gol-01
        build_number: 0.03.01
        module: test_module
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
