#!/bin/bash

# Install oapi-codegen
echo "Installing oapi-codegen..."
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Function to generate API code for a service
generate_api_code() {
  local service_dir="$1"
  local openapi_spec="$service_dir/openapi.yml"
  local output_file="$service_dir/api.gen.go"

  if [ -f "$openapi_spec" ]; then
    echo "Generating API code for $service_dir from $openapi_spec..."
    oapi-codegen -generate types,chi-server -package api -o "$output_file" "$openapi_spec"
    echo "API code generated at $output_file"
  else
    echo "No OpenAPI spec found at $openapi_spec. Skipping generation for $service_dir."
  fi
}

# Generate API code for each service
generate_api_code "src/services/blowup/api"
generate_api_code "src/services/shortn/api"
