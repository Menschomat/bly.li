#!/bin/bash

# Install oapi-codegen
echo "Installing oapi-codegen..."
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Function to generate API code for a service
generate_api_code() {
  local service_name="$1"
  local service_dir="src/services/$service_name"
  local openapi_spec="src/api/$service_name.openapi.yml"
  local output_file="$service_dir/api/api.gen.go"

  if [ -f "$openapi_spec" ]; then
    echo "Generating API code for $service_dir from $openapi_spec..."
    oapi-codegen -generate types,chi-server -import-mapping ./shared.openapi.yml:github.com/Menschomat/bly.li/shared/api -package api -o "$output_file" "$openapi_spec"
    echo "API code generated at $output_file"
  else
    echo "No OpenAPI spec found at $openapi_spec. Skipping generation for $service_dir."
  fi
}
oapi-codegen -generate types,skip-prune -package api -o "src/shared/api/api.gen.go" "src/api/shared.openapi.yml"
# Generate API code for each service
generate_api_code "blowup"
generate_api_code "shortn"
generate_api_code "dasher"
