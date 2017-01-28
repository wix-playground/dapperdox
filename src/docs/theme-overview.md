Title: Themes and styling
Description: Overview of the themeing and styling of documentation in DapperDox
Keywords: theme, css, style, presentation, layout

# Themes and styling

The HTML templates, CSS, images and javascript files used by DapperDox are known as
assets. Assets that collectively provide a particular look-and-feel are known as a theme.

By providing local assets, individual files can be overridden or customised to tailor
DapperDox's presentation to meet your requirements, alternatively an entire theme can be 
created and imported into DapperDox.

If you provide local asset files, these are placed in the same [assets](/docs/author-concepts#assets) directory structure as your authored content and guides.

## Supplying local assets

DapperDox is directed to your local assets through the `-assets-dir` [configuration](/docs/configuration-guide) option. These assets work within, and perhaps customise, the theme that you are
using.

For details on what asset files you can and should override, refer to the section on [customisation](/docs/theme-customisation).

## Changing themes

To switch themes, pass DapperDox the `-theme` [configuration](/docs/configuration-guide) 
option giving the name of the theme you want to use. For example, select to the 
`sectionbar` theme
(see [built-in themes](/docs/theme-overview#built-in-themes)) with the option:

```
-theme=sectionbar
```

## Built-in themes

DapperDox ships with two built-in themes.

1. `default`
2. `sectionbar`

Both these themes are built using [bootstrap](http://getbootstrap.com/), 
as this is one of the most flexible, well documented and well understood frameworks,
making these themes simple to customise and extend.

#### default theme

This theme combines the navigation for reference documentation and authored guides into a
single page side-navigation.

#### sectionbar theme

This theme separates the reference documentation from the authored guides for a specification,
and provides a navigation bar at the top of the screen to allow the user
to switch between the guide and reference sections. The page side-navigation provides 
navigation within the current section.

## Downloading a theme

Overtime the DapperDox team and third parties will release additional themes. Once such example
is the GOV.UK theme, produced by [CompaniesHouse](#).
This is available from their [GitHub repository](https://github.com/companieshouse/dapperdox-theme-gov-uk).

Other themes are available, search for `dapperdox-theme` on GitHub or your favourite search
(see [Dapperdox themes on GitHub](/docs/theme-creating#dapperdox-themes-on-github)).

## Adding new themes

If you are providing your own themes, either created yourself or downloaded, then they 
should be placed within a common directory and DapperDox configured to look for them there.
To do this, pass DapperDox the `-theme-dir` [configuration](/docs/configuration-guide)
option.

DapperDox will fall back to its default theme to find any asset files not provided by an
imported theme.

For details on using and creating new themes, refer to the section on [creating a theme](/docs/theme-creating).

