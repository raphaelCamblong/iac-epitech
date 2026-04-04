#!/usr/bin/env bash
# Creates a GitHub Actions service account, Workload Identity pool/provider,
# and binds the repo so Actions can authenticate to GCP without a JSON key.
#
# Usage (environment variables):
#   export YOUR_PROJECT_ID="your-gcp-project-id"
#   export YOUR_GITHUB_ORG="your-org-or-user"
#   export YOUR_REPO="your-repo-name"
#   ./scripts/gcp/setup-github-actions-wif.sh
#
# Usage (positional parameters — same order as above):
#   ./scripts/gcp/setup-github-actions-wif.sh your-gcp-project-id your-org-or-user your-repo-name

set -euo pipefail

if [[ "${#}" -eq 3 ]]; then
  YOUR_PROJECT_ID="${1}"
  YOUR_GITHUB_ORG="${2}"
  YOUR_REPO="${3}"
fi

if [[ -z "${YOUR_PROJECT_ID:-}" || -z "${YOUR_GITHUB_ORG:-}" || -z "${YOUR_REPO:-}" ]]; then
  echo "Missing configuration. Set YOUR_PROJECT_ID, YOUR_GITHUB_ORG, and YOUR_REPO," >&2
  echo "or pass them as three arguments: <project_id> <github_org> <repo>" >&2
  exit 1
fi

echo "Using project: ${YOUR_PROJECT_ID}, repository: ${YOUR_GITHUB_ORG}/${YOUR_REPO}"

# 1. Service account GitHub Actions will impersonate
gcloud iam service-accounts create "github-actions-sa" \
  --project="${YOUR_PROJECT_ID}" \
  --display-name="GitHub Actions Service Account"

# 2. Permissions for GKE, SQL, networking, etc.
gcloud projects add-iam-policy-binding "${YOUR_PROJECT_ID}" \
  --member="serviceAccount:github-actions-sa@${YOUR_PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/owner"

# 3. Workload Identity pool
gcloud iam workload-identity-pools create "github-actions-pool" \
  --project="${YOUR_PROJECT_ID}" \
  --location="global" \
  --display-name="GitHub Actions Pool"

# 4. OIDC provider (trust between Google Cloud and GitHub)
gcloud iam workload-identity-pools providers create-oidc "github-actions-provider" \
  --project="${YOUR_PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --display-name="GitHub Actions Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository == '${YOUR_GITHUB_ORG}/${YOUR_REPO}'" \
  --issuer-uri="https://token.actions.githubusercontent.com"

# 5. Project number (needed for the workload identity principal)
PROJECT_NUMBER="$(gcloud projects describe "${YOUR_PROJECT_ID}" --format="value(projectNumber)")"
echo "Your Project Number is: ${PROJECT_NUMBER}"

# 6. Allow this GitHub repo to obtain credentials for the service account
gcloud iam service-accounts add-iam-policy-binding "github-actions-sa@${YOUR_PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${YOUR_PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/${YOUR_GITHUB_ORG}/${YOUR_REPO}"
