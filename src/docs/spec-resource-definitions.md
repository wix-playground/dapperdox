Title: Resource definitions
Description: How to get the best out of your specification resource definitions
Keywords: resource, title, including, excluding, controll, openAPI, swagger

# Resource definitions

## Resource title

When specifying a resource [schema object](http://swagger.io/specification/#schemaObject), either as a 
[definitions object](http://swagger.io/specification/#definitionsObject) that can be consumed or produced by
an operation, or as an inline schema object, DapperDox **requires** that the optional 
[schema object](http://swagger.io/specification/#schemaObject) 
`title` member is present.

DapperDox uses the `title` member to give the resource a name in the documentation it produces.

For example, the operation `GET /store/inventory` in the petstore example specification
`examples/specifications/petstore/swagger.json`, has a response schema with a `title` of "Quantities":

```
{
    "/store/inventory": {
        "get": {
            "tags": ["store"],
            "summary": "Returns pet inventories by status",
            "description": "Returns a map of status codes to quantities",
            "operationId": "getInventory",
            "produces": ["application/json"],
            "parameters": [],
            "responses": {
                "200": {
                    "description": "successful operation",
                    "schema": {
                        "type": "object",
                        "title":"Quantities",
                        "additionalProperties": {
                            "type": "integer",
                            "format": "int32"
                        }
                    }
                }
            },
            "security": [{
                "api_key": []
            }]
        }
    },
}
```

Without a `title` member, DapperDox would produce the error:

```
Error: GET /store/inventory references a model definition that does not have a title member.
```

## Excluding request body members from operations

In some scenarios, you might want your reference documentation to show a filtered request body resource.

To understand why, consider the following trivial `Orders` resource used by a fictional REST API:

```
{
    "order_number" : "string",
    "order_date" : "date-time",
    "reference": "string",
    "order_status": "string"
}
```

In this hypothetical example, some of these members are read-only, being automatically populated when the resource
is created, and some are supplied by the client when creating the resource which cannot be modified afterwards. 

The following table shows which members can be set or modified by a `POST` or `PUT` and are returned by a `GET`:

| Member         | POST  |  PUT  |  GET  |
| -------------- | :---: | :---: | :---: |
| `order_number` | ✘ | ✘ | ✔ |
| `order_date`   | ✔ | ✘ | ✔ |
| `reference`    | ✔ | ✔ | ✔ |
| `order_status` | ✘ | ✔ | ✔ |

For instance, the `reference` member can be set, modified and returned by all three methods, whereas
`order_date` can be set on creation (`POST`), and returned by the `GET`, but not subsequently modified (`PUT`).

As a REST API, it is expected that the same resource would be produced and consumed by all methods
acting on its URI. However, as there are subtle differences in which members can be written or read for each
method, it might be tempting to define a *different* resource for each method so that only the significant members
for the consumed or produced resource are shown in the documentation.
Doing so would break this *"identification of resources"* principle, where the URI represents a single resource.

It is usually in the documentation of such APIs that complexity arises, not in the API specification itself.

DapperDox will automatically exclude read-only resource members from the documentation of an operation's request
body (since they can only be sent to writable operations). It also gives you control over which members should be
excluded from an operation's request body documentation, through the `x-excludeFromOperations` member, added
to relevant resource properties.

`x-excludeFromOperations` takes an array of string, with each value being the operation name (which is either
the HTTP method name `get`, `post`, `put` and so on, or the `x-operationName`, if defined.
See [controlling method names](/docs/spec-method-names) for further details).

The following specification segment illustrates the use of `readOnly` and `x-excludeFromOperations`
to control the documentation of the resource when presented as a request body:

```
{
    "Orders": {
        "title" : "Orders",
        "type": "object",
        "properties": {
            "order_number": {
                "type": "string",
                "readOnly": true
            },
            "order_date": {
                "type": "string",
                "format": "date-time",
                "x-excludeFromOperations": ["put"]
            },
            "reference": {
                "type": "string"
            },
            "order_status": {
                "type": "string",
                "x-excludeFromOperations": ["post"]
            }
        }
    }
}
```












!!!HIGHLIGHT!!!
