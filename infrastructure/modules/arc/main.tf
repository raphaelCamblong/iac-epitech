resource "helm_release" "gha_runner_scale_set_controller" {
  name             = "arc-systems"
  repository       = "oci://ghcr.io/actions/actions-runner-controller-charts"
  chart            = "gha-runner-scale-set-controller"
  namespace        = "arc-systems"
  create_namespace = true
}

resource "helm_release" "gha_runner_scale_set" {
  name             = "arc-runner-set"
  repository       = "oci://ghcr.io/actions/actions-runner-controller-charts"
  chart            = "gha-runner-scale-set"
  namespace        = "arc-runners"
  create_namespace = true

  # GitHub's REST API rejects repo paths that include a trailing ".git" (404 on registration-token).
  set {
    name  = "githubConfigUrl"
    value = trimsuffix(var.github_repo_url, ".git")
  }

  set_sensitive {
    name  = "githubConfigSecret.github_token"
    value = var.github_pat
  }

  depends_on = [
    helm_release.gha_runner_scale_set_controller
  ]
}
