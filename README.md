# Emitto

[![Travis Build Status](https://api.travis-ci.org/google/emitto.svg?branch=master)](https://travis-ci.org/google/emitto)
[![Go Report Card](https://goreportcard.com/badge/github.com/google/emitto)](https://goreportcard.com/report/github.com/google/emitto)

## About

Emitto is a service which provides a robust and targeted way to manage, store,
and administer [Suricata](https://suricata-ids.org/) intrusion detection system
(IDS) rules to a distributed network sensor monitoring deployment.

## Building

1) Install [Bazel](https://bazel.build/)

2) Download the source code

3) Build the project:

```bash
cd emitto/
bazel build //...
```

### Docker

Both the sever and sensor client can be built with the Bazel Docker rules located in their respective BUILD files.

For example, to build and run the sensor client locally in Docker:

```bash
bazel --bazelrc=.bazelrc build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 source/sensor:sensor_image.tar
docker load -i bazel-bin/source/sensor/sensor_image.tar
docker run bazel/source/sensor:sensor_image
```

The current `go_image` rule uses a [distroless](https://github.com/GoogleContainerTools/distroless)
runtime base image, and will run using the default Go binary flag values. The `container_image` rule can be extended to support a custom `env` and other image attributes.

Bazel supports the ability to push an image to a remote repository, and also to pull a remote image to be
used as a custom base dependency. More information about Bazel Docker rules can be found [here](https://github.com/bazelbuild/rules_docker).

### Prerequisites

The following services and products must be established and configured before
running Emitto.

#### Fleetspeak

Emitto leverages [Fleetspeak](https://github.com/google/fleetspeak) for
reliable, multi-homed communication between the server and network sensors.

While this code is working internally with an internal deployment of
Fleetspeak,
the open source version is still in development. See the Fleetspeak
[status](https://github.com/google/fleetspeak#status) section for more
information and updates.

Please read the Fleetspeak
[disclaimer](https://github.com/google/fleetspeak#disclaimer).

#### Google Cloud Platform

By default, Emitto uses [Google Cloud Datastore](https://cloud.google.com/datastore/)
and [Google Cloud Storage](https://cloud.google.com/storage/) for object and rule file
storage, respectively.

## Discussions & Announcements

The [Emitto](https://groups.google.com/forum/#!forum/emitto) Google Groups
forum
will be used for community discussions and announcements.

## DISCLAIMER

This is not an officially supported Google product.
