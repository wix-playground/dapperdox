// --------------------------------------------------------------------------------------
//
var apiExplorer = { _apiKeys: {} };

apiExplorer.addApiKey = function(name,key) {
    this._apiKeys[name] = key;
}
apiExplorer.listApiKeys = function(){
    return Object.keys(this._apiKeys);
}
apiExplorer.getApiKey = function(name){
    return this._apiKeys[name];
}
apiExplorer.injectApiKeysIntoPage = function() {
    var select = document.getElementById("api-key-select");

    if( select == null ) return;

    var names = this.listApiKeys();
    var len   = names.length;

    if( len == 0 ) {
        $('#api-key-select').hide();
        $('#api-key-input').show();
        return;
    } 

    $('#api-key-select').show();
    $('#api-key-input').hide();

    for (var i = 0; i < len; i++) {
        var option = document.createElement("option");
        option.text = names[i];
        option.setAttribute("value", this.getApiKey(option.text) );
        select.appendChild(option);
    }
};
apiExplorer.setBeforeSendCallback = function( func ) {
    this._extendCallback = func;
}

// Read the API get from the explorer input parameters.
apiExplorer.readApiKey = function() {
    return $('#api-key-select').val() || $('#api-key-input').val() || "";
};
apiExplorer.readAccessToken = function() {
    return $('#access-token-input').val() || "";
};
apiExplorer.readBasicUsername = function() {
    return $('#basic-username-input').val() || "";
};
apiExplorer.readBasicPassword = function() {
    return $('#basic-password-input').val() || "";
};
apiExplorer.getBasicAuthentication = function() {
       var token = this.readBasicUsername() + ":" + this.readBasicPassword();
       if( token == ":" ) {
           return ""
       }
       return btoa(token); 
}

// --------------------------------------------------------------------------------------
var _process = function(text, status, xhr, fullhost) {
    var content = xhr.getResponseHeader('Content-Type');

    // Clean up previously opened result blocks
    $('#html_block').hide();
    $('#body_block').hide();

    if( content == null )
    {
        content = "text";
    }

    try {
        if( xhr.status == 0 )
        {
            $('#body_block').show();
            $('#response_body').text( "Failure while contacting API. This is usually a result of Cross-Origin Resource Sharing protection (CORS). Please check javascript domains registered against the authentication credentials you are using." );
        }
        else
        {
            if( content.match(/json/) )
            {
                $('#body_block').show();
                $('#response_body').html( hljs.highlight( 'json', JSON.stringify(JSON.parse(text), null, 2) ).value );

            } else if( content.match(/xml/) )
            {
                $('#body_block').show();
                $('#response_body').html( hljs.highlight( 'xml', text ).value );
            }
            else if( content.match(/html/) )
            {
                var iframe = $('#html_block')[0];
                var doc    = (iframe.contentWindow) ? iframe.contentWindow.document : iframe.contentDocument;

                text = text.replace( /<head>/mi, '<head><base href="//'+fullhost+'">');

                doc.open();
                doc.write( text );
                doc.close();

                $('#html_block').show();
            }
            else
            {
                $('#body_block').show();
                $('#response_body').html( hljs.highlightAuto( text ).value );
            }
        }
    }
    catch(err) {
        $('#body_block').show();
        $('#response_body').text( "Unexpected error: " + err.message + ' ' + err.line );
    }

    $('#results').show();
    $('#response').fadeIn().show();

    $('#response_code').text( xhr.status + ' ' + xhr.statusText );
    $('#response_headers').html( hljs.highlight( 'http', xhr.getAllResponseHeaders() ).value );

    $('#exploreButton').removeAttr('disabled');
}

// --------------------------------------------------------------------------------------

var _set_headers = function(request, headers ) {

    for( var i = 0; i < headers.length; i++ )
    {
        request.setRequestHeader( headers[i].name, headers[i].value );
    }
}

// --------------------------------------------------------------------------------------

var _get_header_text = function( headers ) {

    var text = '';

    for( var i = 0; i < headers.length; i++ )
    {
        text = text + '\n' + headers[i].name + ': ' + headers[i].value;
    }
    return text;
}

// --------------------------------------------------------------------------------------

var _get_url = function(url,query) {

    full_url = url + '?' + $.param(query);

    var urlp = document.createElement('a');
    urlp.href = full_url;
    // Make sure there is a leading / on the path
    urlp.pathname = urlp.pathname.replace(/(^\/?)/,"/");
    var port = '';
    if( urlp.port )
    {
        port = ':'+urlp.port;
    }

    return { fullUrl: full_url, fullhost:urlp.hostname + port, requestUrl: urlp.pathname + urlp.search };
};

