load(
    "//:repositories.bzl",
    "com_github_yugui_rules_ruby",
    "rules_proto_grpc_repos",
)

def ruby_repos(**kwargs):
    rules_proto_grpc_repos(**kwargs)
    com_github_yugui_rules_ruby(**kwargs)
