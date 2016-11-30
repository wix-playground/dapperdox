Title: Customising the API explorer
Description: How to customise the API explorer in DapperDox
Keywords: API, explorer, openAPI, Swagger, swagger-ui, customise, API key, OAuth2, methods

<h1>Customising the API explorer</h2>

<p>The <code>apiExplorer</code> JavaScript object provides a method to add API keys to an internal list, and a method to inject those
API keys into the explorer page, so that the user can select the key from a drop-down menu instead of having to type it in.</p>

<p>By default, the template fragment
<code>assets/themes/default/templates/fragments/scripts.tmpl</code> (shown below) controls the
handling of authentication credentials for the explorer:</p>

<pre><code>&lt;!-- Additional scripts to be loaded at end of page 
  -- This should be overridden to take control of the auth process (adding keys to the explorer request).
  -->
&lt;script>
   $(document).ready(function(){
        // Register callback to add authorisation parameters to request before it is sent
        apiExplorer.setBeforeSendCallback( function( request ) {
            var apiKey      = apiExplorer.readApiKey();      // Read API key from explorer input
            var accessToken = apiExplorer.readAccessToken(); // Read access token from explorer input
            var basicAuth   = apiExplorer.getBasicAuthentication(); // Create basic auth string

            // Favour access tokens over api keys
            if( accessToken != "" ) { request.headers = {Authorization: "Bearer "+accessToken}; }
            else if( apiKey != "" ) { request.params  = {key: apiKey}; }
            else if( basicAuth != "" ) { request.headers = {Authorization: "Basic "+basicAuth}; }
        });
    });
&lt;/script></code></pre>

<p>It would be usual for this template fragment to be overridden by an installation, so that a
list of valid API keys can be built and injected into the explorer page. This allows a user of
the explorer to select an API key from a menu, to be used in the request.
Generally the keys would be fetched from some AJAX endpoint that you provide, once the user as
gone though some sign-in process.</p>

<h2><a name="controlling-authentication-credential-passing" class="anchor" href="#controlling-authentication-credential-passing" rel="nofollow" aria-hidden="true"></a>Controlling authentication credential passing</h2>

<p>The example <code>examples/apikey_injection/assets/templates/fragments/scripts.tmpl</code> shows how the default <code>scripts.tmpl</code> can be overridden to inject an API key
(hard-coded for the benefit of this example), as a list-of-one into the explorer page:</p>

<pre><code>&lt;!-- Inject API key(s). This would probably be done via an ajax request to a server to request the API
     keys for the signed in user.
     Register callback to appropriately add the authentication credentials (as a Basic auth header)
     to the request before it is sent.
  -->
&lt;script type="text/javascript">
    $(document).ready(function(){
        apiExplorer.addApiKey("A test key","2jnD1-ZnGBsT2ST7mSm9ASaGxO7BPWU4iz9TlfE6");
        apiExplorer.injectApiKeysIntoPage();

        // Register callback to add authorisation parameters to request before it is sent
        apiExplorer.setBeforeSendCallback( function( request ) {
            var apiKey = apiExplorer.readApiKey(); // Read API key from explorer input
            // Set the API key in the request as an Authorization header, using BASIC authentication
            request.headers = {Authorization:"Basic " + btoa(apiKey + ":")};
        } );
    });
&lt;/script></code></pre>

<p>To run this example, DapperDox needs to be told about the example assets directory for it to pick up the override. 
This is achieved through the configuration parameter <code>-assets-dir</code>, passed to DapperDox when starting: </p>

<pre><code>-assets-dir=examples/apikey_injection/assets</code></pre>

<p>By default, DapperDox will automatically attach the API key if supplied, to the request URL as a <code>key=</code> query parameter.
This behaviour can be customised to satisfy the authentication requirements of your API.</p>

<p>The template fragment <code>assets/themes/default/templates/fragments/scripts.tmpl</code>, which is included at the end of the common
page template <code>layout.tmpl</code> contains the following JavaScript snippet:</p>


<h2>apiExplorer methods</h2>

<p>The API explorer object <code>apiExplorer</code> provides a number of methods:</p>

