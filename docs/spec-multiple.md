# Multiple specifications

Swaggerly can document a suite of API specifications at once. This allows you to present
a portfolio of API products under a single website.

For example, an organisation could provide the following APIs:

1. A public data API
2. A streaming API
3. A product ordering API
4. An OAuth2 authentication service

These APIs collectively form a coherent suite of products, in that they all make use of the authentication 
service and provide different types of access to the same data underlying business data. The choice of which
API to use depends on the customer or business need.

If Swaggerly is asked to render documentation for the above four specifications, by passing it multiple
`-spec-filename=` options, then a page similar to the following would be shown:

![](/images/api_suite.png "Multiple API Specification page")

Each specification loaded by Swaggerly will be assigned its own [specification ID](/docs/spec-concepts.html#specification-identifiers).
