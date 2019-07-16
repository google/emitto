load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.18.5/rules_go-0.18.5.tar.gz"],
    sha256 = "a82a352bffae6bee4e95f68a8d80a70e87f42c4741e6a448bec11998fcc82329",
)

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "87fc6a2b128147a0a3039a2fd0b53cc1f2ed5adb8716f50756544a572999ae9a",
    strip_prefix = "rules_docker-0.8.1",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.8.1.tar.gz"],
)

# Go
load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

# Docker
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)
container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)
_go_image_repos()

# Gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
gazelle_dependencies()

# Dependencies
go_repository(
    name = "com_github_google_uuid",
    commit = "c2e93f3ae59f2904160ceaab466009f965df46d6",
    importpath = "github.com/google/uuid",
)

go_repository(
    name = "com_github_google_fleetspeak",
    commit = "bc95dd6941494461d2e5dff0a7f4c78a07ff724d",
    importpath = "github.com/google/fleetspeak",
)

go_repository(
    name = "com_github_spf13_afero",
    commit = "588a75ec4f32903aa5e39a2619ba6a4631e28424",
    importpath = "github.com/spf13/afero",
)

go_repository(
    name = "com_github_google_go_cmp",
    commit = "917e382dab80060fd1f094402bfbb5137ec3c4ff",
    importpath = "github.com/google/go-cmp",
)

go_repository(
    name = "com_github_fatih_camelcase",
    commit = "9db1b65eb38bb28986b93b521af1b7891ee1b04d",
    importpath = "github.com/fatih/camelcase",
)

go_repository(
    name = "org_golang_google_api",
    commit = "b50168921e183c95644f04416271d702a730550c",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "5ea6847e42e62cd77aab6fb7047b15132174a634",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "aaccbc9213b0974828f81aaac109d194880e3014",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "com_github_googleapis_gax_go",
    commit = "bd5b16380fd03dc758d11cef74ba2e3bc8b0e8c2",
    importpath = "github.com/googleapis/gax-go",
)

go_repository(
    name = "io_opencensus_go",
    commit = "6325d764b2d4a66576c5623aa1e6010b4148a429",
    importpath = "go.opencensus.io",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    commit = "59383c442f7d7b190497e9bb8fc17a48d06cd03f",
    importpath = "github.com/hashicorp/golang-lru",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "532a0b98cb9580f72fd376b539f9eb984f92e054",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "com_github_golang_glog",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    importpath = "github.com/golang/glog",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "eb0b1bdb6ae60fcfc41b8d907b50dfb346112301",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "com_github_google_emitto",
    commit = "0c93e985f54f1fedf41251553458150c12642e5a",
    importpath = "github.com/google/emitto",
)
