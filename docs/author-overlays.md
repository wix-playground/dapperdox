# Content overlays

Additional content can be added to any Swaggerly generated reference page by providing overlay files.
These pages are authored in Github Flavoured Markdown (GFM) and contain special markdown references that
target particular sections within API, Method or Resource pages.

Additional directories are added to your `assets` directory to accomplish this. As Swaggerly can consume and serve
multiple OpenAPI specifications, each is given its own section within Swaggerly, allowing you to provide guides and
overlay documentation appropriate to an OpenAPI specification. 

See [Controlling method names](/docs/author-method-names.html) for a discussion on what an *operation* name is, and
how it differs from an HTTP method name.

For example, the following GFM file adds additional content to the *request* and *response* sections
of the **Estimates of price** `get` method for the `Uber API` OpenAPI specification, where `{local_assets}` is the directory
pointed to by the `-assets-dir=` configuration parameter:

```
> cat {local_asset}/sections/uber-api/reference/estimates-of-price/get.md
```
```gfm
Overlay: true

[[request]]
It is important that this request be called with valid geo-location coordinates.

[[response]]
The response is always an array of response objects, if successful.
```

For a GFM page to be treated as an overlay, it must contain the metadata `Overlay: true` at the start
of the file (see [Supported metadata](#supported-metadata)).

There are three ways to overlay a reference page, globally, per-specification or on a page-by-page basis. Swaggerly will
look at the following file patterns in the order defined below to find any appropriate overlays, and will stop once it finds one.
For example `sections/[spec-ID]/reference/<API name>.md` takes precident over `sections/[spec-ID]/reference/api.md`.

| Reference page | Overlay filename | Description |
| -------------- | ---------------- | ----------- |
| API      | `sections/[spec-ID]/reference/<API name>.md`               | Overlay applied to a specific API page. |
| API      | `sections/[spec-ID]/reference/api.md`                      | Overlay applied to all API pages for the named specification. |
| API      | `templates/reference/api.md`                               | Overlay applied to all API pages. |
| Method   | `sections/[spec-ID]/reference/<API name>/<operation name>.md` | Overlay applied to a specific method page of a specific API of a specific openAPI specification. |
| Method   | `sections/[spec-ID]/reference/<API name>/method.md`        | Overlay applied to all method pages of a specific API. |
| Method   | `sections/[spec-ID]/reference/<operation name>.md`         | Overlay applied to all method pages for <operation name> across all APIs in a specification. |
| Method   | `sections/[spec-ID]/reference/method.md`                   | Overlay applied to all method pages of all APIs in a particular specification. |
| Method   | `templates/reference/<operation name>.md`                  | Overlay applied to the all methods of <operation name> of all APIs across all specifications.  |
| Method   | `templates/reference/method.md`                            | Overlay applied to all method pages of all APIs across all specifications.  |
| Resource | `sections/[spec-ID]/resource/<resource name>.md`           | Overlay applied to a specific resource page of a specific API.  |
| Resource | `sections/[spec-ID]/resource/resource.md`                  | Overlay applied to all resource pages of a specific API.  |
| Resource | `templates/resource/resource.md`                           | Overlay applied to all resource pages of all APIs across all specifications.  |

Shown as a directory tree:

- `{local_assets}/`
    - `templates/`
        - `guides/` - Authored documentation presented when not viewing an OpenAPI specification section.
        - `reference/` - Custom overlay GFM content.
          - `api.md` - Overlay applied to all API pages.
          - `[operation-name].md` - Overlay applied to a specific method, identified by operation, for all APIs in all specifications.
          - `method.md` - Overlay applied to all API methods.
        - `resource/`
          - `resource.md` - Overlay applied to all resource pages.
    - `sections/` - Contains additional documentation and overlays for specific OpenAPI specifications.
      - `[specification-ID]` - The kabab case of the OpenAPI `info.title` member of the specification the overlays are for.
            - `guides/` - Directory containing guides for appropriate for the named OpenAPI specification.
            - `reference/` - 
              - `api.md` - Overlay applied to all API pages of this specification.
              - `[api-name].md` - Overlay applied to a specific API page of this specification.
              - `[api-name]/`
                    - `[operation-name].md` - Overlay applied to a specific method, identified by operation, for this named API.
                    - `method.md` - Overlay applied to all methods of this named API.
              - `[operation-name].md` - Overlay applied to all methods with the given operation name, for all APIs in the specification.
              - `method.md` - Overlay applied to all methods of all APIs in the specification.
            - `resource/`
              - `resource.md` - Overlay applied to all resource pages of this API.

## Enabling author debug mode

By passing swaggerly the configuration parameter `-author-show-assets=true`, swaggerly will display an overlay search
path pane at the foot of each API reference page. This helps you see exactly which path and filenames you can use to
overlay content onto the page you are viewing (see [Configuration parameters](#configuration-parameters)).

## Page overlay targets

Each of the three auto-generated reference page types (api, method and resource) have their own set of overlay targets:


### API page

| GFM section reference | Page section |
| --------------------- | ------------ |
| `[[banner]]`      | Inserts content at the start of the page, before the API list. |

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
| `[[example]]`           | Inserts content before the API explorer. |


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

!!!HIGHLIGHT!!!
