# Content overlays

Additional content can be added to any Swaggerly generated reference page by providing overlay files.
These pages are authored in Github Flavoured Markdown (GFM) and contain special markdown references that
target particular sections within API, Method or Resource pages.

Additional directories are added to your `assets` directory to accomplish this. As Swaggerly can consume and serve
multiple OpenAPI specifications, each is given its own section within Swaggerly, allowing you to provide guides and
overlay documentation appropriate to an OpenAPI specification. 

See [Controlling method names](/docs/author-method-names) for a discussion on what an *operation* name is, and
how it differs from an HTTP method name.

For example, the following GFM file adds additional content to the *request* and *response* sections
of the **Find pet by ID** `GET` method for the example Petstore OpenAPI specification, where `{local_assets}` is the directory
pointed to by the `-assets-dir=` configuration parameter:

```
> cat {swaggerly-source-directory}/examples/overlay/assets/sections/swagger-petstore/templates/reference/everything-about-your-pets/get-pet-by-id.md
```

```gfm
Overlay: true

[[banner]]
> This is a banner. Content overlaid onto this method page comes from
`assets/sections/swagger-petstore/templates/reference/everything-about-your-pets/get-pet-by-id.md`
and is therefore applied to **all** GET methods defined by the Petstore specification.

[[description]]
This is some overlaid *description* text. Defined by file
`assets/sections/swagger-petstore/templates/reference/everything-about-your-pets/get-pet-by-id.md`

[[response]]
Here is some overlaid *response* text. Defined by file
`assets/sections/swagger-petstore/templates/reference/everything-about-your-pets/get-pet-by-id.md`
```

