# Loading specifications

DapperDox can document a single API specification, or a suite of API specification at once.

With a single specification, DapperDox will present the [specification summary page](/docs/glossary-terms#specification-summary-page) as the homepage (`http://localhost:3128/`)
where the all APIs defined by the specification are listed.

With multiple specifications, DapperDox presents the [specification list page](/docs/glossary-terms#specification-list-page)
which catalogues each of the API specifications available.

## Single specifications

Point DapperDox at the directory containing your Swagger specification using the `-spec-dir` configuration
parameter, where DapperDox will look for the file `swagger.json`. If your specification has a different 
filename, specify this with the `-spec-filename` parameter.

On launching DapperDox, the first web page shown will be the 
[specification summary page](/docs/glossary-terms#specification-summary-page)
where the all APIs defined by the specification are listed.

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

DapperDox can be asked to render documentation for all of the above four specifications, by placing
them beneath a common parent directory, specified by the `-spec-dir` option, and passing each relative
specification path with multiple `-spec-filename=` options.

For example, this specification directory tree could be structured as:
```
/user/api_specs/streaming/swagger.json
/user/api_specs/public/swagger.json
/user/api_specs/ordering/swagger.json
/user/api_specs/authorisation/swagger.json
```
for which, DapperDox would be configured with:
```
-spec-dir=/user/api_specs \
-spec-filename=streaming/swagger.json \
-spec-filename=public/swagger.json \
-spec-filename=ordering/saggier.json \
-spec-filename=authorisation/swagger.json
```

producing an initial documentation page of:

<div class="img-border"><div class="fiximage img-responsive"><img src="/images/api_suite.png" /></div></div>

Each specification loaded by DapperDox will be assigned its own [specification ID](#specification-identifiers).

## Specification identifiers

DapperDox assigns a specification it's own identifier, and it's own base documentation URL. It derives the ID from the
specifications `info.title` member, which it lower-cases and hyphen delimits 
([kebab casing](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles)).
For example, the
streaming API has a title of "Streaming Data API", which DapperDox converts into a specification ID of `streaming-data-api`.

The base documentation URL for a specification will be:

```http://localhost:3123/{specification-ID}/```

Thus, the streaming API example above would have a base URL of `http://localhost:3123/streaming-data-api/`.

It is important to know the specification ID when authoring [guides](/docs/author-guides) and documentation
[overlays](/docs/author-overlays).

