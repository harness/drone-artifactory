A plugin to upload files to Jfrog artifactory.

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
  -v $(pwd):/drone \
  plugins/artifactory
```

## Community and Support

* [Harness Community Slack](https://join.slack.com/t/harnesscommunity/shared_invite/zt-25b35u8j5-qAvb~7FJ1NFXbiW4AN101w) - Join the #drone slack channel to connect with our engineers and other users running Drone CI.
* [Drone FAQs on Harness Developer Hub](https://developer.harness.io/kb/continuous-integration/drone-faqs)
* [Harness Community](https://developer.harness.io/community)
* [Harness Events](https://www.harness.io/events) - You can check out previous events on [YouTube](https://www.youtube.com/watch?v=Oq34ImUGcHA&list=PLXsYHFsLmqf3zwelQDAKoVNmLeqcVsD9o).
