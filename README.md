# Emitto

### About

Emitto is a service which provides a robust and targeted way to manage, store,
and administer [Suricata](https://suricata-ids.org/) intrusion detection system
(IDS) rules to a distributed network sensor monitoring deployment.

### Building

1) Install [Bazel](https://bazel.build/)

2) Download the source code

3) Build the project:

```bash
cd emitto/
bazel build //...
```

#### Prerequisites

The following services and products must be established and configured before
running Emitto.

##### Fleetspeak

Emitto leverages [Fleetspeak](https://github.com/google/fleetspeak) for
reliable, multi-homed communication between the server and network sensors.

While this code is working internally with an internal deployment of
Fleetspeak,
the open source version is still in development. See the Fleetspeak
[status](https://github.com/google/fleetspeak#status) section for more
information and updates.

Please read the Fleetspeak
[disclaimer](https://github.com/google/fleetspeak#disclaimer).

##### Google Cloud Platform

By default, Emitto uses [Cloud Storage](https://cloud.google.com/storage/) and
[Cloud Filestore](https://cloud.google.com/filestore/) for object and rule file
storage, respectively.

### Discussion & Announcements

The [Emitto](https://groups.google.com/forum/#!forum/emitto) Google Groups
forum
will be used for community discussions and announcements.

### DISCLAIMER

This is not an officially supported Google product.
