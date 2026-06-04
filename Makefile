bundle-openapi:
	vacuum bundle docs/openapi.yaml docs/openapi_bundled.yaml

mock-server:
	prism mock docs/openapi_bundled.yaml