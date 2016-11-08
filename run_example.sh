./swaggerly \
    -assets-dir=./examples/guides/assets \
    -spec-dir=examples/specifications/petstore/ \
    -bind-addr=0.0.0.0:3123 \
    -site-url=http://0.0.0.0:3123 \
    -spec-rewrite-url=petstore.swagger.io=PETSTORE.swagger.io \
    -force-specification-list=true \
    -author-show-assets=false \
    -theme=default \
    -log-level=info
