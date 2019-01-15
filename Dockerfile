## BUILDER IMAGE
# this intermediate image is only used to compile the go binary
FROM golang:1.11.4
WORKDIR /go
COPY . /go/src/simple-database-provisioner
RUN GOOS=linux go install simple-database-provisioner...

## SIMPLE DATABASE PROVISIONER SERVICE IMAGE
# this image is the service image to be run inside the kubernetes cluster
FROM centos:7
LABEL maintainer="opensource@ecodia.de"
RUN groupadd -g 10000 appuser && useradd -r -u 10000 -g appuser appuser
RUN mkdir -p /persistence && mkdir -p /app/crds
WORKDIR /app/
RUN touch /app/config.yaml
COPY crds ./crds
COPY --from=0 /go/bin/controller .
RUN chown -R appuser:appuser /app && chown -R appuser:appuser /persistence
USER appuser
CMD ["./controller"]
