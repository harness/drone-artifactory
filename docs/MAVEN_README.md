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

# Maven Build and Publish
- Maven build step is used to build the maven project and create artifacts. 
- Publish step is used to publish the maven project artifacts to the artifactory repositories.
- Authentication for Jfrog artifactory can be done using Username and Password or Access Token. Refer to below examples.

### Maven Build step example using Username and Password:
```yaml
- step:
  type: Plugin
  name: MavenBuildTest
  identifier: MavenBuildTest
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: mvn
      username: user
      password: <+secrets.getValue("jfrog_user")>
      pom: pom.xml
      goals: clean install
      build_name: t2
      build_number: t4
      url: https://URL.jfrog.io/artifactory/artifactory-test
      resolver_id: resolve_gen_maven
      resolve_release_repo: mvn_repo_resolve_releases
      resolve_snapshot_repo: mvn_repo_resolve_snapshots
```

### Maven  Publish step example using Username and Password:
```yaml
- step:
  type: Plugin
  name: MavenPublishTest
  identifier: MavenPublishTest
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: mvn
      command: publish
      url: https://trialqlrico.jfrog.io
      username: user
      password: <+secrets.getValue("jfrog_user")>
      build_name: t2
      build_number: t4
      deployer_id: mvn-deployer
      deploy_release_repo: mvn_repo_deploy_releases
      deploy_snapshot_repo: mvn_repo_deploy_snapshots
```

### Maven Build step example using Access Token:
```yaml
- step:
  type: Plugin
  name: MavenBuildTest
  identifier: MavenBuildTest
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: mvn
      access_token: <+secrets.getValue("jfrog_user")>
      pom: pom.xml
      goals: clean install
      build_name: t2
      build_number: t4
      url: https://URL.jfrog.io/artifactory/artifactory-test/
      resolver_id: resolve_gen_maven
      resolve_release_repo: mvn_repo_resolve_releases
      resolve_snapshot_repo: mvn_repo_resolve_snapshots
```

### Maven Publish step example using Access Token:
```yaml
- step:
  type: Plugin
  name: MavenPublishTest
  identifier: MavenPublishTest
  spec:
    connectorRef: account.harnessImage
    image: plugins/artifactory:linux-amd64
    settings:
      build_tool: mvn
      command: publish
      url: https://URL.jfrog.io/artifactory/artifactory-test/
      access_token: <+secrets.getValue("jfrog_user")>
      resolver_id: resolve_gen_maven
      build_name: t2
      build_number: t4
      deployer_id: mvn-deployer
      deploy_release_repo: mvn_repo_deploy_releases
      deploy_snapshot_repo: mvn_repo_deploy_snapshots
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
