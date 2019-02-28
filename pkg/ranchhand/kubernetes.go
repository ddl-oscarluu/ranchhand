package ranchhand

import (
	"os"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
)

var tpl *template.Template

const (
	RKEConfigFile = "rancher-cluster.yml"
	RKETemplate   = `# DO NOT EDIT THIS - GENERATED BY RANCHHAND
ssh_key_path: {{ .SSHKeyPath }}
ignore_docker_version: false
nodes:
{{- range .Nodes }}
  - address: {{ . }}
    user: {{ $.SSHUser }}
    port: {{ $.SSHPort }}
    role: [controlplane,worker,etcd]
{{- end }}

services:
  etcd:
    snapshot: true
    creation: 6h
    retention: 24h
`
)

func installKubernetes(cfg *Config) error {
	// todo: add check to skip when already installed

	// generate rke config
	file, err := os.Create(RKEConfigFile)
	if err != nil {
		return errors.Wrapf(err, "cannot create %s", RKEConfigFile)
	}
	defer file.Close()

	if err := tpl.Execute(file, cfg); err != nil {
		return errors.Wrap(err, "rke template render failed")
	}

	// execute rke up
	cmd := exec.Command("rke", "up", "--config", RKEConfigFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "cannot install kubernetes")
	}

	// todo: add cluster-ready check

	return nil
}

func init() {
	tpl = template.Must(template.New("rke-tmpl").Parse(RKETemplate))
}