<div class="table-responsive">
    <table class="table table-striped">
        <thead>
            <th>Method</th>
            <th>Description</th>
        </thead>
        <tbody>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.addApiKey(<span class="hljs-params">name</span>, <span class="hljs-params">key</span>)</span></code> </td><td>This method adds the named key to the internal list.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.listApiKeys()</span></code></td><td>Returns an array of key names.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.getApiKey(<span class="hljs-params">name</span>)</span></code> </td><td>Returns the key associated with <code class="hljs-params">name</code>.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.injectApiKeysIntoPage()</span></code></td><td> Injects the named API keys into the explorer, building a pulldown menu that can be selected from by the user.</td></tr>

            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.injectApiKeysIntoPage()</span></code></td><td>Injects the named API keys into the explorer, building a pulldown menu that can be selected from by the user.</td></tr>

            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.<br />setBeforeSendCallback(<span class="hljs-params">function</span>)</span></code></td><td>Register a callback function to be called just prior to sending the API request. The callback is passed a request object, with <code>headers</code> and <code>params</code> properties. See <a href="#setting-headers-and-query-parameters">setting headers and query parameters</a>.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.readApiKey()</span></code></td><td>Returns the value of the API key input field or menu selected API key.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.readAccessToken()</span></code></td><td>Returns the value of the access_token input field.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.readBasicUsername()</span></code></td><td>Returns the value of the Basic authorization username input field.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.readBasicPassword()</span></code></td><td>Returns the value of the Basic authorization password input field.</td></tr>
            <tr><td><code><span class="hljs-object">apiExplorer</span><span class="hljs-keyword">.getBasicAuthentication()</span></code></td><td>Fetches the username and password from the relevant input fields, and returns a Basic Authorization encode string (excluding the prefix <code>Basic</code>).</td></tr>
        </tbody>
    </table>
</div>


<h2><a name="setting-headers-and-query-parameters" class="anchor" href="#setting-headers-and-query-parameters" rel="nofollow" aria-hidden="true"></a>Setting headers and query parameters</h2>

<p>The javascript examples on this page register a callback via the explorer method <code>apiExplorer.setBeforeSendCallback()</code>, which gets invoked while the explorer is building the 
request to send to the API. The callback will be receive an empty object which has two properties that can be set as needed,
<code>request.headers</code> - items that are sent as request HTTP headers,  and <code>request.params</code> - items that are sent as query parameters:</p>

<pre><code class="JavaScript">{
    <span class="hljs-attr">headers</span>: <span class="hljs-value">{}</span>,
    <span class="hljs-attr">params</span>: <span class="hljs-value">{}</span>
}</code></pre>

<p>Both the <code>headers</code> and <code>params</code> objects contain zero or more name/value pairs:</p>

<pre><code class="JavaScript">{
    <span class="hljs-attr">key1</span>: <span class="hljs-string">value</span>,
    ..
    ..
    <span class="hljs-attr">key_n</span>: <span class="hljs-string">value_n</span>
}</code></pre>

<p>For example:</p>

<pre><code class="JavaScript">request.headers = { header: "value" };
request.headers = { header1: "value1", header2: "value2" }</code></pre>

<p>To put this into practice, if you wanted to change the authentication credential passing mechanism to instead supply the API key
as an Authorization header, then create a <code>scripts.tmpl</code> within your own assets directory to override this. For example, the
DapperDox example file <code>examples/apikey_injection/assets/templates/fragments/scripts.tmpl</code> passes the API key in the 
Authorization header using BASIC authentication:</p>

<pre><code>$(document).ready(function(){
    // ... other code cut from here ...

    // Register callback to add authorisation parameters to request before it is sent
    apiExplorer.setBeforeSendCallback( function( request ) {
        var apiKey = apiExplorer.readApiKey(); // Read API key from explorer input
        request.headers = {Authorization:"Basic " + btoa(apiKey + ":")};
    });
});</code></pre>

<p>DapperDox then needs to be told about this local assets directory for it to pick up the override, which is achieved through
the configuration parameter <code>-assets-dir</code>, passed to DapperDox when starting. For example, to pick up the example above, use
<code>-assets-dir=examples/apikey_injection/assets</code>.</p>

<p>See <a href="/docs/author-concepts">Authoring content</a> and associated pages for further information about creating custom assets.</p>
<script>hljs.initHighlightingOnLoad();</script>
