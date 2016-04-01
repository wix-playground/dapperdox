Swaggerly
=========

## Quickstart

First build swaggerly (this assumes that you have your golang environment configured correctly):
```bash
go get && go build
```

Then start up Swaggerly, pointing it to your swagger specifications:
```
./swaggerly -bind-addr=0.0.0.0:3128 -swagger-dir=<location of swagger 2.0 spec>
```

Swaggerly looks for the file path `spec/swagger.json` at the `-swagger-dir` location, and builds documentation for the swagger specification it finds. For example, the obligitary *petstore* swagger specification is provided in the `petstore` directory, so
passing parameter `-swagger-dir=petstore` will build the petstore documentation.

**TODO: Support multiple API specifications**

The `-bind-addr` parameter configures the address and port number that Swaggerly will serve the documentation from, so in
this case you can point your browser to `http://0.0.0.0:3123`, `http://127.0.0.1:3123` or `http://localhost:3123`

### Next steps
While simply running Swaggerly and pointing it at your swagger specifications will give you some documentation quickly, there
will be a number of things that you will want to configure or change:

1. The URLs picked up from the swagger specifications will probably not match your environment.
2. You will want to supplement the auto-generated resource documentation with your own authored text and guides.
3. The default authentication credential passing may not match your API requirements.

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
-site-url=http://localhost:3123 \
-default-assets-dir=<swaggerly install directory>/assets
```

**Swaggerly does not currently does not translate the swagger Host: member.**

#### NOTES ####

- Swagger - `Host:` member needs setting to API address. `-api-host-url` configuration. **This will not work for multiple api specifications.**
- Swagger - `$ref:` members need setting to Swaggerly address. `-swagger-rewrite-urls` instead of `-rewrite-url` (should accept a list)
- Documentation - Various URLs needed for pages. `-document-rewrite-urls` (a list of key=value,) **done**


## Customising the documentation
Swaggerly presents two classes of documentation:

1. API reference documentation, derived from Swagger specifications
2. Guides and other authored documentation

Documentation is built from assets, which mostly consist of styles, page templates and template fragments, grouped together
to form a theme. To customise the documentation: 

1. Alternative themes may be used
2. Parts of a theme may be overridden
3. Additional assets may be provided to extend the generated documentation, such as guides

First, and explanation of **assets** is required, and an introduction to its directory structure.

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
Markdown (support for which is currently incomplete). These documents are grouped under an application-specific `assets` directory.
The structure of this assets directory is important, as it also allows you to override assets that are normally supplied
by the chosen theme, and so must follow the structure set out in **Theme assets** above.

Custom documentation introduces a new directory to the assets structure: `assets/templates/guides/`.
Under this directory is placed all the authored content you wish to publish. The content may be a `.tmpl` HTML template or
a `.md` Github Flavoured Markdown file, and will be rendered within the common `theme name/templates/layout.tmpl` page
template.

Swaggerly will serve this content from the `/guides/` URL path, for example `http://localhost:3128/guides/my_example_page`.


#### Adding navigation

Swaggerly will automatically build the necessary navigation for the swagger spec generated content, but it cannot do this
for authored content within the `assets/templates/guides` directory. It leaves you to manage this appropriately. This is
straightforward to do, as you can override the theme `templates/fragments/sidenav.tmpl` and introduce any navigation links you require. Simply copy this file from the theme directory into your application-specific assets as `assets/templates/fragments/sidenav.tmpl` and amend as appropraite.

**(TODO: Investigate whether directives could be added to /guides/* files, giving navigation hints for auto-building of navigation links, which may be sufficient in many cases)**

### Overriding theme assets

Any theme assets may be overridden, as shown in the previous section with `sidenav.tmpl`. This allows you to modify 


### Customisation of auto-generated pages

Each of the pages auto-generated from the swagger specification may be customised. These pages fall into three categories:

1. API pages. These list and document all the methods available for an API endpoint.
2. Method pages. These document a particular method, acting on an API endpoint.
2. Resource pages. These document the resource that is submitted to or returned from an API endpoint, depending on method.

There are three theme templates that provide the presentation structure for each of these pages, and can therefore be
overridden if you wish to affect a documentation change consistently across all instances of a particular page type.
Alternatively, you may override a page template on an API, method and resource basis.

**TODO: complete this**

# The API explorer

## Controlling authentication credential capture

The Swaggerly in-page API explorer detects when a method is configured as authenticated, and prompts for appropriate
credentials to be supplied as part of the request being explored. These could be one of API key or an OAuth2 access token.

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
fetched from some ajax endpoint, having previously gone through user authentication.

The supplied example, `examples/apikey_injection/assets/templates/fragments/scripts.tmpl` demonstrates the addition of an
API key (hardcoded for the benefit of this example), and injects the list of one into the explorer page.

To run this example, Swaggerly needs to be told about the example assets directory for it to pick up the override. 
This is achieved through the configuration parameter `-assets-dir`, passed to swaggerly when starting: 
`-assets-dir=examples/apikey_injection/assets`.


## Controlling authentication credential passing

By default, Swaggerly will automatically attach the API key, if supplied, to the request URL as a `key=` query parameter.
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

This snippet registers a callback with the `apiExplorer` object which is invoked while the explorer is building the request
to send to the API. This callback is passed an empty object which has two properties that can be set as needed, `request.headers` and `request.params`:

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

