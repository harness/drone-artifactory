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

#  Upload artifact to Jfrog Artifactory
This step uploads the artifact to Jfrog Artifactory.
A valid spec or a spec path given as an argument is mandatory.
The spec json format should be the same as Jfrog spec format

### Upload artifact to Jfrog Artifactory using spec path example:
```yaml
- step:
    type: Plugin
    name: UploadStep
    identifier: UploadStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: upload
        username: user
        password: <+secrets.getValue("jfrog_user")>
        build_name: gol-01
        build_number: 0.03.01
        url: https://URL.jfrog.io/artifactory/artifactory-test 
        module: pre-prod
        project: dyn_023
        spec_path: /harness/upload_spec.json
```
### Upload artifact to Jfrog Artifactory using spec as inline example:
```yaml
- step:
    type: Plugin
    name: UploadStep
    identifier: UploadStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: upload
        username: user
        password: <+secrets.getValue("jfrog_user")>
        build_name: gol-01
        build_number: 0.03.01
        url: https://URL.jfrog.io/artifactory/artifactory-test 
        spec: |
          {
          "files": [
            {
                "pattern": "/harness/game-of-life-master/**/*.jar",
                "target": "mvn_repo_deploy_snapshots_01/upload2/"
              },
              {
                "pattern": "/harness/game-of-life-master/**/*.java",
                "target": "mvn_repo_deploy_snapshots_01/upload/src2/"
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
