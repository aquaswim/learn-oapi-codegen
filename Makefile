OAPI_CODEGEN_VERSION=latest
OAPI_GENERATE_SERVER_CFG=api/generate-server.cfg.yaml
OAPI_GENERATE_MODEL_CFG=api/generate-model.cfg.yaml
OAPI_SPEC=api/contract.yaml

build: generate
	go build -o api-server cmd/api-server/main.go

.PHONY: __install_deps clean generate

__install_deps:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}

generate: __install_deps
	oapi-codegen -config ${OAPI_GENERATE_MODEL_CFG} ${OAPI_SPEC}
	oapi-codegen -config ${OAPI_GENERATE_SERVER_CFG} ${OAPI_SPEC}

clean:
	@echo "cleaning generated file"
	@rm -fv internal/api_server/api_server.gen.go
	@rm -fv internal/dto/model.gen.go
	@rm -fv api-server
