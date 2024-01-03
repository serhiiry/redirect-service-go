# Go-based Redirect Service

## Overview
This repository contains a Go-based HTTP redirect service. 
The service redirects incoming HTTP requests to various destinations based on domain pools, maintaining URL parameters and structure. 
It also includes path-based domain routing and custom header functionality.

## Getting Started

### Prerequisites
- Docker
- Go (optional, for local development without Docker)

### Build and Run

#### Using Docker
1. **Build the Docker image**:

```bash
make build
```
This command builds the Docker image for the redirect service.

2. **Run the Docker container**:

```bash
make run
```

This command runs the Docker container, exposing the service on port 80.

### Configuration
The service is configured through a JSON file (`redirect-config.json`), where you can define domain pools, path-based domains, and custom headers.

## Testing the Service
Once the service is running, you can test the redirects by sending HTTP requests to `http://localhost/redirect/{pool_id}/{path}`.

## Health Check
The service includes a `/health/` endpoint for health checks, accessible at `http://localhost/health/`.
