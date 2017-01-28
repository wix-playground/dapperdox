Title: Creating a theme
Description: Creating a theme
Keywords: dapperdox, theme, customisation, presentation
GFMMap: <li>:<li class="bullet-list">

# Creating a theme

To create your own theme, follow the guidance give in [theme customisation](/docs/theme-customisation), but instead of storing your overridden assets within `{local_assets}/templates`
and  `{local_assets}/static`, create them within a theme directory.

The following steps illustrate the process:

1. Create a directory within which you store all the themes you want to use with DapperDox
(see [Adding new themes](/docs/theme-overview#adding-new-thmes)). For example, this could be the directory `${HOME}/themes`:

```
$ mkdir ${HOME}/themes
```

2. Create a directory for your theme. For this example, we assume the theme name `my-theme` and follow the theme naming convention of having a `dapperdox-theme-` prefix (see [DapperDox thees on GitHub](#dapperdox-themes-on-github)):

```
$ mkdir ${HOME}/themes/dapperdox-theme-my-theme
```

### DapperDox themes on GitHub

It is recommended that GitHub repositories for DapperDox themes are  given a `dapperdox-theme-`
prefix to make them easy to find by searching GitHub or via a search engine.