// --------------------------------------------------------------------------------------
//
apiExplorer.go = function( method, url ){
    var query   = [];
    var form    = [];
    var file    = [];
    var body;
    var body_name;
    var headers = [];
    var gotjson = false;
    var errors  = [];
    var display_content_type = "";
    var content_type = "";
    var form_data = new FormData();

    $('#apiexplorer :input').each( function() {
        var $input   = $(this);
        var type     = $input.data('type');
        var required = $input.attr('required');
        var val      = $input.val(); //.trim();
        var name     = $input.prop('name');

        $input.removeClass("errorfield");

        // Pick up any missing mandatory fields
        if( required && val == "" ) {
            errors.push( $input );
        }

        var obj = { "name":name, "value":val}

        if( type=='path' ) {
            url = url.replace('{'+name+'}', val);
        }
        if( type=='query' && val ) {
            query.push( obj );
        }
        if( type=='header' && val ) {
            headers.push( obj );
        }
        if( type=='form' && val ) {
            form.push( obj );
        }
        if( type=='file' && val ) {
            obj.file = $input[0].files[0];
            file.push( obj );
        }
        if( type=='body' && val ) {
            body = val;
            // Parse might fail
            try {
                $.extend(body, JSON.parse(val));
            } catch(e) {
                errors.push( $input );
                $('#jsonerror').html( e );
                $('#jsonerror').show();
            }
            gotjson = true; // FIXME This might not be JSON! Use "produces" from swagger spec to decide JSON, XML etc
            body_name = name;
        }
    });

    // Handle errors
    if( errors.length ) {
        $.each( errors, function( index, value ) {
            // Target outer "group" with has-error class, as well as the input field. This gives a bit of flexibility
            $('#'+value.prop('name')+'-group').addClass("has-error");
            value.addClass("has-error");
            value.wiggle();
        });
        $('#results').hide();
        return;
    }

    var body_text;

    $('#request_body').hide();
    $('#jsonerror').hide();

    if( form.length )
    {
        body_text = $.param(form);
        content_type = "application/x-www-form-urlencoded";// Just for display purposes. formData will set it's own.
        display_content_type = "\nContent-Type: "+content_type;
        $('#request_body').html( hljs.highlightAuto( body_text ).value );
        $('#request_body').show();

        for( var p in form ) {
            form_data.append( form[p].name, form[p].value);
        }
    }

    if( file.length )
    {
        content_type = "application/x-www-form-urlencoded"; // Just for display purposes. formData will set it's own.
        display_content_type = "\nContent-Type: "+content_type;
        for( var p in file ) {
            form_data.append( file[p].name, file[p].file);
        }
    }

    // Create display headers before custom headers/queries are added, as these are internal.
    var display_headers = _get_header_text( headers );
    
    // Add internal queries and headers, if the extend callback is configured.
    // This allows authentication paramerters, for example, to be added to the request.
    if( this._extendCallback ) {
        var req = {};
        this._extendCallback(req); 

        if( req.headers ) {
            for( var p in req.headers ) {
                if( req.headers.hasOwnProperty(p) ) {
                    headers.push( { name: p, value:req.headers[p] } );
                }
            }
        }

        if( req.params ) {
            for( var p in req.params ) {
                if( req.params.hasOwnProperty(p) ) {
                    query.push( { name: p, value:req.params[p] } );
                }
            }
        }
    }
    var display_url = _get_url(url, query).requestUrl;

    // Tricky to see if formData object is empty. This works. Not elegant.
    var got_form_data = false;
    for(var pair of form_data.entries()) {
        got_form_data = true;
        break;
    }

    var data;
    if( gotjson ) // TODO need to handle outer mime types, such as XML, text etc
    {
        body_text = JSON.stringify(body, null, 2); 

        // If we don't have an empty JSON document, display it
        if( body_text != '{}' ) {
            $('#request_body').html( hljs.highlightAuto( body ).value );
            $('#request_body').show();
            content_type = "application/json"; // Display and actual content type

            // Only send body as multipart formData (above) IFF there is also form data or file upload.
            // Otherwise, send as playin body_text body.
            if( got_form_data ) {
                var blob = new Blob([JSON.parse(body_text)], {type: content_type});
                form_data.append(body_name, blob, body_name); // name, content, filename
                content_type = "application/x-www-form-urlencoded"; // Just for display purposes. multipart will contain real
            } else {
                data = body;
            }
            display_content_type = "\nContent-Type: "+content_type;
        }
    }

    if( got_form_data ) {
        data = form_data;
    }

    // Construct request URL bits from the url and extended query
    var constructed_request = _get_url( url, query );

    // FIXME Get protocol from passed in url
    $('#request_url').html( hljs.highlight( 'http', method.toUpperCase() + ' ' + display_url + ' HTTP/1.1\nHost: ' + constructed_request.fullhost + display_content_type + display_headers ).value );

    $('#exploreButton').attr('disabled', 'disabled');

    $.support.cors = true;

    $.ajax({
        url: constructed_request.fullUrl,
        async: true,
        data: data,
        type: method,
        dataType: "text", // Prevents .ajax from parsing response (JSON.parse in _process does this)
        contentType: got_form_data ? false : content_type, // MUST be false if we've got formData, else real type.
        traditional: true,
        processData: false, // Must be False for FormData

        success:  function( text, status, xhr)  { _process(text, status, xhr, constructed_request.fullhost) },
        error:    function( xhr,  status, text) { _process(xhr.responseText,  status, xhr, constructed_request.fullhost) },
        beforeSend: function( request ) {
            _set_headers( request, headers );
            $('#progress').stop(1,0).hide().delay(800).fadeIn();
            $('#response').stop(1,0).delay(10).hide();
        },
        complete:   function() { $('#progress').stop(1,0).hide(); }
        //,statusCode: {
        //    404: function() {
        //        alert("Not found");
        //    }
        //}
    });
}

// --------------------------------------------------------------------------------------
