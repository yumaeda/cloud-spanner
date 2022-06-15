# cloud-spanner
Connect to the specified Cloud Spanner and render retrieved data

## Prep

```zsh
export PROJECT_ID=hello-world-352201 \
       IMG_VERSION=1.0.10 \
       REGION=us-central1 \
       REPOSITORY=hello-world \
       IMG_NAME=hello-world \
       DEPLOYMENT=hello-world
```

## Build
```zsh
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMG_NAME}:${IMG_VERSION} .
```

## Push the Docker image to Artifact Registry
```zsh
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMG_NAME}:${IMG_VERSION}
```

## Update Kubernetes Deployment
```zsh
kubectl set image deployment/${DEPLOYMENT} ${IMG_NAME}=${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMG_NAME}:${IMG_VERSION}
```

## Configure environment variables for the Deployment 
```zsh
kubectl set env deployment/${DEPLOYMENT} SPINNER_PROJECT_ID=xxx
kubectl set env deployment/${DEPLOYMENT} DB_INSTANCE=xxx
kubectl set env deployment/${DEPLOYMENT} DB_NAME=xxx
kubectl set env deployment/${DEPLOYMENT} GOOGLE_APPLICATION_CREDENTIALS=/opt/app/key.json
```

## List the environment variables defined in the Development
```zsh
kubectl set env deployment/${DEPLOYMENT} --list
```

&nbsp;

## Run locally
```zsh
docker run --rm -d \
    -e SPINNER_PROJECT_ID=xxx \
    -e DB_INSTANCE=xxx \
    -e DB_NAME=xxx \
    -e GOOGLE_APPLICATION_CREDENTIALS=/opt/app/key.json \
    -p 8080:8080 \
    "${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${IMG_NAME}:${IMG_VERSION}"
```

## Misc
### Init Modules
```zsh
go mod init tokyo-takeout.com/api 
```
### Install Modules
```zsh
go install cloud.google.com/go/spanner
```

## References
- https://dockerlabs.collabnix.com/kubernetes/cheatsheets/kubectl.html
- https://cloud.google.com/spanner/docs/getting-started/go
- https://cloud.google.com/docs/authentication/production
- https://spanner.gcpug.jp/
