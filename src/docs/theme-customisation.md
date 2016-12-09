Title: Customisation
Description: Customising the presentation
Keywords: header, footer, customisation, presentation
GFMMap: <li>:<li class="bullet-list">

# Customisation

The default DapperDox themes are built using [bootstrap](http://getbootstrap.com/), 
as this is one of the most flexible, well documented and well understood frameworks,
making these themes simple to customise and extend.

## Overriding assets

Overriding any of the assets in a theme is reletively straightforward, and plenty of
provision has been made in the makeup of the default themes to facilitate this.

It will be common for users to want to alter the styling, through CSS, to achieve a
particular look and feel. Additionally, users will want to tailor the header and footer
of the pages to meet their needs and achieve a level of integration with their existing
web presence.

These frequently overridden assets are described here.

### Frequently overriden assets

Taking the default theme directory `assets/themes/default/`:

Create a copies of these file in your `{local_assets}/templates` directory (see [local assets](/docs/author-concepts))
and customise to tailor the presentation to your needs.


| Asset filename | Description |
| -------------- | ----------- |
| `templates/fragments/header_bar.tmpl` | Provides the header bar for all pages. It pulls in fragments `templates/fragments/header_bar_title.tmpl` and `templates/fragments/header_bar_right.tmpl` |
| `templates/fragments/header_bar_title.tmpl` | Contains the branding for all pages and the provides the title of the specification being viewed. |
| `templates/fragments/header_bar_right.tmpl` | Supplies the content for the right-hand side of the header bar. By default provides the `All specifications` navigation, enabled when multiple specifications are being served. |
| `templates/fragments/footer.tmpl` | This template fragment provides the footer content for all pages. |
| `templates/fragments/theme.tmpl` | This template fragment imports the theme specific styles. Override this to provide styles *in addition* to the default. |
| `templates/fragments/fonts.tmpl` | This template fragment imports the fonts used by DapperDox. |
| `templates/fragments/scripts.tmpl` | This fragment sources any javascript required by pages and is loaded at the end of the page. The default `scripts.tmpl` registers an [apiExplorer callback](/docs/explorer-customising#controlling-authentication-credential-passing) to control API explorer authentication.  |

## Introducing other resources

You may need to introduce images and other files such as CSS into your local assets tree
to achieve the presentation you are after. These should be added to appropriate
sub-directories within your `{local_assets}/static/` directory.  DapperDox will make
these assets available using the file's path, excluding the `/{local_assets}/static` stem.

DapperDox will import and register files that match one of the following MIME types:

- `image*`
- `text/css`
- `javascript*`

Therefore, a local asset of `{local_assets}/static/images/my_corporate_logo.png` will be
served by DapperDox as `/images/my_corporate_logo.png`.


