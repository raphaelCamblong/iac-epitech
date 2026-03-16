resource "kubernetes_namespace" "task_manager" {
  count = var.namespace == "default" ? 0 : 1

  metadata {
    name = var.namespace
  }
}

resource "kubernetes_secret" "app_secrets" {
  metadata {
    name      = "${var.release_name}-secrets"
    namespace = var.namespace
  }

  data = {
    "DATABASE_URL" = base64encode(var.database_url)
    "JWT_SECRET"   = base64encode(var.jwt_secret)
  }

  type = "Opaque"
}

resource "helm_release" "task_manager" {
  name      = var.release_name
  chart     = "${path.module}/../charts/task-manager"
  namespace = var.namespace

  set {
    name  = "existingSecret"
    value = kubernetes_secret.app_secrets.metadata[0].name
  }

  set {
    name  = "image.repository"
    value = var.image_repository
  }

  set {
    name  = "image.tag"
    value = var.image_tag
  }

  set {
    name  = "replicaCount"
    value = var.replica_count
  }

  set {
    name  = "ingress.hosts[0].host"
    value = var.ingress_host
  }

  set {
    name  = "hpa.minReplicas"
    value = var.hpa_min_replicas
  }

  set {
    name  = "hpa.maxReplicas"
    value = var.hpa_max_replicas
  }

  depends_on = [kubernetes_secret.app_secrets]
}
