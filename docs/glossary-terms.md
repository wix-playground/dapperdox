# Terms used by Swaggerly

## `specification-ID`

The identifier assigned to each specification Swagger loads. It is derrived from the specifications `info.title` member,
which Swaggerly lowercases and hyphen delimits ([kebab casing](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles)).

For example, a specification title of "My API suite" produces a `specification-ID` of `my-api-suite`.

## `api-group`

A logical grouping of API operations. Grouping will be made based on the tag if tagging is used by the specification,
in which case `api-group` is set to the [kebab case](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles) of
the tag's `description` member if present, or the tag `name` if not.

If tagging is not in use, then `api-group` is set to the 
[kebab case](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles) 
of the operation's `x-pathName` member if present, otherwise the operation's `summary` member is used.

## `operation-ID`

The identifier given to an operation. It is most importantly used to determine the correct content [overlay](/docs/author-overlays) to apply.

It is formed by taking the [kebab case](https://en.wikipedia.org/wiki/Letter_case#Special_case_styles) of the
first in the following list of operation members that yields a value:

1. `operationID`
2. `x-opererationName`
3. `summary`
4. method name

where method name is `method-name` below. Therefore, it is possible for `operation-ID` and `method-name` to be the same value.


## `method-name`

The name of the HTTP method for an operation, such as `get`, `post` and `put`.

## Specification list page

This page lists the API specifications that are available for reference, and is the default homepage when multiple
specifications are loaded by Swaggerly. It is not normally displayed when a single specification is loaded.

## Specification summary page

This page lists the all APIs defined by a specification, grouped by [`api-group`](/docs/glossary-terms#api-group).
It is the first page displayed when a specification is selected from the `specification-list`, and the default homepage
when a single specification is loaded by Swaggerly.

## API group summary page

This page lists the API operations within a single `api-group`.



