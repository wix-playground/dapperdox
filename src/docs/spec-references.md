Title: Resolving specification references
Description: How to correctly resolve references in your OpenAPI specification
Keywords: openAPI, swagger, specification, resolve, references

# Resolving specification references

It is common that a specification that is split over multiple files will fail to load the first time, caused by its `$ref`
members not resolving to the address that DapperDox is serving the files from.

To correct this, you can tell DapperDox to automatically rewrite all specification references so that they correctly
resolve to the running DapperDox instance.

For example, if the above example specification has been written with `http://mydomain.com/swagger-2.0/` as its base URL,
such that the main swagger specification would be found at `http://mydomain.com/swagger-2.0/spec/swagger.json`, and if
DapperDox is serving files from address `http://localhost:3123`, at the directory `/users/api_specs/` and down,
then we need to rewrite all references to `http://mydomain.com/swagger-2.0/` in the specification as `http://localhost:3123/`.

To do this, pass the `-spec-rewrite-url` option to DapperDox:

```
-spec-rewrite-url=http://mydomain.com/swagger-2.0 \
-site-url=http://localhost:3123
```

This rewrites the part of a URL that matches the value given by `-spec-rewrite-url` as the value given by `-site-url`.

> Note that only **absolute** `$ref` references are allowed in a Swagger specification. Relative references in JSON
files generally do not resolve consistently, because the concept of what the reference is relative *to* changes from
file to file.

See [Rewriting specification URLs](/docs/rewrite-spec) for further details.
!!!HIGHLIGHT!!!
