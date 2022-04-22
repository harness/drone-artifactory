# this script is used by the continuous integration server to
# build and publish the docker image for a commit to master.
$ErrorActionPreference = "Stop"

$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED="0"

if (-not (Test-Path env:VERSION)) {
    $env:VERSION="1809"
}

echo $env:GOOS
echo $env:GOARCH
echo $env:VERSION

# build the binary
Write-Host "+ go build -o release/windows/amd64/plugin.exe";
go build -o release/windows/amd64/plugin

# build and publish the docker image
docker login -u ${env:USERNAME} -p ${env:PASSWORD}
Write-Host "+ docker build -f docker/Dockerfile.windows.nano.amd64.${env:VERSION} -t plugins/artifactory:windows-${env:VERSION}-nano-amd64 .";
docker build -f docker/Dockerfile.windows.amd64.${env:VERSION} -t plugins/artifactory:windows-${env:VERSION}-nano-amd64 .
Write-Host "+ docker push plugins/artifactory:windows-${env:VERSION}-nano-amd64"
docker push plugins/artifactory:windows-${env:VERSION}-nano-amd64

# remove images from local cache
Write-Host "+ docker rmi plugins/artifactory:windows-${env:VERSION}-nano-amd64"
docker rmi plugins/artifactory:windows-${env:VERSION}-nano-amd64
