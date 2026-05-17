# Task Manager API

Production-ready Task Manager REST API built with Go, Clean Architecture, PostgreSQL, and deployed via Helm/Terraform.

## Requirements

- Go 1.23+
- PostgreSQL (CloudSQL, RDS, or local)
- Docker (for building image)
- Kubernetes cluster with Helm
- Terraform 1.9+ (for deployment; CI pins 1.9.8)


# [Deployment instructions](./DEPLOYMENTS.md)

# [Application](./app/README.md)

## Architecture diagrams

### Dev environment — zonal GKE in `europe-west9-a`

![Dev environment infrastructure](./docs/infra-dev.png)

<sub>Source: [`docs/infra-dev.html`](./docs/infra-dev.html)</sub>

### Prod environment — regional GKE across 3 zones

![Prod environment infrastructure](./docs/infra-prod.png)

<sub>Source: [`docs/infra-prod.html`](./docs/infra-prod.html)</sub>

### CI/CD pipeline — App CI → Infra CD via `workflow_run`

![CI/CD pipeline](./docs/cicd-pipeline.png)

<sub>Source: [`docs/cicd-pipeline.html`](./docs/cicd-pipeline.html)</sub>
