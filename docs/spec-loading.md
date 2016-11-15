# Loading specifications

DapperDox can document a single API specification, or a suite of API specification at once.

With a single specification, DapperDox will present the [specification summary page](/docs/glossary-terms#specification-summary-page) as the homepage (`http://localhost:3128/`)
where the all APIs defined by the specification are listed.

With multiple specifications, DapperDox presents the [specification list page](/docs/glossary-terms#specification-list-page)
which catalogues each of the API specifications available.

## Multiple specifications

DapperDox can document a suite of API specifications at once. This allows you to present
a portfolio of API products under a single website.

For example, an organisation could provide the following APIs:

1. A public data API
2. A streaming API
3. A product ordering API
4. An OAuth2 authentication service

These APIs collectively form a coherent suite of products, in that they all make use of the authentication 
service and provide different types of access to the same data underlying business data. The choice of which
API to use depends on the customer or business need.

If DapperDox is asked to render documentation for the above four specifications, by passing it multiple
`-spec-filename=` options, then a page similar to the following would be shown:

<div class="img-border"><div class="fiximage img-responsive"><img src="/images/api_suite.png" /></div></div>

Each specification loaded by DapperDox will be assigned its own [specification ID](#specification-identifiers).

## Specification identifiers

DapperDox assigns a specification it's own identifier, and it's own base documentation URL. It derives the ID from the
specifications `info.title` member, which it lowercases and hyphen delimits 
([kebab casing](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles)).
For example, the
streaming API has a title of "Streaming Data API", which DapperDox converts into a specification ID of `streaming-data-api`.

The base documentation URL for a specification will be:

```http://localhost:3123/{specification-ID}/```

Thus, the streaming API example above would have a base URL of `http://localhost:3123/streaming-data-api/`.

It is important to know the specification ID when authoring [guides](/docs/author-guides) and documentation
[overlays](/docs/author-overlays).

