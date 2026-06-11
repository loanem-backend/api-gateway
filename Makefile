bundle-openapi:
	vacuum bundle docs/openapi.yaml docs/openapi_bundled.yaml

mock-server:
	prism mock docs/openapi_bundled.yaml

gen-mock-client:
	mockgen \
	-destination=internal/mocks/server/mock_$(service)_client.go \
	-package=server_mock \
	github.com/loanem-backend/protos/pb/proto/services/$(svc)/v1 \
	$(interface)
