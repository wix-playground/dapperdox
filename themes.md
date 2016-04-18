Swaggerly Themes
================

**INCOMPLETE**

## Theme assets

A typical theme assets structure is:

- `theme name/`
    - `static/`
        - `css/` - Theme specific stylesheets
        - `js/` - Theme specific javascript
    - `templates/`
        - `layout.tmpl` - The common page template
        - `default-api.tmpl` - The default API page template
        - `default-method.tmpl` - The default API method page template
        - `default-resource.tmpl` - The default API resource page template
        - `fragments/` - Common HTML fragments used across pages
            - `reference/` - Fragments such as the page side navigation and authorisation details
            - `explorer/` - API explorer fragments, such as input fields

### Using an alternative theme
```
   -themes-dir=<installation directory of additional themes>
   -theme=<name of theme to use>
```

### Adding custom documentation

One of the main features of Swaggerly is its ability to serve specification generated reference documentation along with
authored guides and reference notes, allowing the production of a complete unified documentation suite for an API, or set of
APIs.

To do this, Swaggerly would be configured to pick up and serve additional documentation, written in HTML or Github Flavoured 
Markdown (styling support for which is currently incomplete). These documents are grouped under an application-specific `assets` directory.
The structure of this assets directory is important, as it also allows you to override assets that are normally supplied
by the chosen theme, and so must follow the structure set out in the section called [Theme assets](#theme-assets).

Custom documentation introduces a new directory into the assets structure, that of `assets/templates/guides/`.
Under this directory is placed all the authored content you wish to publish. The content may be a `.tmpl` HTML template or
a `.md` Github Flavoured Markdown file, and will be rendered within the common `theme name/templates/layout.tmpl` page
template.

Swaggerly will serve this content from the `/guides/` URL path, for example `http://localhost:3123/guides/my_example_page`.


#### Adding navigation

Swaggerly will automatically build the necessary navigation for the swagger spec generated content, but it cannot do this
for authored content within the `assets/templates/guides` directory. It leaves you to manage this appropriately. This is
straightforward to do, as you can override the theme `templates/fragments/sidenav.tmpl` and introduce any navigation links you require. Simply copy this file from the theme directory into your application-specific assets as `assets/templates/fragments/sidenav.tmpl` and amend as appropraite.

**(TODO: Investigate whether directives could be added to /guides/* files, giving navigation hints for auto-building of navigation links, which may be sufficient in many cases)**

### Overriding theme assets

Any theme assets may be overridden, as shown in the previous section with `sidenav.tmpl`. This allows you to modify any of
the presentation to suit your local requirements.


### Customisation of auto-generated pages

Every reference page, auto-generated from the swagger specification, can be customised. These pages fall into three categories:

1. API pages. These list and document all the methods available for an API endpoint.
2. Method pages. These document a particular method, acting on an API endpoint.
2. Resource pages. These document the resource that is submitted to or returned from an API endpoint, depending on method.

There are three theme templates that provide the presentation structure for each of these pages, and can therefore be
overridden if you wish to affect a documentation change consistently across all instances of a particular page type.
Alternatively, you may override a page template on an API, method and resource basis.

*The consequence of overriding theme templates is that you cannot change themes without changing your overridden pages. A
Git Flavoured Markdown approach is being developed which will make theme template overriding unnecessary for a majority of
cases.*

#### API pages

An API page will have a URL formed with the following pattern:

```
/reference/<api-group-name>
```

where `api-group-name` is the kabab formatted name derived from the tag description or name, or `x-pathName` as described
in [Specification requirements](#specification-requirements).

To override the API page for an endpoint called `example-api`, when using the default theme:

1. Copy `assets/themes/default/templates/default-api.tmpl` to `local-assets/assets/templates/reference/example-api.tmpl`
2. Edit `example-api.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for the API `example-api`, it will use your overridden template.
		
#### Method pages

A method page will have a URL formed with the following pattern:

```
/reference/<api-group-name>/<method-name>
```

where `api-group-name` is the kabab formatted name derived from the tag description or name, or `x-pathName` as described
in [Specification requirements](#specification-requirements), and `method-name` will be `GET`, `POST` etc

To override the `GET` method page for an endpoint called `example-api`, when using the default theme:

1. Copy `assets/themes/default/templates/default-method.tmpl` to `local-assets/assets/templates/reference/example-api/get.tmpl`
2. Edit `get.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for that API method, it will use your overridden template.

#### Resource pages

An API page will have a URL formed with the following pattern:

```
/resources/<resource-name>
```

There can only be one resource of a given name.

To override the resource page for an example resource `example-resource`, when using the default theme:

1. Copy `assets/themes/default/templates/default-resource.tmpl` to `local-assets/assets/templates/resources/example-resource.tmpl`
2. Edit `example-resource.tmpl` to add any additional content you required.

Whenever Swaggerly renders the page for this resource, it will use your overridden template.

