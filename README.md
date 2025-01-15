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

# Testing

Execute the plugin from your current working directory:

```text
docker run --rm \
  -e PLUGIN_USERNAME=foo \
  -e PLUGIN_PASSWORD=bar \
  -e PLUGIN_URL=<url> \
  -e PLUGIN_SOURCE=/drone/README.md \
  -e PLUGIN_TARGET=/pcf \
  -e PLUGIN_BUILD_NAME=buildName \
  -e PLUGIN_BUILD_NUMBER=4 \
  -e PLUGIN_TAREGT_PROPS='key1=value1,key2=value2'
  -v $(pwd):/drone \
  plugins/artifactory
```

## Harness CI Example:
```yaml
              - step:
                  type: Plugin
                  name: jFrog-Test
                  identifier: Email_Plugin
                  spec:
                    connectorRef: account.harnessImage
                    image: plugins/artifactory:linux-amd64
                    settings:
                      access_token: <JFROG_ACCESS_TOKEN>
                      url: https://URL.jfrog.io/artifactory/artifactory-test/
                      source: /harness/cache.txt
                      target: newdemo
                      build_name: <+pipeline.name>
                      build_number: <+pipeline.executionId>
                      target_props: key1=value1,key2=value2
```
### Maven Build and Publish reference
[Go to Maven reference](./docs/MAVEN_README.md)

### Gradle Build and Publish reference
[Go to Gradle reference](./docs/GRADLE_README.md)

## Community and Support
[Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.

[Harness Community Forum](https://community.harness.io/) - Ask questions, find answers, and help other users.

[Report and Track A Bug](https://community.harness.io/c/bugs/17) - Find a bug? Please report in our forum under Drone Bugs. Please provide screenshots and steps to reproduce. 

[Events](https://www.meetup.com/harness/) - Keep up to date with Drone events and check out previous events [here](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
