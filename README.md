Swaggerly
=========

> Themed documentation generator, server and API explorer for OpenAPI (Swagger) Specifications. Helps you build integrated, browsable reference documentation and guides.

## Quickstart

First build swaggerly (this assumes that you have your golang environment configured correctly):
```bash
go get && go build
```

Then start up Swaggerly, pointing it to your OpenAPI 2.0 specification file:
```
./swaggerly -spec-dir=<location of OpenAPI 2.0 spec>
```
Swaggerly looks for the file path `spec/swagger.json` at the `-spec-dir` location, and builds documentation for the OpenAPI specification it finds. For example, the obligitary *petstore* OpenAPI specification is provided in the `petstore` directory, so
passing parameter `-spec-dir=petstore` will build the petstore documentation.

Swaggerly will default to serving documentation from port 3123 on all interfaces, so you can point your web browser to
either http://0.0.0.0:3123, http://127.0.0.1:3123 or http://localhost:3123.

## Guide Contents
- [Next steps](#next-steps)
- [Specifying an OpenAPI specification](#specifying-an-openapi-specification)
- [Configuring the address of the server](#specifying-the-address-of-the-server)
- [Rewriting URLs](#rewriting-urls)
- [Customising the documentation](#customising-the-documentation)
  - [Themes](#theme-assets)
  - [Custom documentation pages](#adding-custom-documentation)
  - [Adding nagivation](#adding-navigation)
- [The API explorer](#the-api-explorer)
  - [Customising authentication credential capture](#customising-authentication-credential-capture)
    - [apiExplorer methods](#apiexplorer-methods)
  - [Controlling authentication credential passing](#controlling-authentication-credential-passing)
- [Versioning](#versioning)
- [Reverse proxying to other resources](#reverse-proxying-through-to-other-resources)

## Next steps
While simply running Swaggerly and pointing it at your swagger specifications will give you some documentation quickly, there
will be a number of things that you will want to configure or change:

1. The URLs picked up from the swagger specifications will probably not match your environment.
2. You will want to supplement the auto-generated resource documentation with your own authored text and guides.
3. The default authentication credential passing may not match your API requirements.

## Specifying an OpenAPI specification

Out of the box, Swaggerly will look for the OpenAPI specification `spec/swagger.json` under the directory specified by the
command line option `-spec-dir`. To change this, you can supply the `-spec-filename` option to Swaggerly. For example,
`-spec-filename=spec/swagger.json` does the same as the default.

Multiple API specifications are not currently supported, but are on the feature list.


## Configuring the address of the server

Swaggerly will start serving content from http://0.0.0.0:3123. You can change this through the `-bind-addr` configuration
parameter, the format of which being `IP address:port number`, such as `-bind-address=0.0.0.0:3123`.

## Rewriting URLs
### Documentation URLs
The swagger specification often does not contain API or resource URLs that are correct for the environment being documented.
For example, the swagger specifications may contain the production URLs, which are not appropriate when the specification and
documentation is being served up in a development or test environment.

Swaggerly allows you to rewrite URLs on the fly, so that they match the environment they are being served from. To do this,
you specify the URL pattern that should be rewritten *from* and *to*, by passing the `-document-rewrite-url` configuration
parameter. The parameter takes a `from=to` pair, such as

```
-document-rewrite-url=domain.name.from=domain.name.to
```

You may also choose to use placeholders in your documentation, instead of real URLs, so that you replace the placeholder with
a valid URL:

```html
<a href="@MY_DOMAIN/some/document.html">Some document</a>
```

which would be rewritten with:

```bash
-document-rewrite-url=@MY_DOMAIN=http://www.mysite.com
```

You may pass multiple `-document-rewrite-url` parameters to Swaggerly, to have it replace multiple URLs or placeholders,
particularly useful if you additionally need to configure URLs such as CDN location.

### Specification URLs

If your swagger specification is split over multiple files, and therefore contain absolute `$ref:` object references, these
references will not be followed correctly unless they resolve to the running Swaggerly instance serving the files.

For example, if the swagger specification uses the absolute references of `http://mydomain.com/swagger-2.0/....`, and
Swagger is serving content from `http://localhost:3123`, then the additional configuration parameters to pass to Swaggerly
to correct this would be:

```
-spec-rewrite-url=http://mydomain.com/swagger-2.0 \
-site-url=http://localhost:3123 
```

Sometimes you will want to map a specification URL to one that is not the `site url`, for example changing the URL that the
API is served from to be live instead of test. To do this, supply `-spec-rewrite-url` with a `from=to` pair. 

```
-spec-rewrite-url=http://api.test.domain.com=http://api.live.domain.com
```

You may pass multiple `-spec-rewrite-url` parameters to Swaggerly, to have it replace multiple URLs or placeholders.

## Customising the documentation
Swaggerly presents two classes of documentation:

1. API reference documentation, derived from Swagger specifications
2. Guides and other authored documentation

Documentation is built from assets, which mostly consist of styles, page templates and template fragments, grouped together
to form a theme. To customise the documentation: 

1. Alternative themes may be used
2. Parts of a theme may be overridden
3. Additional assets may be provided to extend the generated documentation, such as guides

First an explanation of **assets** is required, and an introduction to its directory structure.

## Theme assets

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

### Using an alternative theme
```
   -themes-dir=<installation directory of additional themes>
   -theme=<name of theme to use>
```

### Adding custom documentation

One of the main features of Swaggerly is its ability to serve specification generated reference documentation along with
authored guides and reference notes, allowing the production of a complete unified documentation suite for an API, or set of
APIs.

To do this, Swaggerly would be configured to pick up and serve additional documentation, written in HTML or Github Flavoured 
Markdown (styling support for which is currently incomplete). These documents are grouped under an application-specific `assets` directory.
The structure of this assets directory is important, as it also allows you to override assets that are normally supplied
by the chosen theme, and so must follow the structure set out in the section called [Theme assets](#theme-assets).

Custom documentation introduces a new directory into the assets structure, that of `assets/templates/guides/`.
Under this directory is placed all the authored content you wish to publish. The content may be a `.tmpl` HTML template or
a `.md` Github Flavoured Markdown file, and will be rendered within the common `theme name/templates/layout.tmpl` page
template.

Swaggerly will serve this content from the `/guides/` URL path, for example `http://localhost:3123/guides/my_example_page`.


#### Adding navigation

Swaggerly will automatically build the necessary navigation for the swagger spec generated content, but it cannot do this
for authored content within the `assets/templates/guides` directory. It leaves you to manage this appropriately. This is
straightforward to do, as you can override the theme `templates/fragments/sidenav.tmpl` and introduce any navigation links you require. Simply copy this file from the theme directory into your application-specific assets as `assets/templates/fragments/sidenav.tmpl` and amend as appropraite.

**(TODO: Investigate whether directives could be added to /guides/* files, giving navigation hints for auto-building of navigation links, which may be sufficient in many cases)**

### Overriding theme assets

Any theme assets may be overridden, as shown in the previous section with `sidenav.tmpl`. This allows you to modify any of
the presentation to suit your local requirements.


### Customisation of auto-generated pages

Every reference page, auto-generated from the swagger specification, can be customised. These pages fall into three categories:

1. API pages. These list and document all the methods available for an API endpoint.
2. Method pages. These document a particular method, acting on an API endpoint.
2. Resource pages. These document the resource that is submitted to or returned from an API endpoint, depending on method.

There are three theme templates that provide the presentation structure for each of these pages, and can therefore be
overridden if you wish to affect a documentation change consistently across all instances of a particular page type.
Alternatively, you may override a page template on an API, method and resource basis.

*The consequence of overriding theme templates is that you cannot change themes without changing your overridden pages. A
Git Flavoured Markdown approach is being developed which will make theme template overriding unnecessary for a majority of
cases.*

#### API pages

An API page will have a URL formed with the following pattern:

```
/docs/<api-name>
```

To override the API page for an endpoint called `example-api`, when using the default theme:

1. Copy `assets/themes/default/templates/default-api.tmpl` to `local-assets/assets/templates/docs/example-api.tmpl`
2. Edit `example-api.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for the API `example-api`, it will use your overridden template.
		
#### Method pages

A method page will have a URL formed with the following pattern:

```
/docs/<api-name>/<method-name>
```

where `method-name` will be `GET`, `POST` etc

To override the `GET` method page for an endpoint called `example-api`, when using the default theme:

1. Copy `assets/themes/default/templates/default-method.tmpl` to `local-assets/assets/templates/docs/example-api/get.tmpl`
2. Edit `get.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for that API method, it will use your overridden template.

#### Resource pages

An API page will have a URL formed with the following pattern:

```
/resources/<resource-name>
```

There can only be one resource of a given name.

To override the resource page for an example resource `example-resource`, when using the default theme:

1. Copy `assets/themes/default/templates/default-resource.tmpl` to `local-assets/assets/templates/resources/example-resource.tmpl`
2. Edit `example-resource.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for this resource, it will use your overridden template.


# The API explorer

The Swaggerly in-page API explorer is similar in function to **swagger-ui**, as it allows users to try out API calls
from within the reference page, without needing to write any client code.

The Swaggerly in-page API explorer detects when a method is configured as authenticated, and prompts for appropriate
credentials to be supplied as part of the request being explored. These could be one of API key or an OAuth2 access token.

If you have, or are building, a developer site that allows users to regiser for and manage their own API keys, you may want 
Swaggerly to integrate with that site, so that a users API keys are automatically available to the explorer once the user has
signed-in. Swaggerly provides a simple Javascript interface via which you can pass API keys to the explorer, through a piece
of custom Javascript.


## Customising authentication credential capture

The `apiExplorer` javascript object provides a method to add API keys to an internal list, and a method to inject those
API keys into the explorer page, so that the user can select the key from a dropdown menu instead of having to type it in.

### apiExplorer methods

`apiExplorer.addApiKey(name, key)`

This method adds the named key to the internal list.


`apiExplorer.listApiKeys()`

Returns an array of key names.


`apiExplorer.getApiKey( name )`

Returns the key associated with name `name`.


`apiExplorer.injectApiKeysIntoPage()`

Injects the named API keys into the explorer, building a pulldown menu that can be selected from by the user.


It would be usual for the template fragment `assets/themes/default/templates/fragments/scripts.tmpl` to be overridden, so
that it can build a list of valid API keys to be used and inject them into the explorer page. Generally the keys would be
fetched from some ajax endpoint that you provide, once the user as gone though some sign-in process.

The supplied example, `examples/apikey_injection/assets/templates/fragments/scripts.tmpl` demonstrates the addition of an
API key (hardcoded for the benefit of this example), and injects the list of one into the explorer page.

To run this example, Swaggerly needs to be told about the example assets directory for it to pick up the override. 
This is achieved through the configuration parameter `-assets-dir`, passed to swaggerly when starting: 
`-assets-dir=examples/apikey_injection/assets`.


## Controlling authentication credential passing

By default, Swaggerly will automatically attach the API key if supplied, to the request URL as a `key=` query parameter.
This behaviour can be customised to satisfy the authentication requirements of your API.

The template fragment `assets/themes/default/templates/fragments/scripts.tmpl`, which is included at the end of the common
page template `layout.tmpl` contains the following javascript snippet:

```javascript
<!-- Additional scripts to be loaded at end of page -->
<!-- This should be overridden to take control of the authorisation process (adding keys to the explorer request). -->
<script>
    $(document).ready(function(){
        // Register callback to add authorisation parameters to request before it is sent
        apiExplorer.setBeforeSendCallback( function( request ) {
            var apiKey = apiExplorer.readApiKey(); // Read API key from explorer input
            request.params = {key: apiKey};
        });
    });
</script>
```

The above snippet registers a callback with the `apiExplorer` object which gets invoked while the explorer is building the 
request to send to the API. This callback will be receive an empty object which has two properties that can be set as needed,
`request.headers` - items that are sent as request HTTP headers,  and `request.params` - items that are sent as query parameters:

```javascript
{
    headers: {},
    params: {}
}
```

Both the `headers` and `params` objects contain zero or more name/value pairs:

```javascript
{
    key1: value,
    ..
    ..
    key_n: value_n,
}
```

For example:
```javascript
request.headers = { header: "value" };
request.headers = { header1: "value1", header2: "value2" }
```

To put this into practice, if you wanted to change the authentication credential passing mechanism to instead supply the API key
as an Authorization header, then create a `scripts.tmpl` within your own assets directory to override this. For example, the
Swaggerly example file `examples/apikey_injection/assets/templates/fragments/scripts.tmpl` passes the API key in the 
Authorization header using BASIC authentication:

```javascript
$(document).ready(function(){
    // ... other code cut from here ...

    // Register callback to add authorisation parameters to request before it is sent
    apiExplorer.setBeforeSendCallback( function( request ) {
        var apiKey = apiExplorer.readApiKey(); // Read API key from explorer input
        request.headers = {Authorization:"Basic " + btoa(apiKey + ":")};
    });
});
```

Swaggerly then needs to be told about this local assets directory for it to pick up the override, which is achieved through
the configuration parameter `-assets-dir`, passed to swaggerly when starting. For example, to pick up the example above, use
`-assets-dir=examples/apikey_injection/assets`.

## Versioning

To be completed

## Reverse proxying through to other resources

To create an integrated developer hub. Such resources could be:

1. API key management tools
2. Forum

To be completed

## Dependencies

### `go-swagger`

Swaggerly depends on a fork of [go-swagger](https://github.com/zxchris/go-swagger) as it's specification parser. This
fork adds missing support for complex object response (arrays of objects), and the Swaggerly specific versioning scheme.


