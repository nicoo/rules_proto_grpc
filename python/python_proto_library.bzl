load("//python:python_proto_compile.bzl", "python_proto_compile")

def python_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create python library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ],
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

# Alias
py_proto_library = python_proto_library
