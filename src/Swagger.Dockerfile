# Use the official Swagger UI Docker image
FROM swaggerapi/swagger-ui

# Create the directory for OpenAPI specs
RUN mkdir -p /usr/share/nginx/html/openapi

# Copy the OpenAPI specifications
COPY services/shortn/api/openapi.yml /usr/share/nginx/html/openapi/shortn.yml
COPY services/blowup/api/openapi.yml /usr/share/nginx/html/openapi/blowup.yml

# Set environment variables to configure Swagger UI
ENV URLS='[{"url":"./openapi/shortn.yml","name":"Shortn API"},{"url":"./openapi/blowup.yml","name":"Blowup API"}]'
ENV URLS_PRIMARY_NAME='Shortn API'
ENV SWAGGER_UI_PRESORT=true
ENV DEEP_LINKING=true
ENV LAYOUT='StandaloneLayout'
