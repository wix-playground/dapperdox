export PROXY_PATH=/developer=https://developer.companieshouse.gov.uk,/fred/=https://google.co.uk

./swaggerly \
    -assets-dir=./examples/overlay/assets \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://0.0.0.0:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -force-specification-list=true \
    -author-show-assets=true \
    -theme=default \
    -tls-certificate=server.rsa.crt \
    -tls-key=server.rsa.key \
    -log-level=info
    #-proxy-path=/developer=https://developer.companieshouse.gov.uk \
    #-proxy-path=/fred/=https://google.co.uk \
