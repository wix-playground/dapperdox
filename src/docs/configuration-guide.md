Title: Configuring DapperDox
Description: DapperDox configuration parameters
Keywords: Configuration, command line, environment variables, customisation

# DapperDox configuration

DapperDox can be configured using command line options or, in line with the [twelve-factor app](https://12factor.net/) recommendations, through environment variables. 

## Configuration options

| Option | Environment variable | Description |
| ------ | -------------------- | ----------- |
| `-assets-dir` | `ASSETS_DIR` | Assets to serve. Effectively the document root. |
| `-bind-addr` | `BIND_ADDR` | IP address and port number that DapperDox should serve content from. Value takes the form `<IP address>:<port>`. |
| `-default-assets-dir` | `DEFAULT_ASSETS_DIR` | Used to point DapperDox at its default assets. You Will need to provide this option if you are not running DapperDox from within it's distribution directory. For example, if you invoke DapperDox as `<path_to_dapperdox_distribution>/dapperdox` then you must pass this configuration option as `-default-assets-dir=<path_to_dapperdox_distribution>/assets`.|
| `-document-rewrite-url` | `DOCUMENT_REWRITE_URL` | Specify a URL that is to be rewritten in the documentation. May be multiply defined. Format is from=to. This is applied to assets, not to OpenAPI specification generated text. |
| `-log-level` | `LOGLEVEL` | Log level (`info`, `debug`, `trace`) |
| `-spec-dir` | `SPEC_DIR` | OpenAPI specification (swagger) directory. |
| `-spec-filename` | `SPEC_FILENAME` | The filename of the OpenAPI specification file within the spec-dir. Defaults to spec/swagger.json |
| `-site-url` | `SITE_URL` | Public URL of the documentation service. Used by `-spec-rewrite-url` if no `=to` is supplied to the rewrite. |
| `-spec-rewrite-url` | `SPEC_REWRITE_URL` | The URLs in the OpenAPI specifications to be rewritten as `site-url`, or to the `to` URL if the value given is of the form from=to. Applies to OpenAPI specification text, not asset files. |
| `-theme` | `THEME` | Name of the theme to render documentation. |
| `-themes-dir` | `THEMES_DIR` | Directory containing installed themes. |
| `-force-specification-list` | `FORCE_SPECIFICATION_LIST` | When DapperDox is serving a single OpenAPI specification, then by default it will show the API summary page when serving the homepage. Instead, you can force DapperDox to show the list of available specifications (as it would if there were more than one specification) with this option. This is necessary if you have global documentation guides which live outside the specification. Takes the value `true` or `false`. |
| `-author-show-assets` | `AUTHOR_SHOW_ASSETS` | Setting this value to `true` will enable an *assets search path* pane at the foot of every API reference page. This shows the path order that DapperDox will scan to find GFM content overlay or replacement files. Takes the value `true` or `false`. |
| `-proxy-path` | `PROXY_PATH` | Configures a path prefix that is to be reverse-proxied to another service. Value takes the form `source-path=service-host/destination/path`. If `destination-path` is given, this will prefix the `source-path` that is passed to the service. See [reverse proxy](/docs/proxy-configure) for further details.. |
| `-tls-certificate` | `TLS_CERTIFICATE` | Enables HTTPS. The path and filename of the TLS certificate to load. Must be accompanied by the `-tls-key` configuration. |
| `-tls-key` | `TLS_KEY` | Enables HTTPS. The path and filename of the TLS private key to load. Must be accompanied by the `-tls-certificate` configuration. |

## Multiple values

Some configuration parameters can take multiple values, either by specifying the parameter multiple times on the command lines, or by
comma separating multiple values when using environment variables.

