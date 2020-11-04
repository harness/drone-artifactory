A plugin to upload files to Jfrog artifactory.

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t drone/drone-artifactory -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

```text
docker run --rm \
  -e PLUGIN_USERNAME=foo \
  -e PLUGIN_PASSWORD=bar \
  -e PLUGIN_URL=<url> \
  -e PLUGIN_SRC_FILE=/drone/README.md \
  -e PLUGIN_TARGET_FILE=/pcf \
  -v $(pwd):/drone \
  drone/drone-artifactory
```