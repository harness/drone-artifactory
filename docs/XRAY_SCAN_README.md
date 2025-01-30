A plugin to scan artifacts in Jfrog Artifactory using Xray

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

#  Scan artifacts in Jfrog Artifactory using Xray
This step scans the artifacts in JFrog Artifactory using Xray. Suitable
watches and policies should be setup in JFrog for Xray in JFrog Artifactory.

### Scan artifacts in Jfrog Artifactory using Xray
```yaml
- step:
    type: Plugin
    name: XrayScanStep
    identifier: XrayScanStep
    spec:
      connectorRef: account.harnessImage
      image: plugins/artifactory:linux-amd64
      settings:
        command: scan
        url: https://URL.jfrog.io/xray
        username: user
        password: <+secrets.getValue("jfrog_user")>
        build_name: gol-01
        build_number: 0.03.01
```

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
