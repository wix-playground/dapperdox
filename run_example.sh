./dapperdox \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://localhost:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -theme=default  \
    -log-level=trace \
    #-tls-certificate=server.rsa.crt \
    #-tls-key=server.rsa.key \
    #-author-show-assets=true \
    #-force-specification-list=true \
    #-assets-dir=./examples/overlay/assets \
    #-proxy-path=/developer=https://developer.some-dev-site.com 
