./dapperdox \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://localhost:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -theme=dapperdox-theme-gov-uk  \
    -theme-dir=../../companieshouse \
    -log-level=info \
    -force-specification-list=true \
    -author-show-assets=true \
    #-assets-dir=./examples/guides/assets \
    #-assets-dir=./examples/gds/assets \
    #-tls-certificate=server.rsa.crt \
    #-tls-key=server.rsa.key \
    #-proxy-path=/developer=https://developer.some-dev-site.com 
