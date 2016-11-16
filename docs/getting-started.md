# Getting started

The Getting Started guide shows you how to download and install DapperDox; how to quickly
get an OpenAPI Swagger 2.0 specification rendered, and some useful configuration options
to tune the presentation.

## Download DapperDox

[Download the latest release](/downloads/download) of DapperDox for you operating system, and unpack it.

```
tar -zxf swaggerly-*.tgz
cd swaggerly
```

## Basic configuration

DapperDox requires little configuration to get started, just where to find your OpenAPI Swagger 2.0
specifications. Configuration parameters can be specified on the command line, or by setting
environment variables.

For this example, we will specify the Swagger specification location on the command line:
```
./swaggerly -spec-dir=examples/specifications/petstore
```

This will start serving the example <em>petstore</em> specification included in the DapperDox distribution.

By default, DapperDox will start serving reference documentation on port 3123, so point your web browser at
[http://localhost:3123](http://localhost:3123)

## Serving your own specification

Point DapperDox at the directory containing your Swagger specification using the `-spec-dir` configuration parameter,
where DapperDox will look for the file `swagger.json`. If your specification has a different filename, specify
this with the `-spec-filename` parameter.

If your specification is broken down into multiple files, with some in sub-directories, then you must point
DapperDox at the parent directory, beneath which it can find all the files it needs. If the main Swagger
specification file exists in a sub-directory of this parent, then configure that using the `-spec-filename` parameter.

For example, consider the following file structure:

```
/user/api_specs/spec/swagger.json
/user/api_specs/operations/get_user.json
/user/api_specs/operations/post_user.json
/user/api_specs/definitions/user_resource.json
```

Since the common parent directory of all these resource files is `api_specs`, this is where you need to tell DapperDox
to look for specification files:
```
-spec-dir=/user/api_specs
```

and because the main swagger specification file is within a sub-directory of the common parent `api_specs`, give the
relative path to this using the `-spec-filename` parameter:
```
-spec-filename=spec/swagger.json
```

### Resolving references

You may find that your specification does not load first time if it consists of multiple files
referenced by `$ref:` members. Please refer to the section on [resolving references](/docs/spec-references) if this is the case.

!!!HIGHLIGHT!!!
