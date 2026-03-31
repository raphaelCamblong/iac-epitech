resource "helm_release" "ingress_nginx" {
  name             = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = var.ingress_controller_namespace
  create_namespace = true
  version          = var.ingress_nginx_chart_version
  wait             = true
  timeout          = 600

  set {
    name  = "controller.ingressClass"
    value = "nginx"
  }

  set {
    name  = "controller.ingressClassResource.name"
    value = "nginx"
  }

  dynamic "set" {
    for_each = var.ingress_static_ip != "" ? [var.ingress_static_ip] : []

    content {
      name  = "controller.service.loadBalancerIP"
      value = set.value
    }
  }
}

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
    "DATABASE_URL" = var.database_url
    "JWT_SECRET"   = var.jwt_secret
  }

  type = "Opaque"

  depends_on = [kubernetes_namespace.task_manager]
}

resource "helm_release" "task_manager" {
  name      = var.release_name
  chart     = "${path.module}/../../charts/task-manager" # Adjust path relative to the module
  namespace = var.namespace
  wait      = true
  timeout   = 600

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

  depends_on = [
    helm_release.ingress_nginx,
    kubernetes_secret.app_secrets,
  ]
}
