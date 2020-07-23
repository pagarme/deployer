job "{{ .Environment.name }}-nginx" {
  type        = "system"
  region      = "global"
  datacenters = ["dc1"]

  constraint {
    attribute = "${node.class}"
    value     = "nginx"
  }

  constraint {
    attribute = "${meta.environment}"
    value     = "{{ .Environment.name }}"
  }

  update {
    stagger      = "5s"
    max_parallel = 1
  }

  group "web" {
    task "nginx" {
      driver = "docker"

      config {
        image      = "{{ .Image }}"
        force_pull = true

        port_map {
          "http" = 80
        }

        dns_servers = ["${attr.unique.network.ip-address}"]
      }

      env {
          NGINX_PORT = "80"
      }

      service {
        name = "nginx"

        tags = ["nginx", "${meta.environment}"]

        port = "http"

        check {
          name = "NGINX Health Check"
          type = "http"
          port = "http"
          path = "/"
          interval = "5s"
          timeout = "2s"
        }
      }

      resources {
        cpu    = "2400"
        memory = "1024"

        network {
          mbits = 100

          port "http" {
            static = 80
          }
        }
      }
    }
  }
}