For a GFM page to be treated as an overlay, it must contain the metadata `Overlay: true` at the start
of the file (see [Supported metadata](#supported-metadata)).

## Prioritising overlays

There are three ways to overlay a reference page, globally, per-specification or on a page-by-page basis.
To find an overlay for a page, Swaggerly follows the file patterns defined in the table below
**in the order in which they are listed**, using the first one it finds.
For example `sections/{specification-ID}/templates/reference/{api-group}.md` takes precidence over `sections/{specification-ID}/templates/reference/api.md`.

In the references below, the `operation-ID` for a method is determined by taking the
[kebab case](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles) of the
first member in the following list that yields a value from an operations's specification:

1. `operationID`
2. `x-opererationName`
3. `summary`
4. method name

where method name (and `method-name` below) is the actual HTTP method name of the operation, such as `get`, `post` and `put`.
Therefore, it is possible for `operation-ID` and `method-name` to be the same value.

`api-group` is calculated in two ways, depending on whether tagging is in use or not, and is also kebab-cased. If tagging is in use, then `api-group` 
takes the tag's `description`, falling back to the tag `name` if there no description. If tagging is not in use, then the 
operation's `x-pathName` member is used, otherwise it takes the operation's `summary`. 

### Overlay file priority

This table defines the overlay file lookup priority for API, method and resource pages:

| Reference page | Overlay filename | Description |
| -------------- | ---------------- | ----------- |
| API      | `sections/{specification-ID}/templates/reference/{api-group}.md`               | Overlay applied to a specific API page. |
| API      | `sections/{specification-ID}/templates/reference/api.md`                      | Overlay applied to all API pages for the named specification. |
| API      | `templates/reference/api.md`                               | Overlay applied to all API pages. |
| Method   | `sections/{specification-ID}/templates/reference/{api-group}/{operation-ID}.md` | Overlay applied to a specific method page of a specific API of a specific openAPI specification. |
| Method   | `sections/{specification-ID}/templates/reference/{api-group}/{method-name}.md` | Overlay applied to all method pages with this HTTP method name in the named openAPI specification. |
| Method   | `sections/{specification-ID}/templates/reference/{api-group}/method.md`        | Overlay applied to all method pages of a specific API. |
| Method   | `sections/{specification-ID}/templates/reference/{operation-ID}.md`         | Overlay applied to all method pages for {operation-ID} across all APIs in a specification. |
| Method   | `sections/{specification-ID}/templates/reference/{method-name}.md`           | Overlay applied to all method pages with this HTTP method name, across all APIs in the named specification. |
| Method   | `sections/{specification-ID}/templates/reference/method.md`                   | Overlay applied to all method pages of all APIs in the named specification. |
| Method   | `templates/reference/{method-name}.md`                  | Overlay applied to the all methods with this HTTP method name for all APIs across all specifications.  |
| Method   | `templates/reference/method.md`                            | Overlay applied to all method pages of all APIs across all specifications.  |
| Resource | `sections/{specification-ID}/templates/resource/{resource-name}.md`           | Overlay applied to a specific resource page of a specific API.  |
| Resource | `sections/{specification-ID}/templates/resource/resource.md`                  | Overlay applied to all resource pages of a specific API.  |
| Resource | `templates/resource/resource.md`                           | Overlay applied to all resource pages of all APIs across all specifications.  |

This can be visualised as a directory tree (though precedence is not maintained in this representation):

- `{local_assets}/`
    - `templates/`
        - `guides/` - Authored documentation presented when not viewing an OpenAPI specification section.
        - `reference/` - Custom overlay GFM content.
          - `api.md` - Overlay applied to all API pages.
          - `{operation-ID}.md` - Overlay applied to a specific method, identified by operation, for all APIs in all specifications.
          - `{method-name}.md` - Overlay applied to all methods with this HTTP method name, across all APIs in all specifications.
          - `method.md` - Overlay applied to all API methods.
        - `resource/`
          - `resource.md` - Overlay applied to all resource pages.
    - `sections/` - Contains additional documentation and overlays for specific OpenAPI specifications.
      - `{specification-ID}` - The kebab case of the OpenAPI `info.title` member of the specification the overlays are for.
            - `templates/`
                - `guides/` - Directory containing authored documentation for the named OpenAPI specification.
                - `reference/` - 
                  - `api.md` - Overlay applied to all API pages of this specification.
                  - `{api-group}.md` - Overlay applied to a specific API page of this specification.
                  - `{api-group}/`
                        - `{operation-ID}.md` - Overlay applied to a specific method, identified by operation, for this named API.
                        - `{method-name}.md` - Overlay applied to all methods with this HTTP method name in the named API.
                        - `method.md` - Overlay applied to all methods of this named API.
                  - `{operation-ID}.md` - Overlay applied to all methods with the given operation name, for all APIs in the specification.
                  - `{method-name}.md` - Overlay applied to all methods with this HTTP method name, across all APIs in the specification.
                  - `method.md` - Overlay applied to all methods of all APIs in the specification.
                - `resource/`
                  - `{resource-name}.md` - Overlay applied to a specific resource page of this API.
                  - `resource.md` - Overlay applied to all resource pages of this API.

## Enabling author debug mode

By passing swaggerly the configuration parameter `-author-show-assets=true`, swaggerly will display an overlay search
path pane at the foot of each API reference page. This helps you see exactly which path and filenames you can use to
overlay content onto the page you are viewing (see [Configuration parameters](#configuration-parameters)).

> Try enabling this mode for the Petstore specification to see for yourself the file paths scanned, and how they relate
to the specification.

The recommendation is to always enable this debug mode when you are writing overlay content.

## Page overlay targets

Each of the three auto-generated reference page types (api, method and resource) have their own set of overlay targets:


### API page

| GFM section reference | Page section |
| --------------------- | ------------ |
| `[[banner]]`      | Inserts content at the start of the page, before the description. |
| `[[description]]` | Adds content before the API method list. |
| `[[additional]]`  | Inserts content at the end of the API page. |

### Method page

| GFM section reference | Page section |
| --------------------- | ------------ |
| `[[banner]]`            | Inserts content at the start of the page, before the description header. |
| `[[description]]`       | Adds content after the method description. |
| `[[request]]`           | Adds content before the method request URL. |
| `[[path-parameters]]`   | Adds content before the path parameters block. |
| `[[query-parameters]]`  | Adds content before the query parameters block. |
| `[[request-headers]]`   | Adds content before the header parameters block. |
| `[[form-parameters]]`   | Adds content before the form parameters block. |
| `[[request-body]]`      | Adds content before the body block. |
| `[[security]]`          | Adds content before the security section. |
| `[[response]]`          | Adds content before the response section. |
| `[[example]]`           | Inserts example content after the response section. |
| `[[additional]]`        | Inserts additional beofre the API explorer. |


### Resource page

| GFM section reference | Page section |
| --------------------- | ------------ |
| `[[banner]]`          | Inserts banner content at the start of the page, before the description. |
| `[[description]]`     | Adds a description block to the start of the page. |
| `[[methods]]`         | Inserts content before the methods list. |
| `[[resource]]`        | Inserts content before the resource schema. |
| `[[example]]`         | Adds content before the resource example, if it exists. |
| `[[properties]]`      | Inserts content before the resource properties table. |
| `[[additional]]`      | Inserts content at the end of the resource page. |

## Example

The Swaggerly assets example `examples/overlay/assets/` contains two overlay files:

1. `templates/reference/method.md`
2. `sections/swagger-petstore/templates/reference/get.md`

The first provides `[[banner]]` and `[[request]]` overlay text for all method pages, across all
specifications, except where it is overridden by a higher precidence overlay, such as the second
overlay. This second overlay targets the Petstore specification `GET` method pages, and
provides `[[banner]]`, `[[description]]` and `[[response]]` overlay texts.

To run these examples, pass the following assets configuration to Swaggerly:

```bash
-assets-dir=<Swaggerly-source-directory>/examples/overlay/assets 
```

If you view any `GET` method page of the Petstore API, you will see the text from the second
overlay file injected into the method's `banner`, `description` and `response` sections. If you
view any other method page, you will see `banner` and `request` texts injected from the first
overlay file.

!!!HIGHLIGHT!!!
