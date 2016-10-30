# Rewriting URLs

## Specification URLs

If your swagger specification is split over multiple files, and therefore contain absolute `$ref:` object
references, these references will not be followed correctly unless they resolve to the running Swaggerly instance serving
the files.

For example, if the swagger specification uses the absolute references of `http://mydomain.com/swagger-2.0/....`, and
Swagger is serving content from `http://localhost:3123`, then the additional configuration parameters to pass to Swaggerly
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

See [Configuration guide](/docs/configuration-guide) for further information on configuring Swaggerly.

## Documentation URLs
The authored documentation you are combining with your swagger specifications often will not contain URLs
that are correct for the environment being documented.

For example, the specification guides may contain the production URLs, which are not appropriate when the documentation
is being served up in a development or test environment.

Swaggerly allows you to rewrite these documentation URLs on the fly, so that they match the environment they are being
served from. To do this, you specify the URL pattern that should be rewritten *from* and *to*, by passing the
`-document-rewrite-url` configuration parameter. The parameter takes a `{from}={to}` pair, such as:

```
-document-rewrite-url=domain.name.from=domain.name.to
```

You may also choose to use placeholders in your documentation, instead of real URLs, so that you replace the placeholder with
a valid URL:

```html
<a href="MY_DOMAIN/some/document">Some link</a>
```

which would be rewritten with:

```
-document-rewrite-url=MY_DOMAIN=http://www.mysite.com
```

There is nothing special about `MY_DOMAIN` in the above example, it is merely a convention. You can use any expansion text you like.

You may pass multiple `-document-rewrite-url` parameters to Swaggerly, to have it replace multiple URLs or placeholders,
particularly useful if you additionally need to configure URLs such as CDN location.

See [Configuration guide](/docs/configuration-guide) for further information on configuring Swaggerly.

!!!HIGHLIGHT!!!
