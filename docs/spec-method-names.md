# Controlling method names

## Methods and operations

Swaggerly allows you to control the name of each method operation, and how methods are represented in the navigation.

By default the HTTP method is used (`GET`, `POST`, `PUT` etc), as is usual for RESTful API specifications. Even so, it is
often the case that a resource will be given two `GET` methods, one which returns a single resource, and one that
returns a list. This is clearly confusing for the reader, so in the latter case it would be clearer for the list method
to be identified as such.

To do this, Swaggerly introduces the concept of a method **operation**, which usually has the same value as the
HTTP method, but may be overridden on a per-method basis by the `x-operationName` extension:

```json
{
    "paths": {
        "/products": {
            "get": {
                "x-operationName": "list",
                "summary": "List products",
                "description": "Returns a list of products"
                "tags": ["Products"]
            },
            "post": {
                "summary": "Create a product",
                "description": "Create a new product type"
                "tags": ["Products"]
            }
        },
        "/products/{id}": {
            "get": {
                "summary": "Get a product",
                "description": "Returns a products"
                "tags": ["Products"]
            }
        }
    }
}
```

The methods in the above example would be grouped together because they are similarly tagged, and the confusion between
the two `GET` methods resolved by giving one the operation name of `list`. Thus, the three methods would be described in the
navigation and API list page as being `list`, `post` and `get`. The HTTP methods assigned to each (get, post and get) remain unchanged.

This technique can also be used to override `GET` to `fetch`, `POST` to `create`, `PUT` to `update` and so on.

## Navigating by method name

Where an openAPI specification is describing a non-RESTful set of APIs, they are often grouped together and share the same
HTTP method. For example, two `GET` methods having respective `summary` texts of `lookup product by ID` and `lookup product by barcode` 
would probably be listed together, both being product APIs. As they are both `get` methods, the reader would be unable to tell them
apart if they are both referred to as `get` operations in the API navigation. By adding the `"x-navigateMethodsByName" : true` 
extension to the top level openAPI specification, you can force Swaggerly to describe each method in the navigation using its 
`summary` text instead of its operation name or HTTP method. The methods will continue to be referred to by operation name or
HTTP method in API pages.

```json
{
    "swagger": "2.0",
    "x-navigateMethodsByName": true,
    "info": {
        "title": "Example API",
        "description": "The only way is up",
        "version": "1.0.0"
    }
}
```
!!!HIGHLIGHT!!!
