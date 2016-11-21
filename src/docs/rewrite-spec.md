Title: Rewriting specification URLs
Description: How to rewrite URLs appearing in OpenAPI specifications
Keywords: Rewrite, URL, OpenAPI, Swagger, specification

# Rewriting specification URLs

If your swagger specification is split over multiple files, and therefore contain absolute `$ref:` object
references, these references will not be followed correctly unless they resolve to the running DapperDox instance serving
the files.

For example, if the swagger specification uses the absolute references of `http://mydomain.com/swagger-2.0/....`, and
Swagger is serving content from `http://localhost:3123`, then the additional configuration parameters to pass to DapperDox
to correct this would be:

```
-spec-rewrite-url=http://mydomain.com/swagger-2.0 \
-site-url=http://localhost:3123
```

Multiple `-spec-rewrite-url` options may be given if you have several URLs you need to rewrite, perhaps in the case
where you have embedded links to external documentation. In these scenarios rewriting to a single site-url is insufficient, 
and you will want to use the alternative form of the configuration option, which has a `{from}` and `{to}` component:

```
-spec-rewrite-url=http://mydomain.com/swagger-2.0=http://localhost:3123
```

See [Configuration guide](/docs/configuration-guide) for further information on configuring DapperDox.

!!!HIGHLIGHT!!!
