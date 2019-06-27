package main

var rustWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()`)

var rustLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("//{{ .Lang.Dir }}:rust_proto_lib.bzl", "rust_proto_lib")
load("@io_bazel_rules_rust//rust:rust.bzl", "rust_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_lib = kwargs.get("name") + "_lib"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var rustProtoLibraryRuleTemplate = mustTemplate(rustLibraryRuleTemplateString + `
    # Create lib file
    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
        grpc = False,
    )

    # Create {{ .Lang.Name }} library
    rust_library(
        name = kwargs.get("name"),
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
        ],
        visibility = kwargs.get("visibility"),
    )`)

var rustGrpcLibraryRuleTemplate = mustTemplate(rustLibraryRuleTemplateString + `
    # Create lib file
    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
        grpc = True,
    )

    # Create {{ .Lang.Name }} library
    rust_library(
        name = kwargs.get("name"),
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
            "@io_bazel_rules_rust//proto/raze:grpc",
            "@io_bazel_rules_rust//proto/raze:tls_api",
            "@io_bazel_rules_rust//proto/raze:tls_api_stub",
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeRust() *Language {
	return &Language{
		Dir:  "rust",
		Name: "rust",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "rust_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//rust:rust"},
				WorkspaceExample: rustWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates rust protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "rust_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//rust:rust", "//rust:grpc_rust"},
				WorkspaceExample: rustWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates rust protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "rust_proto_library",
				Kind:             "proto",
				Implementation:   rustProtoLibraryRuleTemplate,
				WorkspaceExample: rustWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates rust protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "rust_grpc_library",
				Kind:             "grpc",
				Implementation:   rustGrpcLibraryRuleTemplate,
				WorkspaceExample: rustWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates rust protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
