Title: Creating a theme
Description: Creating a theme
Keywords: dapperdox, theme, customisation, presentation
XGFMMap: <li>:<li class="bullet-list">

# Creating a theme

To create your own theme, follow the guidance give in [theme customisation](/docs/theme-customisation), but instead of storing your overridden assets within `{local_assets}/templates`
and  `{local_assets}/static`, create them within a theme directory.

The following steps illustrate the process:

1. Create a directory within which you store all the themes you want to use with DapperDox
(see [Adding new themes](/docs/theme-overview#adding-new-themes)). For example:
    ```
    $ mkdir ${HOME}/themes
    ```

2. Create a directory for your theme. For this example, we assume the theme name `my-theme` and follow the theme naming convention of having a `dapperdox-theme-` prefix (see [DapperDox themes on GitHub](#dapperdox-themes-on-github)):
    ```
    $ mkdir ${HOME}/themes/dapperdox-theme-my-theme
    ```

4. Create directories to contain your static theme assets and templates:
    ```
    $ mkdir ${HOME}/themes/dapperdox-theme-my-theme/static
    $ mkdir ${HOME}/themes/dapperdox-theme-my-theme/templates
    $ mkdir ${HOME}/themes/dapperdox-theme-my-theme/templates/fragments
    ```

5. Copy as many files as you need from DapperDox's default theme into your theme directory tree. For example, lets change the default page header title and provide an additional stylesheet to change the background colour:
    ```
    $ cp <directory_of_dapperdox_install>/assets/themes/default/templates/fragments/header_bar_title.tmpl ${HOME}/themes/dapperdox-theme-my-theme/templates/fragments
    ```
    and edit it to change the page title:
    ```HTML
    [&colon; if .Info.Title :]
    <a class="navbar-brand" href="[:$.SpecPath:]/reference">My fantastic.com API</a>
    <span>[: .Info.Title :]</span>
    [&colon; else :]
    <a class="navbar-brand" href="/">APIs available at fantastic.com</a>
    [&colon; end :]
    ```

6. Create your stylesheet:
    ```
    $ echo "html { background-color: lightgrey; }" > ${HOME}/themes/dapperdox-theme-my-theme/static/mystyle.css
    ```
    and get DapperDox to load it:
    ```
    $ echo '<link href="/css/mystyle.css" rel="stylesheet">' > ${HOME}/themes/dapperdox-theme-my-theme/templates/fragments/theme.tmpl
    ```

Finally, tell DapperDox where to find your themes, and which theme to use:
```
dapperdox -theme-dir=${HOME}/themes -theme=dapperdox-theme-my-theme
```

If you intend to make your theme public, then ideally publish it on GitHub. Initialise your
theme directory as a git repository, and push that to GitHub (this assumes you're already
profficient with the use of `git` and GitHub):
```
$ cd ${HOME}/themes/dapperdox-theme-my-theme
$ git init
```

### DapperDox themes on GitHub

It is recommended that GitHub repositories for DapperDox themes are  given a `dapperdox-theme-`
prefix to make them easy to find by searching GitHub or via a search engine.
