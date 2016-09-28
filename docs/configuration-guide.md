# Configuration parameters

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
| `-force-root-page` | `FORCE_ROOT_PAGE` | When Swaggerly is serving a single OpenAPI specification, then by default Swaggerly will show the specification index page when serving the homepage. You can force Swaggerly to show the root index page with this option. Takes the value `true` or `false`. |
| `-author-show-assets` | `AUTHOR_SHOW_ASSETS` | Setting this value to `true` will enable an *assets search path* pane at the foot of every API reference page. This shows the path order that Swaggerly will scan to find GFM content overlay or replacement files. Takes the value `true` or `false`. |

Some configuration parameters can take multiple values, either by specifying the parameter multiple times on the command lines, or by
comma seperating multiple values when using environment variables.

