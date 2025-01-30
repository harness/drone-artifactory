A plugin to download files from Jfrog artifactory.

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

#  Download artifact to Jfrog Artifactory
This step downloads the artifacts from Jfrog Artifactory.
A valid spec or a spec path given as an argument is mandatory.
The spec json format should be the same as Jfrog spec format

### Download artifact to Jfrog Artifactory using spec path example:
```yaml
- step:
    type: Plugin
    name: DownloadStep
    identifier: DownloadStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: download
        username: user
        password: <+secrets.getValue("jfrog_user")>
        build_name: gol-01
        build_number: 0.03.01
        url: https://URL.jfrog.io/artifactory
        module: pre-prod
        project: dyn_023
        spec_path: /harness/download_spec.json
```

### Download artifact to Jfrog Artifactory using spec as inline example:
```yaml
- step:
    type: Plugin
    name: DownloadStep
    identifier: DownloadStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: download
        username: user
        password: <+secrets.getValue("jfrog_user")>
        build_name: gol-01
        build_number: 0.03.01
        url: https://URL.jfrog.io/artifactory 
        spec: |
          {
            "files": [
              {
                "pattern": "mvn_repo_deploy_snapshots_01/**/*.jar",
                "target": "./downloads/jars/",
                "flat": "false",
                "recursive": "true"
              },
              {
                "pattern": "mvn_repo_deploy_snapshots_01/**/*.java",
                "target": "./downloads/src/",
                "flat": "false",
                "recursive": "true"
              }
            ]
          } 
        module: pre-prod
        project: dyn_023
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
