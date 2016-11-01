./swaggerly \
    -assets-dir=./examples/overlay/assets \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://0.0.0.0:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -force-specification-list=false \
    -author-show-assets=false \
    -log-level=info
