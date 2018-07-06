package cmd

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "Chart.yaml",
		FileModTime: time.Unix(1530546712, 0),
		Content:     string("apiVersion: v1\ndescription: kube-sre-essentials\nname: kube-sre-essentials\nversion: 0.3.0\n"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "requirements.yaml",
		FileModTime: time.Unix(1530801058, 0),
		Content:     string("dependencies:\n- name: nginx-ingress\n  version: 0.22.1\n  repository: \"https://kubernetes-charts.storage.googleapis.com/\"\n- name: prometheus\n  version: 6.8.0\n  repository: \"https://kubernetes-charts.storage.googleapis.com/\""),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "secret-ingress-auth-prometheus.yml",
		FileModTime: time.Unix(1530861369, 0),
		Content:     string("apiVersion: v1\nkind: Secret\nmetadata:\n  name: ingress-auth-prometheus\ntype: Opaque\ndata:\n  auth: \"{{ .Prometheus.Ingress.Auth }}\""),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "values.yml",
		FileModTime: time.Unix(1530821548, 0),
		Content:     string("nginx-ingress:\n  controller:\n    hostNetwork: true\n    service:\n      externalIPs:\n        - {{ .Ingress.ExternalIp }}\n    nodeSelector:\n      kubernetes.io/hostname: {{ .Ingress.NodeSelector }}\nprometheus:\n  alertmanager:\n    persistentVolume:\n      enabled: False\n  server:\n    persistentVolume:\n      enabled: False\n    ingress:\n      enabled: True\n      annotations:\n        kubernetes.io/ingress.class: nginx\n        nginx.ingress.kubernetes.io/auth-type: basic\n        nginx.ingress.kubernetes.io/auth-secret: ingress-auth-prometheus\n        nginx.ingress.kubernetes.io/auth-realm: \"Authentication required to access Prometheus.\"\n      hosts:\n        - {{ .Prometheus.Host }}\n      tls:\n      - hosts:\n        - {{ .Prometheus.Host }}\n        secretName: {{ .Prometheus.Host }}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1530887324, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "Chart.yaml"
			file3, // "requirements.yaml"
			file4, // "secret-ingress-auth-prometheus.yml"
			file5, // "values.yml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1530887324, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"Chart.yaml":                         file2,
			"requirements.yaml":                  file3,
			"secret-ingress-auth-prometheus.yml": file4,
			"values.yml":                         file5,
		},
	})
}
