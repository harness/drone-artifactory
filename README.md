A plugin to upload files to Jfrog artifactory.

# Building

Build the plugin binary:

```bash
scripts/build.sh
```

Build the plugin image:

```bash
docker build -t plugins/artifactory  -f docker/Dockerfile .
```

# Testing

Execute the plugin from your current working directory:

```bash
docker run --rm \
  -e PLUGIN_USERNAME=foo \
  -e PLUGIN_PASSWORD=bar \
  -e PLUGIN_URL=<url> \
  -e PLUGIN_SOURCE=/drone/README.md \
  -e PLUGIN_TARGET=/pcf \
  -v $(pwd):/drone \
  plugins/artifactory
```
