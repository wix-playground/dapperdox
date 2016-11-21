Title: Operation tagging
Description: Using operation tagging to control what documentation is built
Keywords: tagging, operation, path, openAPI, swagger, navigation, documentation

# Operation tagging

For most APIs, the grouping of operations into related sets helps users understand and reference the
documentation. 

DapperDox will try to group your API operations as best it can, and by default will group by API
path, but you can take control of this grouping by using [tags](http://swagger.io/specification/#tagObject).

## Grouping by tag

DapperDox will look for the specification's top level `tags` declaration member and, if 
present, it restricts itself to documenting only those API operations that have tags declared in
this set. Operations will be grouped by tag and listed in the navigation in the order in which
they are declared. 
This allows you to control which reference documentation gets presented. The name of each group
is taken from the tag's `description` member if it has one, otherwise it takes the tag's `name`
member, and is used as the group description or heading in pages and navigation.

If a path has an `x-pathName` member, then its value will override the group name inherited from 
the `tags` `description` or `name` member.

This group name is also used to form the `api-group` identifier,
used in URLs and [overlays](/docs/author-overlays).
See the [Glossary of terms](/docs/glossary-terms#api-group) for further details.

```json
{
  "swagger": "2.0",

  "tags": [
    { 
        "name": "Products",
        "description": "A more verbose description of tag"
    },
    { "name": "Estimates Price" }
  ],
    "paths": {
        "/products": {
            "get": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Read product types"
            },
            "post": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Create product types"
            }
        }
    }
}
```

This incomplete specification example shows how documentation order and filtering is controlled by
`tags`. The top level tags member declares that API operations tagged with `Estimates Price` and 
`Products` should be built, in that order.  The `description` member of the `Products` tag is used
to name all operations grouped by that tag. The name of the `Estimates Price` tag would be used to
name all operations grouped by that tag, as there is no `description` member.

This mechanism for naming and grouping API operations gives you the most control.

## Grouping without tags

If tags cannot be used, DapperDox must still have a title to use for an operation path, and will 
fall back to using the `summary` member from one of the operations for a path. This often does not 
produce the desired results, unless the `summary` members of all operations for a path are set to 
the same text, as in the example above, but will rarely be the case.

To force the group name of all operations declared for a path, the DapperDox specific `x-pathName`
member may be specified in the Path Item object. This will always take effect, and will even
override the description or name inherited from the top level `tags` member. However, tags are the
most flexible approach to name method groups. 
See the [Glossary of terms](/docs/glossary-terms#api-group) for further details.


```json
{
  "swagger": "2.0",

  "tags": [
    { 
        "name": "Products",
        "description": "A more verbose description of tag"
    },
    { "name": "Estimates Price" }
  ],
    "paths": {
        "/products": {
            "x-pathName": "Types of Product",
            "get": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Read product types"
            },
            "post": {
                "tags": [
                  "Products"
                ],
                "summary": "Product Types",
                "description": "Create product types"
            }
        }
    }
}
```

!!!HIGHLIGHT!!!
