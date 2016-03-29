Swaggerly
=========

## Quickstart

First build swaggerly (this assumes that you have your golang environment configured correctly):
```bash
go get && go build
```

Then start up Swaggerly, pointing it to your swagger specifications:
```
./swaggerly -bind-addr=0.0.0.0:3128 -swagger-dir=<location of swagger 2.0 specifications>
```

Swaggerly looks for the file path `spec/swagger.json` at the `-swagger-dir` location, and builds documentation for the swagger specification it finds. It then starts serving this documentation from IP address 0.0.0.0 port 3123.

Point your browser to `http://0.0.0.0:3123` or `http://127.0.0.1:3123` or `http://localhost:3123`

### Next steps
Where simply running Swaggerly and pointing it at your swagger specifications will give you some documentation quickly, there
will be a number of things that you will want to configure or change:
1. The URLs picked up from the swagger specifications will probably not match your environment.
2. You will want to supplement the auto-generated resource documentation with your own authored text and guides.

## Rewriting URLs in the documentation
The swagger specification often does not contain API or resource URLs that are correct for the environment being documented.
For example, the swagger specifications may contain the production URLs, which are not appropriate when the specification and
documentation is being served up in a development or test environment.

Swaggerly allows you to rewrite URLs on the fly, so that they match the environment they are being served from. To do this,
you specify the URL pattern that should be rewritten *from*, by passing the `-rewrite-url` configuration parameter, along with
a `site-url` specifyig the URL ssttern URLs hould be rewritten *to*:

For example, if the swagger specification uses the URL root `http://mydomain.com/swagger-2.0/....` which should be changed to
`http://localhost:3123`, then the additional configuration parameters to pass to Swaggerly are:

```
-rewrite-url=http://mydomain.com/swagger-2.0 \
-site-url=http://localhost:3123 \
```

## Creating your own documentation pages
Swaggerly presents two classes of documentation:

1. API reference documentation, derived from Swagger specifications
2. Guides and other authored documentation.

## Assets
Swaggerly documentation is built from assets, which mostly consist of page templates and template fragments, grouped together to form a theme.

A typical theme assets structure is:

- `theme name/`
    - `static/`
        - `css/` - Theme specific stylesheets
        - `js/` - Theme specific javascript
    - `templates/`
        - `layout.tmpl` - The common page template
        - `default-api.tmpl` - The default API page template
        - `default-method.tmpl` - The default API method page template
        - `default-resource.tmpl` - The default API resource page template
        - `fragments/` - Common HTML fragments used across pages
            - `docs/` - Fragments such as the page side navigation and authorisation details
            - `explorer/` - API explorer fragments, such as input fields

## Using an alternative theme
```
   -themes-dir=<installation directory of additional themes>
   -theme=<name of theme to use>
   -default-assets-dir=<locaton of swaggerly default assets> (usually swaggerly_install_directory/assets)
```
