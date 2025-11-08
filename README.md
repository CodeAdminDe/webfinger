# Webfinger Go Server

A simple Go application that implements a Webfinger server. This server can be configured to respond to specific `acct` resources and supports domain-based validation for more flexible handling.

## Features

-   **Webfinger Protocol Implementation**: Serves JSON Resource Descriptors (JRD) for Webfinger requests.
-   **Configurable Resource and Issuer URL**: Easily set the target `acct` resource and issuer URL via environment variables.
-   **Domain-Based `acct` Validation**: Optionally allows validation of any username within a configured domain, providing flexibility for multi-user domains.

## Configuration

The server is configured using environment variables:

-   `WEBFINGER_RESOURCE`: The primary `acct` resource (e.g., `acct:user@example.com`) that the server is configured for. This is used for exact matching or domain extraction.
-   `WEBFINGER_ISSUER_URL`: The URL of the OpenID Connect issuer to be included in the JRD response.
-   `WEBFINGER_ALLOW_DOMAIN_WILDCARD`: Set to `true` to enable domain-based validation. If `true`, the server will respond to any `acct` resource with a matching domain to `WEBFINGER_RESOURCE`. If `false` or unset, only an exact match to `WEBFINGER_RESOURCE` will be accepted.

## How to build and run from src

1.  **Build the application**:
    ```bash
    cd ./src && go build -o webfinger-server .
    ```
2.  **Set environment variables**:
    ```bash
    export WEBFINGER_RESOURCE="acct:anyuser@yourdomain.com"
    export WEBFINGER_ISSUER_URL="https://your-issuer.com"
    export WEBFINGER_ALLOW_DOMAIN_WILDCARD="true" # Or "false" or omit for exact match
    ```
3.  **Run the server**:
    ```bash
    ./webfinger-server
    ```

The server will listen on port 8080.

## How to build and run - dockerized

1.  **Build the Docker image**:
    ```bash
    docker build -t webfinger-server .
    ```
2.  **Run the Docker container**:
    ```bash
    docker run -p 8080:8080 \
        -e WEBFINGER_RESOURCE="acct:anyuser@yourdomain.com" \
        -e WEBFINGER_ISSUER_URL="https://your-issuer.com/issuer" \
        -e WEBFINGER_ALLOW_DOMAIN_WILDCARD="true" \
        webfinger-server
    ```

## License notice

(c) 2025 Frederic Roggon
This project is licensed under the terms of GNU AFFERO GENERAL PUBLIC (LICENSE)[./LICENSE].
You should have received a copy of the GNU AGPL with the program and/or source files.
Otherwise, see (https://github.com/CodeAdminDe/webfinger/blob/main/LICENSE)[https://github.com/CodeAdminDe/webfinger/blob/main/LICENSE].