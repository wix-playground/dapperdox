# Getting started

The Getting Started guide shows you how to download and install Swaggerly; how to quickly
get an OpenAPI Swagger 2.0 specification rendered, and some useful configuration options
to tune the presentation.

## Download Swaggerly

[Download the latest release](/) of Swaggerly for you operating system, and unpack it.

```
tar -zxf swaggerly-*.tgz
cd swaggerly
```

## Basic configuration

Swaggerly requires little configuration to get started, just where to find your OpenAPI Swagger 2.0
specifications. Configuration parameters can be specified on the command line, or by setting
environment variables.

For this example, we will specify the Swagger specification location on the command line:
```
./swaggerly -spec-dir=examples/specifications/petstore
```

This will start serving the example <em>petstore</em> specification included in the Swaggerly distribution.

By default, Swaggerly will start serving reference documentation on port 3123, so point your web browser at
[http://localhost:3123](http://localhost:3123)

## Serving your own specification

Point Swaggerly at the directory containing your Swagger specification using the `-spec-dir` configuration parameter,
where Swaggerly will look for the file `swagger.json`. If your specification has a different filename, specify
this with the `-spec-filename` parameter.

If your specification is broken down into multiple files, with some in sub-directories, then you must point
Swaggerly at the parent directory, beneath which it can find all the files it needs. If the main Swagger
specification file exists in a sub-directory of this parent, then configure that using the `-spec-filename` parameter.

For example, consider the following file structure:

```
/user/api_specs/spec/swagger.json
/user/api_specs/operations/get_user.json
/user/api_specs/operations/post_user.json
/user/api_specs/definitions/user_resource.json
```

Since the common parent directory of all these resource files is `api_specs`, this is where you need to tell Swaggerly
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

It is common that a specification that is split over multiple files will fail to load the first time, caused by its `$ref`
members not resolving to the address that Swaggerly is serving the files from.

To correct this, you can tell Swaggerly to automatically rewrite all specification references so that they correctly
resolve to the running Swaggerly instance.

For example, if the above example specification has been written with `http://mydomain.com/swagger-2.0/` as its base URL,
such that the main swagger specification would be found at `http://mydomain.com/swagger-2.0/spec/swagger.json`, and if
Swaggerly is serving files from address `http://localhost:3123`, at the directory `/users/api_specs/` and down,
then we need to rewrite all references to `http://mydomain.com/swagger-2.0/` in the specification as `http://localhost:3123/`.

To do this, pass the `-spec-rewrite-url` option to Swaggerly:

```
-spec-rewrite-url=http://mydomain.com/swagger-2.0 \
-site-url=http://localhost:3123
```

This rewrites the part of a URL that matches the value given by `-spec-rewrite-url` as the value given by `-site-url`.

> Note that only **absolute** `$ref` references are allowed in a Swagger specification. Relative references in JSON
files generally do not resolve consistently, because the concept of what the reference is relative *to* changes from
file to file.

See [Rewriting URLs](/docs/rewrite-urls.html) for further details.
!!!HIGHLIGHT!!!
