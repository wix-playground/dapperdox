# Configuration parameters

| Option | Environment variable | Description |
| ------ | -------------------- | ----------- |
| `-assets-dir` | `ASSETS_DIR` | Assets to serve. Effectively the document root. |
| `-bind-addr` | `BIND_ADDR` | IP address and port number that Swaggerly should serve content from. Value takes the form `<IP address>:<port>`. |
| `-default-assets-dir` | `DEFAULT_ASSETS_DIR` | Used to point Swaggerly at its default assets. You wil need to provide this option if you are not running Swaggerly from within it's distribution directory. For example, if you invoke Swaggerly as `<path_to_swaggerly_distribution>/swaggerly` then you must pass this configuration option as `-default-assets-dir=<path_to_swaggerly_distribution>/assets`.|
| `-document-rewrite-url` | `DOCUMENT_REWRITE_URL` | Specify a URL that is to be rewritten in the documentation. May be multiply defined. Format is from=to. This is applied to assets, not to OpenAPI specification generated text. |
| `-log-level` | `LOGLEVEL` | Log level (`info`, `debug`, `trace`) |
| `-spec-dir` | `SPEC_DIR` | OpenAPI specification (swagger) directory. |
| `-spec-filename` | `SPEC_FILENAME` | The filename of the OpenAPI specification file within the spec-dir. Defaults to spec/swagger.json |
| `-site-url` | `SITE_URL` | Public URL of the documentation service. Used by `-spec-rewrite-url` if no `=to` is supplied to the rewrite. |
| `-spec-rewrite-url` | `SPEC_REWRITE_URL` | The URLs in the OpenAPI specifications to be rewritten as `site-url`, or to the `to` URL if the value given is of the form from=to. Applies to OpenAPI specification text, not asset files. |
| `-theme` | `THEME` | Name of the theme to render documentation. |
| `-themes-dir` | `THEMES_DIR` | Directory containing installed themes. |
| `-force-specification-list` | `FORCE_SPECIFICATION_LIST` | When Swaggerly is serving a single OpenAPI specification, then by default it will show the API summary page when serving the homepage. Instead, you can force Swaggerly to show the list of available specificatons (as it would if there were more than one specification) with this option. This is necessary if you have global documentation guides which live outside the specification. Takes the value `true` or `false`. |
| `-author-show-assets` | `AUTHOR_SHOW_ASSETS` | Setting this value to `true` will enable an *assets search path* pane at the foot of every API reference page. This shows the path order that Swaggerly will scan to find GFM content overlay or replacement files. Takes the value `true` or `false`. |

Some configuration parameters can take multiple values, either by specifying the parameter multiple times on the command lines, or by
comma seperating multiple values when using environment variables.

