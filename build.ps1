if ([string]::IsNullOrEmpty($IMAGE_VERSION)) {
    Write-Output "Environment variable IMAGE_VERSION is not set, exiting..."
    return
}


$imageName="ghcr.io/katasec/aproxy:$IMAGE_VERSION"
docker build . -t $imageName
docker push $imageName
#docker run -it ghcr.io/katasec/aproxy:v0.0.2
#docker run -e APROXY_TARGET_URL="https://go.dev" -e APROXY_TARGET_PORT="1337" -it ghcr.io/katasec/aproxy:v0.0.2