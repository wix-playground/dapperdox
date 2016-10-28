# Authoring content concepts

Swaggerly combines reference documentation automatically generated from OpenAPI Swagger specifications
with other authored content to produce rich, browsable, API documentation.

## Assets

Any additional content that Swaggerly imports and combines with reference documentation is known
as an asset. Assets can be HTML,
[Github Flavoured Markdown](https://help.github.com/articles/basic-writing-and-formatting-syntax/) (GFM),
and static resources such as images.

Assets are split into two groups:

1. Those applied across the suite of specifications being served
2. Those applied to particular specifications

Within these groups, assets divide into three categories:

1. Authored documentation, refered to as [guides](/docs/author-guides.html)
2. [Content overlays](/docs/author-overlays.html), containing documentation to be injected into API reference pages
3. Static resources such as images, which are used by guides or overlays.

### Assets directory structure

- `{local_assets}/`
  - `templates/`
     - `guides/`
     - `reference/`
     - `resource/`
  - `static/`
  - `sections/`
        - `{specification-ID}/`
          - `templates/`
             - `guides/`
             - `reference/`
             - `resource/`
          - `static/`

For example, `{local_assets}/templates/` contains assets that are applied across the suite of specifications.
`{local_assets}/sections/{specification-ID}/templates/` contains assets that are applied to a particular specification.

### Configuring assets

Swaggerly is directed to your assets through the `--assets-dir` [configuration](/docs/configuration-guide.html) option.
