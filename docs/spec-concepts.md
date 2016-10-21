# Specifications concepts

## Specification identifiers

Swaggerly assigns a specification it's own identifier, and it's own base URL. It derives the ID from the
specifications `info.title` member, which it lowercases and hyphen delimits (kabab casing).  For example, the
streaming API has a title of "Streaming Data API", which Swaggerly converts into a specification ID of `streaming-data-api`.

The default base URL for a specification's documentation pages is:

```http://localhost:3123/{specification-ID}/```

Thus, the streaming API example above would have a base URL of `http://localhost:3123/streaming-data-api/`.

It is important to know the specification ID when authoring [guides](/docs/author-guides.html) and documentation
[overlays](/docs/author-overlays.html).
