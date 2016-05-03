Swaggerly
=========

**BETA release - features being added weekly**.

> Themed documentation generator, server and API explorer for OpenAPI (Swagger) Specifications. Helps you build integrated, browsable reference documentation and guides. For example, the [Companies House Developer Hub](https://developer.companieshouse.gov.uk/api/docs/) built with the alpha version.

## Quickstart

First build swaggerly (this assumes that you have your golang environment configured correctly):
```bash
make
```

Then start up Swaggerly, pointing it to your OpenAPI 2.0 specification file:
```
./swaggerly -spec-dir=<location of OpenAPI 2.0 spec>
```
Swaggerly looks for the file path `spec/swagger.json` at the `-spec-dir` location, and builds reference documentation for the OpenAPI specification it finds. For example, the obligitary *petstore* OpenAPI specification is provided in the `petstore` directory, so
passing parameter `-spec-dir=petstore` will build the petstore documentation.

Swaggerly will default to serving documentation from port 3123 on all interfaces, so you can point your web browser to
either http://0.0.0.0:3123, http://127.0.0.1:3123 or http://localhost:3123.

For an out-of-the-box example, execute the example run script. A description of what this does is given in the section [Swaggerly start up example](#swaggerly-start-up-example), as it makes use of many of the configuration options discussed in this README.

```bash
./run_example.sh
```

*See the section [Why a makefile and not go build?](#why-a-makefile-and-not-go-build) to understand why a makefile is necessary.*

## Guide Contents
- [Next steps](#next-steps)
- [Configuring the address of the server](#specifying-the-address-of-the-server)
- [Specifying an OpenAPI specification](#specifying-an-openapi-specification)
- [Specification requirements](#specification-requirements)
- [Rewriting URLs](#rewriting-urls)
- [The API explorer](#the-api-explorer)
  - [Customising authentication credential capture](#customising-authentication-credential-capture)
    - [apiExplorer methods](#apiexplorer-methods)
  - [Controlling authentication credential passing](#controlling-authentication-credential-passing)
- [Customising the documentation](#customising-the-documentation)
  - [Creating local assets](#creating-local-assets)
  - [Creating authored documentation pages](#creating-authored-documentation-pages)
  - [Customising the 'homepage'](#customising-the-homepage)
- [Versioning](#versioning)
- [Reverse proxying to other resources](#reverse-proxying-through-to-other-resources)
- [Configuration parameters](#configuration-parameters)
- [Swaggerly start up example](#swaggerly-start-up-example)
- [Why a makefile and not go build?](#why-a-makefile-and-not-go-build)

## Next steps
While simply running Swaggerly and pointing it at your swagger specifications will give you some documentation quickly, there
will be a number of things that you will want to configure or change:

1. The URLs picked up from the swagger specifications will probably not match your environment.
2. You will want to supplement the auto-generated resource documentation with your own authored text and guides.
3. The default authentication credential passing may not match your API requirements.

## Configuring the address of the server

Swaggerly will start serving content from http://0.0.0.0:3123. You can change this through the `-bind-addr` configuration
parameter, the format of which being `IP address:port number`, such as `-bind-address=0.0.0.0:3123`.
See [Configuration parameters](#configuration-parameters) for further information on configuring Swaggerly.

## Specifying an OpenAPI specification

Out of the box, Swaggerly will look for the OpenAPI specification `spec/swagger.json` under the directory specified by the
command line option `-spec-dir`. To change this, you can supply the `-spec-filename` option to Swaggerly. For example,
`-spec-filename=spec/swagger.json` does the same as the default.

All JSON specification files found below the `-spec-dir` are served by Swaggerly, maintaining the directory structure.
For the petstore example, which has its OpenAPI specification `swagger.json` stored in a `spec` subdirectory, the 
URL to retrieve the specification is:

```url
http://127.0.0.1:3123/spec/swagger.json
```


See [Configuration parameters](#configuration-parameters) for further information on configuring Swaggerly.

Multiple API specifications are not currently supported, but are on the feature list.

## Specification requirements

### Tags

Swaggerly will try and read read the top level specification object `tags` member, and if it finds one it will only documents
API operations where tags match, and in the order they are listed. This allows you to control what reference documentation gets
presented. In these cases, Swaggerly uses the tag `description` member, or tag `name` member as the API identifier in pages, navigation and URLs.

If tags are not used, Swaggerly falls back to presenting all operations in the OpenAPI specification.

```json
{
  "swagger": "2.0",

  "tags": [
    { 
        "name": "Products",
        "description": "A more verbose description of tag"
    },
    { "name": "Estimates Price" }
  ],
    "paths": {
        "/products": {
            "get": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Read product types"
            },
            "post": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Create product types"
            }
        }
    }
}
```

This incomplete specification example shows how documentation order and filtering is controlled by `tags`. The top level
tags member declares that API operations tagged with `Estimates Price` and `Products` should be built, in that order.
The `description` member of the `Products` tag is used to name all operations grouped by that tag. The name of the
`Estimates Price` tag would be used to name all operations grouped by that tag, as there is no `description` member.

This mechanism for naming and grouping API operations gives you the most control.

However, if tags cannot be used, Swaggerly must still have a title to use for an API path, and will fall back to using
the `summary` member from one of the path operations. This often does not produce the best results, unless
the `summary` members of all operations for a path are set to the same text, as in the example above, but will often not be
the case.

To force the group name of all operations declared for a path, the Swaggerly specific `x-pathName` member may be specified
in the Path Item object. This will always take effect, and will even override the description or name inherited from the top
level `tags` member. However, tags are the most flexible approach.

```json
{
  "swagger": "2.0",

  "tags": [
    { 
        "name": "Products",
        "description": "A more verbose description of tag"
    },
    { "name": "Estimates Price" }
  ],
    "paths": {
        "/products": {
            "x-pathName": "Types of Product",
            "get": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Read product types"
            },
            "post": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Create product types"
            }
        }
    }
}
```

### Response model title

When processing model definitions, Swaggerly needs to know what to call the response schema (or model).
The following snippet shows the API response model `Product`, explicitly named with the `title` member.

```json
"definitions": {
    "Product": {
        "title":"Product resource",
        "properties": {
            "product_id": {
                "type": "string",
                "description": "Unique identifier..."
            }
        }
    }
}
```

Even though `title` is optional in the OpenAPI specification, without it Swaggerly will generate an error:
```
Error: GET /estimates/price references a model definition that does not have a title member.
```


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

There is nothing special about the @ in the above example, it is merely a convention. You can use any expansion syntax you want.

You may pass multiple `-document-rewrite-url` parameters to Swaggerly, to have it replace multiple URLs or placeholders,
particularly useful if you additionally need to configure URLs such as CDN location.

See [Configuration parameters](#configuration-parameters) for further information on configuring Swaggerly.

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
See [Configuration parameters](#configuration-parameters) for further information on configuring Swaggerly.

# The API explorer

The Swaggerly in-page API explorer is similar in function to **swagger-ui**, as it allows user's to try out API calls
from within the reference page, without needing to write any client code.

The Swaggerly in-page API explorer detects when a method is configured as authenticated, and prompts for appropriate
credentials to be supplied as part of the request being explored. These could be one of API key or an OAuth2 access token.

If you have, or are building, a developer site that allows users to register for and manage their own API keys, you may want 
Swaggerly to integrate with that site, so that a user's API keys are automatically available to the explorer once the user has
signed-in. Swaggerly provides a simple Javascript interface via which you can pass API keys to the explorer, through a piece
of custom Javascript.


## Customising authentication credential capture

The `apiExplorer` javascript object provides a method to add API keys to an internal list, and a method to inject those
API keys into the explorer page, so that the user can select the key from a dropdown menu instead of having to type it in.

### apiExplorer methods

| Method | Description |
| ------ | ----------- |
| `apiExplorer.addApiKey(name, key)` | This method adds the named key to the internal list. |
| `apiExplorer.listApiKeys()` | Returns an array of key names. |
| `apiExplorer.getApiKey( name )` | Returns the key associated with name `name`. |
| `apiExplorer.injectApiKeysIntoPage()` | Injects the named API keys into the explorer, building a pulldown menu that can be selected from by the user. |

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

See [Creating local assets](#creating-local-assets) for further information about creating custom assets.

## Customising the documentation
Swaggerly presents two classes of documentation:

1. API reference documentation, derived from Swagger specifications
2. Guides and other authored documentation

Documentation is built from assets, which mostly consist of styles, page templates and template fragments, grouped together
to form a theme. To customise the documentation: 

1. Alternative themes may be used to change the look and feel.
2. Additional assets may be provided to extend the generated documentation, such as guides
3. Parts of a theme may be overridden (advanced topic)

In general, documentation should be written using [Github Flavoured Markdown](https://help.github.com/articles/basic-writing-and-formatting-syntax/), which seamlessly integrates with the reference
documentation generated from the OpenAPI specification.

### Creating local assets

Swaggerly builds documentation for several sets of assets. The primary assets are those which make up the theme being used
for presentation, however Swaggerly will also pick up local assets and serve them along with the reference documentation
it builds from the OpenAPI specification. The local assets directory can be considered equivalent to the `docroot` of a
web server.

Local assets can be images, guides, styles, javascript and *replacement* assets for those provided by the theme.

The directory structure of your local assets must follow a defined directory structure, as Swaggerly needs to understand
what it is serving and whether it is a replacement resource or not. It can do this by matching the assets directory
structure with that provided by the theme:

- `assets/`
    - `static/`
        - `css/` - Local stylesheets
        - `js/` - Local javascript
    - `templates/`
        - `guides/` - Local authored documentation

To have Swaggerly pick up your local assets, pass the `-assets-dir=<directory-path>` option to Swaggerly on start up. See
[Configuration parameters](#configuration-parameters) for further information on configuring Swaggerly.

### Creating authored documentation pages

Authored documentation pages are referred to as *guides*, and have their own directory within an assets structure. Guides may
be authored in HTML as `.tmpl` files, or as [Github Flavoured Markdown](https://help.github.com/articles/basic-writing-and-formatting-syntax/). Writing guides as HTML `.tmpl` files will make those files dependant on the theme in use when they were written,
and therefore not resistant to change. The flexible approach is to use Github Flavoured Markdown.

#### Github Flavoured Markdown

Guides written using [Github Flavoured Markdown](https://help.github.com/articles/basic-writing-and-formatting-syntax/) (GFM)
have a file extension of `.md` and are stored within directory `assets/templates/guides/` of your [local assets](#creating-local-assets). You can organise your files in subdirectories within the `/guides/` directory.

On startup, Swaggerly will locate and build pages for all of your guides, maintaining the directory 
structure if finds below the `/guides/` directory.

For example, the Swaggerly assets example `examples/markdown/assets/templates/` contains two guides:

1. `guides/markdown.md`
2. `guides/level2/markdown2.md`

Passing Swaggerly the `-assets-dir=<Swaggerly-source-directory>/examples/markdown/assets` will build these 
two guides, making them available as http://127.0.0.1:3123/guides/markdown and 
http://127.0.0.1:3123/guides/level2/markdown2

The navigation rendered at the side of the page will show two navigation entries:

- level2
  - markdown2
- markdown

By default, the side navigation will reproduce the directory structure beneath the `guides/` directory.
As the navigation cannot be more than two levels deep, this restricts the depth of your directory structure.

If you need a more elaborate directory structure, or have a file nameing convention that does not lend itself
to navigation titles, you can take control of the side navigation through [metadata](#controlling-guide-behaviour-with-meta-data).

**GFM support is a new feature, and so guides created using GFM are not currently styled correctly. Standard GFM HTML is generated which does not use the appropriate theme CSS. This being addressed in [issue #1](https://github.com/zxchris/swaggerly/issues/1)**

#### Controlling guide behaviour with metadata

Swaggerly allows the integration of guides to be controlled with some simple metadata. This metadata is added
to the beginning of GFM files as a block of lines containing `key: value` pairs. If present, metadata ends at
the first blank line.

Through the metadata, you can control the side navigation hierarchy, grouping and page naming.

For example, the metadata contained within the example `examples/metadata/assets/templates/guides/markdown.md` is:

```http
Navigation: Examples/A markdown example
SortOrder: 210
Note: This top section is just MetaData and gets stripped to the first blank line.

This page was written using Git Flavoured Markdown - With metadata
==================================================================
```

Whereas the example `examples/metadata/assets/templates/guides/level2/level3/markdown2.md` which is *three*
directory levels deep, contains navigation metadata of:

```http
Navigation: Overview/Another example
SortOrder: 110
```

This builds a page side navigation two levels deep:

- Overview
  - Another example
- Examples
  - A markdown example

By using metadata, the side navigation wording and structure is divorced from the underlying file naming
convention, structure and depth.

The ordering of pages within the page side navigation is also controllable with metadata, as described in [SortOrder](#sortorder) below.

#### Supported metadata

The following metadata is recognised by Swaggerly. All other metadata entries will be ignored.

##### Navigation

The `Navigation` metadata entry describes how the page is integrated into the site navigation. The navigation value is a
path that defines the page placement in the navigation tree. With the default theme, guides are placed *before* the
reference documentation in the navigation.

For example, a page containing the metadata `Navigation: Examples/A markdown example` creates a navigation section called
*Examples* and places that page beneath it, with the description *A markdown example*.

##### SortOrder

The order in which guides are listed in the page side navigation is controlled with `SortOrder` metadata.
`SortOrder` can take any alphanumeric string, but may be clearer if numeric only values are used.
Each page is assigned the sort order defined by the metadata, otherwise it is assigned the relative URI path
(`/guides/.....`) as its sort order.

Where pages are grouped by a parent section, the parent section receives the lowest sort order of its
children, unless it is a page in its own right.
Assigning a numeric sort order range per section makes it easy to manage the ordering of sections, 
and the pages within a section. This can be illustrated by the following tree, where pages have been give 
numeric `SortOrder` metadata entries, assigning blocks of 100 per section:

- 100 Overview
  - 100 - A page
  - 150 - Another page
- 210 Getting Started
  - 210 - Getting started page one
  - 250 - Getting started page two
- 300 Examples
  - 300 - Examples page one
  - 350 - Examples page two
- 400 Top level page one
- 420 Top level page two


### Customising the 'homepage'

By default, the homepage that Swaggerly presents is an API reference summary. You can create your own
homepage by providing your own `assets/templates/index.tmpl` or `assets/templates/index.md` - with the
recommendation being that you use markdown instead of HTML, for the reasons described in 
[Creating authored documentation pages](#creating-authored-documentation-pages).

An example of this is demonstrated by the metadata example, which provides its own custom `index.md` file:

`examples/metadata/assets/templates/index.md`

To run this example, pass Swaggerly the option 
`-assets-dir=<Swaggerly-source-directory>/examples/metadata/assets`

The API reference summary will always be available at the `/reference` endpoint.


# Versioning

To be completed.

## Reverse proxying through to other resources

To create an integrated developer hub. Such resources could be:

1. API key management tools
2. Forum

**Coming soon!**

## Dependencies

### `go-swagger`

Swaggerly depends on a fork of [go-swagger](https://github.com/zxchris/go-swagger) as its specification parser. This
fork adds missing support for complex object response (arrays of objects), and the Swaggerly specific versioning scheme.
Versioning is currently implemented on the `feature/swaggerly-versioning-extension` branch, and it is on this branch that
Swaggerly depends.

## Why a makefile and not go build?

As described in [Dependencies](#dependencies), Swaggerly requires a particular branch of `go-swagger` and `go get`, by default,
will checkout the master branch. The supplied makefile ensures the correct `go-swagger` branch is checked out and built.

If you wish to perform the build manually, the following steps should be carried out (as shown in the Makefile):

```bash
go get github.com/go-swagger/go-swagger
cd ../../go-swagger/go-swagger
git checkout 4459770
cd -
go get github.com/zxchris/go-swagger
cd ../go-swagger
git checkout feature/swaggerly-versioning-extension
cd ../swaggerly
go get && go build
```

## Configuration parameters

| Option | Environment variable | Description |
| ------ | -------------------- | ----------- |
| `-assets-dir` | `ASSETS_DIR` | Assets to serve. Effectively the document root. |
| `-bind-addr` | `BIND_ADDR` | Bind address. |
| `-default-assets-dir` | `DEFAULT_ASSETS_DIR` | Default assets. |
| `-document-rewrite-url` | `DOCUMENT_REWRITE_URL` | Specify a URL that is to be rewritten in the documentation. May be multiply defined. Format is from=to. This is applied to assets, not to OpenAPI specification generated text. |
| `-log-level` | `LOGLEVEL` | Log level (info, debug, trace) |
| `-site-url` | `SITE_URL` | Public URL of the documentation service. |
| `-spec-dir` | `SPEC_DIR` | OpenAPI specification (swagger) directory. |
| `-spec-filename` | `SPEC_FILENAME` | The filename of the OpenAPI specification file within the spec-dir. Defaults to spec/swagger.json |
| `-spec-rewrite-url` | `SPEC_REWRITE_URL` | The URLs in the OpenAPI specifications to be rewritten as site-url, or to the `to` URL if the value given is of the form from=to. Applies to OpenAPI specification text, not asset files. |
| `-theme` | `THEME` | Name of the theme to render documentation. |
| `-themes-dir` | `THEMES_DIR` | Directory containing installed themes. |

## Swaggerly start up example

The following command line will start Swaggerly serving the petshop OpenAPI specification, rewriting the API host URL
of the petstore specification from api.uber.com to API.UBER.COM, and picking up the local assets `examples/markdown/assets` 
which brings in the two GFM example guides, rewriting `www.google.com` within them to `www.google.co.uk`.

This start up script can be found as `run_example.sh` in the swaggerly source directory.

```bash
./swaggerly \
    -spec-dir=petstore \
    -bind-addr=0.0.0.0:3123 \
    -spec-rewrite-url=api.uber.com=API.UBER.COM \
    -document-rewrite-url=www.google.com=www.google.co.uk \
    -site-url=http://127.0.0.1:3123 \
    -assets-dir=./examples/markdown/assets \
    -log-level=trace
```
