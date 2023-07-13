package template

const MasterTemplate = `---
ignition:
  version: "2.2.0"
passwd:
  users:
{{ range $index, $user := .Cluster.Kubernetes.SSH.UserList }}
    - name: {{ $user.Name }}
      shell: "/bin/bash"
      groups:
        - "sudo"
        - "docker"
{{ if ne $user.PublicKey "" }}
      sshAuthorizedKeys:
        - "{{ $user.PublicKey }}"
{{ end }}
{{ end }}

systemd:
  units:
  # Start - manual management for cgroup structure
  - name: kubereserved.slice
    path: /etc/systemd/system/kubereserved.slice
    content: |
      [Unit]
      Description=Limited resources slice for Kubernetes services
      Documentation=man:systemd.special(7)
      DefaultDependencies=no
      Before=slices.target
      Requires=-.slice
      After=-.slice
  # End - manual management for cgroup structure
  - name: audit-rules.service
    enabled: true
    dropins:
    - name: 10-Wait-For-Docker.conf
      contents: |
        [Service]
        ExecStartPre=/bin/bash -c "while [ ! -f /etc/audit/rules.d/10-docker.rules ]; do echo 'Waiting for /etc/audit/rules.d/10-docker.rules to be written' && sleep 1; done"
  - name: auditd.service
    enabled: true
  {{range .Extension.Units}}
  - name: {{.Metadata.Name}}
    enabled: {{.Metadata.Enabled}}
    contents: |
      {{range .Content}}{{.}}
      {{end}}{{end}}
  - name: set-certs-group-owner-permission-giantswarm.service
    enabled: true
    contents: |
      [Unit]
      Description=Change group owner for certificates to giantswarm
      Wants=k8s-kubelet.service k8s-setup-network-env.service
      After=k8s-kubelet.service k8s-setup-network-env.service
      [Service]
      Type=oneshot
      ExecStart=/bin/sh -c "find /etc/kubernetes/ssl -type f -print | xargs -i  sh -c 'chown root:giantswarm {} && chmod 640 {}'"
      [Install]
      WantedBy=multi-user.target
  - name: wait-for-domains.service
    enabled: true
    contents: |
      [Unit]
      Description=Wait for etcd and k8s API domains to be available
      [Service]
      Type=oneshot
      ExecStart=/opt/wait-for-domains
      [Install]
      WantedBy=multi-user.target
  - name: os-hardening.service
    enabled: true
    contents: |
      [Unit]
      Description=Apply os hardening
      [Service]
      Type=oneshot
      ExecStartPre=-/bin/bash -c "gpasswd -d core rkt; gpasswd -d core docker; gpasswd -d core wheel"
      ExecStartPre=/bin/bash -c "until [ -f '/etc/sysctl.d/hardening.conf' ]; do echo Waiting for sysctl file; sleep 1s;done;"
      ExecStart=/usr/sbin/sysctl -p /etc/sysctl.d/hardening.conf
      [Install]
      WantedBy=multi-user.target
  - name: k8s-setup-kubelet-environment.service
    enabled: true
    contents: |
      [Unit]
      Description=k8s-setup-kubelet-environment Service
      After=k8s-setup-network-env.service docker.service
      Requires=k8s-setup-network-env.service docker.service
      [Service]
      Type=oneshot
      RemainAfterExit=yes
      TimeoutStartSec=0
      ExecStart=/opt/bin/setup-kubelet-environment master
      [Install]
      WantedBy=multi-user.target
  - name: k8s-setup-kubelet-config.service
    enabled: true
    contents: |
      [Unit]
      Description=k8s-setup-kubelet-config Service
      After=k8s-setup-network-env.service docker.service k8s-setup-kubelet-environment.service
      Requires=k8s-setup-network-env.service docker.service k8s-setup-kubelet-environment.service
      [Service]
      Type=oneshot
      RemainAfterExit=yes
      TimeoutStartSec=0
      Environment=IMAGE={{ .Images.Envsubst }}
      ExecStart=docker run --rm \
        --env-file /etc/network-environment --env-file /etc/kubelet-environment \
        -v /etc/kubernetes/config/:/etc/kubernetes/config/ \
        $IMAGE \
        ash -c "cat /etc/kubernetes/config/kubelet.yaml.tmpl |envsubst >/etc/kubernetes/config/kubelet.yaml"
      [Install]
      WantedBy=multi-user.target
  - name: k8s-setup-apiserver-manifest.service
    enabled: true
    contents: |
      [Unit]
      Description=k8s-setup-apiserver-manifest Service
      After=docker.service
      Requires=docker.service
      [Service]
      Type=oneshot
      RemainAfterExit=yes
      TimeoutStartSec=0
      Environment=IMAGE={{ .Images.Envsubst }}
      ExecStartPre=/opt/bin/setup-apiserver-environment
      ExecStart=docker run --rm \
        --env-file /etc/apiserver-environment \
        -v /etc/kubernetes/:/etc/kubernetes/ \
        $IMAGE \
        ash -c "cat /etc/kubernetes/manifest-templates/k8s-api-server.yaml.tmpl |envsubst >/etc/kubernetes/manifests/k8s-api-server.yaml"
      [Install]
      WantedBy=multi-user.target
  - name: containerd.service
    enabled: true
    contents: |
    dropins:
      - name: 10-change-cgroup.conf
        contents: |
          [Service]
          CPUAccounting=true
          MemoryAccounting=true
          Slice=kubereserved.slice
  - name: docker.service
    enabled: true
    contents: |
    dropins:
      - name: 10-giantswarm-extra-args.conf
        contents: |
          [Service]
          CPUAccounting=true
          MemoryAccounting=true
          Slice=kubereserved.slice
          Environment="DOCKER_CGROUPS={{ if .ForceCGroupsV1 }}--exec-opt native.cgroupdriver=cgroupfs --cgroup-parent=/kubereserved.slice {{ else }}--cgroup-parent=kubereserved.slice {{ end }}--log-opt max-size=25m --log-opt max-file=2 --log-opt labels=io.kubernetes.container.hash,io.kubernetes.container.name,io.kubernetes.pod.name,io.kubernetes.pod.namespace,io.kubernetes.pod.uid"
          Environment="DOCKER_OPT_BIP=--bip={{.Cluster.Docker.Daemon.CIDR}}"
          {{- if .Proxy.HTTP }}
          Environment="HTTP_PROXY={{ .Proxy.HTTP }}"
          {{- end }}
          {{- if .Proxy.HTTPS }}
          Environment="HTTPS_PROXY={{ .Proxy.HTTPS }}"
          {{- end }}
          {{- if .Proxy.NoProxy }}
          Environment="NO_PROXY={{ .Proxy.NoProxy }}"
          {{- end }}
  - name: k8s-setup-network-env.service
    enabled: true
    contents: |
      [Unit]
      Description=k8s-setup-network-env Service
      Wants=network.target docker.service wait-for-domains.service
      After=network.target docker.service wait-for-domains.service
      [Service]
      Type=oneshot
      TimeoutStartSec=0
      Environment="IMAGE={{ .Images.KubernetesNetworkSetupDocker }}"
      Environment="NAME=%p.service"
      ExecStartPre=/usr/bin/mkdir -p /opt/bin/
      ExecStartPre=/usr/bin/docker pull $IMAGE
      ExecStartPre=-/usr/bin/docker stop -t 10 $NAME
      ExecStartPre=-/usr/bin/docker rm -f $NAME
      ExecStart=/usr/bin/docker run --rm --net=host -v /etc:/etc --name $NAME $IMAGE
      ExecStop=-/usr/bin/docker stop -t 10 $NAME
      ExecStopPost=-/usr/bin/docker rm -f $NAME
      [Install]
      WantedBy=multi-user.target
  - name: etcd3.service
    enabled: true
    contents: |
      [Unit]
      Description=etcd3
      Wants=k8s-setup-network-env.service
      After=k8s-setup-network-env.service
      Conflicts=etcd.service etcd2.service
      StartLimitIntervalSec=0
      [Service]
      Restart=always
      RestartSec=0
      TimeoutStopSec=10
      LimitNOFILE=40000
      CPUAccounting=true
      MemoryAccounting=true
      Slice=kubereserved.slice
      Environment=IMAGE={{ .Images.Etcd }}
      Environment=NAME=%p.service
      Environment=ETCD_NAME={{ .Etcd.NodeName }}
      Environment=ETCD_INITIAL_CLUSTER={{ .Etcd.InitialCluster }}
      Environment=ETCD_INITIAL_CLUSTER_STATE={{ .Etcd.InitialClusterState }}
      Environment=ETCD_PEER_CA_PATH=/etc/etcd/server-ca.pem
      Environment=ETCD_PEER_CERT_PATH=/etc/etcd/server-crt.pem
      Environment=ETCD_PEER_KEY_PATH=/etc/etcd/server-key.pem
      EnvironmentFile=/etc/network-environment
      ExecStartPre=-/usr/bin/docker stop  $NAME
      ExecStartPre=-/usr/bin/docker rm  $NAME
      ExecStartPre=-/usr/bin/docker pull $IMAGE
      ExecStartPre=/bin/bash -c "while [ ! -f /etc/kubernetes/ssl/etcd/server-ca.pem ]; do echo 'Waiting for /etc/kubernetes/ssl/etcd/server-ca.pem to be written' && sleep 1; done"
      ExecStartPre=/bin/bash -c "while [ ! -f /etc/kubernetes/ssl/etcd/server-crt.pem ]; do echo 'Waiting for /etc/kubernetes/ssl/etcd/server-crt.pem to be written' && sleep 1; done"
      ExecStartPre=/bin/bash -c "while [ ! -f /etc/kubernetes/ssl/etcd/server-key.pem ]; do echo 'Waiting for /etc/kubernetes/ssl/etcd/server-key.pem to be written' && sleep 1; done"
      ExecStartPre=/bin/bash -c "chmod 0700 /var/lib/etcd/"
      ExecStart=/usr/bin/docker run \
          -v /etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt \
          -v /etc/kubernetes/ssl/etcd/:/etc/etcd \
          -v /var/lib/etcd/:/var/lib/etcd  \
          --net=host  \
          --name $NAME \
          $IMAGE \
          etcd \
          --name ${ETCD_NAME} \
          --trusted-ca-file /etc/etcd/server-ca.pem \
          --cert-file /etc/etcd/server-crt.pem \
          --key-file /etc/etcd/server-key.pem\
          --client-cert-auth=true \
          --peer-trusted-ca-file ${ETCD_PEER_CA_PATH} \
          --peer-cert-file ${ETCD_PEER_CERT_PATH} \
          --peer-key-file ${ETCD_PEER_KEY_PATH} \
          --peer-client-cert-auth=true \
          --advertise-client-urls=https://{{ .Cluster.Etcd.Domain }}:{{ .Etcd.ClientPort }} \
          --initial-advertise-peer-urls=https://${ETCD_NAME}.{{ .BaseDomain }}:2380 \
          --listen-client-urls=https://0.0.0.0:2379 \
          --listen-peer-urls=https://0.0.0.0:2380 \
          --initial-cluster-token k8s-etcd-cluster \
          --initial-cluster ${ETCD_INITIAL_CLUSTER} \
          --initial-cluster-state ${ETCD_INITIAL_CLUSTER_STATE} \
          --experimental-peer-skip-client-san-verification=true \
          --data-dir=/var/lib/etcd \
          --quota-backend-bytes {{ if eq .Etcd.QuotaBackendBytes 0}}8589934592{{ else }}{{ .Etcd.QuotaBackendBytes }}{{ end }} \
          --auto-compaction-mode revision \
          --auto-compaction-retention 1 \
          --logger=zap
      [Install]
      WantedBy=multi-user.target
  - name: etcd3-defragmentation.service
    enabled: false
    contents: |
      [Unit]
      Description=etcd defragmentation job
      After=docker.service etcd3.service
      Requires=docker.service etcd3.service
      [Service]
      Type=oneshot
      EnvironmentFile=/etc/network-environment
      Environment=IMAGE={{ .Images.Etcd }}
      Environment=NAME=%p.service
      ExecStartPre=-/usr/bin/docker stop  $NAME
      ExecStartPre=-/usr/bin/docker rm  $NAME
      ExecStartPre=-/usr/bin/docker pull $IMAGE
      ExecStart=/usr/bin/docker run \
        -v /etc/kubernetes/ssl/etcd/:/etc/etcd \
        --net=host  \
        -e ETCDCTL_API=3 \
        --name $NAME \
        $IMAGE \
        etcdctl \
        --endpoints https://127.0.0.1:2379 \
        --cacert /etc/etcd/server-ca.pem \
        --cert /etc/etcd/server-crt.pem \
        --key /etc/etcd/server-key.pem \
        defrag \
        --command-timeout=60s \
        --dial-timeout=60s \
        --keepalive-timeout=25s
      [Install]
      WantedBy=multi-user.target
  - name: etcd3-defragmentation.timer
    enabled: true
    contents: |
      [Unit]
      Description=Execute etcd3-defragmentation every day at 3.30AM UTC
      [Timer]
      OnCalendar=hourly
      [Install]
      WantedBy=multi-user.target
  - name: k8s-extract.service
    enabled: true
    contents: |
      [Unit]
      Description=k8s-extract Service
      After=docker.service
      Requires=docker.service
      [Service]
      Type=oneshot
      RemainAfterExit=yes
      TimeoutStartSec=0
      Environment=IMAGE={{ .Images.Hyperkube }}
      Environment=CONTAINER_NAME=%p.service
      ExecStartPre=/usr/bin/mkdir -p /opt/bin/
      ExecStartPre=/usr/bin/docker pull $IMAGE
      ExecStartPre=-/usr/bin/docker rm $CONTAINER_NAME
      ExecStartPre=-/usr/bin/docker create --name $CONTAINER_NAME $IMAGE /kubectl
      ExecStart=/opt/k8s-extract $CONTAINER_NAME
      ExecStopPost=-/usr/bin/docker rm $CONTAINER_NAME
      [Install]
      WantedBy=multi-user.target
  - name: k8s-kubelet.service
    enabled: true
    contents: |
      [Unit]
      Wants=k8s-setup-network-env.service k8s-setup-kubelet-config.service k8s-extract.service{{ if eq .Cluster.Kubernetes.CloudProvider "" }} rpc-statd.service{{ end }}
      After=k8s-setup-network-env.service k8s-setup-kubelet-config.service k8s-extract.service{{ if eq .Cluster.Kubernetes.CloudProvider "" }} rpc-statd.service{{ end }}
      Description=k8s-kubelet
      StartLimitIntervalSec=0
      [Service]
      User=root
      TimeoutStartSec=300
      Restart=always
      RestartSec=0
      TimeoutStopSec=10
      Slice=kubereserved.slice
      CPUAccounting=true
      MemoryAccounting=true
      Environment="ETCD_CA_CERT_FILE=/etc/kubernetes/ssl/calico/etcd-ca"
      Environment="ETCD_CERT_FILE=/etc/kubernetes/ssl/calico/etcd-cert"
      Environment="ETCD_KEY_FILE=/etc/kubernetes/ssl/calico/etcd-key"
      EnvironmentFile=/etc/network-environment
      Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin"
      {{- if eq .Cluster.Kubernetes.CloudProvider "azure" }}
      ExecStartPre=/bin/bash -c 'while (curl -s -H Metadata:true --noproxy "*" "http://169.254.169.254/metadata/instance?api-version=2021-02-01" | jq -r .compute.osProfile.computerName >/etc/desired-host-name) && DES="$(cat /etc/desired-host-name)" && [ "$DES" != "" ] && HN="$(hostname)" && [ "$HN" != "$DES" ] ;  do sleep 2s ; echo "hostname is unexpected (want $DES, got $HN)" ;done;'
      {{- end }}
      ExecStart=/opt/bin/kubelet \
        {{ range .Kubernetes.Kubelet.CommandExtraArgs -}}
        {{ . }} \
        {{ end -}}
        --node-ip=${DEFAULT_IPV4} \
        --config=/etc/kubernetes/config/kubelet.yaml \
        --container-runtime=remote \
        --container-runtime-endpoint=unix:///run/containerd/containerd.sock \
        --logtostderr=true \
        {{ if .ExternalCloudControllerManager -}}
        --cloud-provider=external \
        {{ else -}}
        --cloud-provider={{.Cluster.Kubernetes.CloudProvider}} \
        {{ end -}}
        --pod-infra-container-image={{ .Images.Pause }} \
        --register-node=true \
        --register-with-taints=node-role.kubernetes.io/master=:NoSchedule \
        --kubeconfig=/etc/kubernetes/kubeconfig/kubelet.yaml \
        --node-labels="node.kubernetes.io/master,role=master,ip=${DEFAULT_IPV4},{{.Cluster.Kubernetes.Kubelet.Labels}}" \
        --v=2
      [Install]
      WantedBy=multi-user.target
  - name: etcd2.service
    enabled: false
    mask: true
  - name: update-engine.service
    enabled: false
    mask: true
  - name: locksmithd.service
    enabled: false
    mask: true
  - name: fleet.service
    enabled: false
    mask: true
  - name: fleet.socket
    enabled: false
    mask: true
  - name: flanneld.service
    enabled: false
    mask: true
  - name: systemd-networkd-wait-online.service
    enabled: false
    mask: true
  - name: k8s-label-node.service
    enabled: true
    contents: |
      [Unit]
      Description=Adds labels to the node after kubelet startup
      After=k8s-kubelet.service
      Wants=k8s-kubelet.service
      [Service]
      Type=oneshot
      Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin"
      Environment="KUBECONFIG=/etc/kubernetes/kubeconfig/kubelet.yaml"
      ExecStart=/bin/sh -c '\
        while [ "$(kubectl get nodes $(hostname | tr '[:upper:]' '[:lower:]')| wc -l)" -lt "1" ]; do echo "Waiting for healthy k8s" && sleep 20s;done; \
        kubectl label nodes --overwrite $(hostname | tr '[:upper:]' '[:lower:]') node-role.kubernetes.io/master=""; \
        kubectl label nodes --overwrite $(hostname | tr '[:upper:]' '[:lower:]') kubernetes.io/role=master; \
        kubectl label nodes --overwrite $(hostname | tr '[:upper:]' '[:lower:]') node-role.kubernetes.io/control-plane=""; \
        for l in $(echo "{{.Cluster.Kubernetes.Kubelet.Labels}}" | tr "," " "); do \
            kubectl label nodes --overwrite $(hostname | tr "[:upper:]" "[:lower:]") $l; \
        done'
      [Install]
      WantedBy=multi-user.target
  - name: k8s-label-node.timer
    enabled: true
    contents: |
      [Unit]
      Description=Execute k8s-label-node every hour
      [Timer]
      OnCalendar=hourly
      [Install]
      WantedBy=multi-user.target
  - name: k8s-addons.service
    enabled: true
    contents: |
      [Unit]
      Description=Kubernetes Addons
      Wants=k8s-kubelet.service k8s-setup-network-env.service
      After=k8s-kubelet.service k8s-setup-network-env.service
      [Service]
      Type=oneshot
      Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/opt/bin"
      ExecStart=/opt/k8s-addons
      # https://github.com/kubernetes/kubernetes/issues/71078
      ExecStartPost=/usr/bin/systemctl restart k8s-kubelet.service
      [Install]
      WantedBy=multi-user.target
  - name: ensure-kube-proxy-vpa.service
    enabled: true
    contents: |
      [Unit]
      Description=Enable VPA for kube-proxy
      After=k8s-addons.service
      [Service]
      Type=oneshot
      Restart=on-failure
      ExecStartPre=/bin/bash -c "while ! /opt/bin/kubectl --kubeconfig=/etc/kubernetes/kubeconfig/addons.yaml get crd verticalpodautoscalers.autoscaling.k8s.io; do echo 'Waiting for VPA CRD to exists' && sleep 10; done"
      ExecStart=/opt/bin/kubectl --kubeconfig=/etc/kubernetes/kubeconfig/addons.yaml apply -f /srv/kube-proxy-vpa.yaml
      [Install]
      WantedBy=multi-user.target

{{ if .Debug.Enabled }}
  - name: logentries.service
    enabled: true
    contents: |
      [Unit]
      Description=Logentries
      After=systemd-networkd.service
      Wants=systemd-networkd.service
      StartLimitBurst=10
      StartLimitIntervalSec=600

      [Service]
      Restart=on-failure
      RestartSec=5
      Environment=LOGENTRIES_PREFIX={{ .Debug.LogsPrefix }}-master
      Environment=LOGENTRIES_TOKEN={{ .Debug.LogsToken }}
      ExecStart=/bin/sh -c 'journalctl -o short -f | sed \"s/^/${LOGENTRIES_TOKEN} ${LOGENTRIES_PREFIX} \\0/g\" | ncat data.logentries.com 10000'
      [Install]
      WantedBy=multi-user.target
{{ end }}

storage:
  files:
    - path: /boot/coreos/first_boot
      filesystem: root

    - path: /etc/ssh/trusted-user-ca-keys.pem
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;base64,{{ index .Files "conf/trusted-user-ca-keys.pem" }}"

    - path: /srv/calico-crd-installer.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-crd-installer.yaml" }}"

    - path: /srv/calico-crd-installer-rbac.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-crd-installer-rbac.yaml" }}"

    {{- if .CalicoPolicyOnly }}
    - path: /srv/calico-policy-only.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-policy-only.yaml" }}"
    {{- else }}
    - path: /srv/calico-datastore-migrator-pre.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-datastore-migrator-pre.yaml" }}"

    - path: /srv/calico-datastore-migrator-post.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-datastore-migrator-post.yaml" }}"

    - path: /srv/calico-all.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/calico-all.yaml" }}"
    {{- end }}

    {{- if not .DisableIngressControllerService }}
    - path: /srv/ingress-controller-svc.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/ingress-controller-svc.yaml" }}"
    {{- end }}

    - path: /etc/kubernetes/config/proxy-config.yml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "config/kube-proxy.yaml" }}"

    - path: /srv/kube-proxy-config.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "config/kube-proxy.yaml" }}"

    - path: /srv/kube-proxy-sa.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/kube-proxy-sa.yaml" }}"

    - path: /srv/kube-proxy-ds.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/kube-proxy-ds.yaml" }}"

    - path: /srv/rbac_bindings.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/rbac_bindings.yaml" }}"

    - path: /srv/rbac_roles.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/rbac_roles.yaml" }}"

    - path: /srv/priority_classes.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/priority_classes.yaml" }}"

    - path: /srv/psp_policies.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/psp_policies.yaml" }}"

    - path: /srv/psp_roles.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/psp_roles.yaml" }}"

    - path: /srv/psp_binding.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/psp_bindings.yaml" }}"

    - path: /srv/network_policies.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/network_policies.yaml" }}"

    - path: /opt/wait-for-domains
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/wait-for-domains" }}"

    - path: /opt/k8s-addons
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/k8s-addons" }}"

    - path: /opt/k8s-extract
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/k8s-extract" }}"

    - path: /opt/bin/setup-kubelet-environment
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/setup-kubelet-environment" }}"

    - path: /opt/bin/setup-apiserver-environment
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/setup-apiserver-environment" }}"

    - path: /etc/kubernetes/kubeconfig/addons.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/addons.yaml" }}"

    - path: /etc/kubernetes/config/proxy-kubeconfig.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/kube-proxy-master.yaml" }}"

    - path: /etc/kubernetes/kubeconfig/kube-proxy.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/kube-proxy-master.yaml" }}"

    - path: /etc/kubernetes/config/kubelet.yaml.tmpl
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "config/kubelet-master.yaml.tmpl" }}"

    - path: /etc/kubernetes/kubeconfig/kubelet.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/kubelet-master.yaml" }}"

    - path: /etc/kubernetes/kubeconfig/controller-manager.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/controller-manager.yaml" }}"

    - path: /etc/kubernetes/config/scheduler.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "config/scheduler.yaml" }}"

    - path: /etc/kubernetes/kubeconfig/scheduler.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "kubeconfig/scheduler.yaml" }}"

    {{ if not .DisableEncryptionAtREST -}}
    - path: /etc/kubernetes/encryption/k8s-encryption-config.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "k8s-resource/k8s-encryption-config.yaml" }}"

    {{ end -}}
    - path: /etc/kubernetes/policies/audit-policy.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "policies/audit-policy.yaml" }}"
    - path: /etc/kubernetes/manifests/k8s-api-healthz.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "manifests/k8s-api-healthz.yaml" }}"

    - path: /etc/kubernetes/manifest-templates/k8s-api-server.yaml.tmpl
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "manifests/k8s-api-server.yaml.tmpl" }}"

    - path: /etc/kubernetes/manifests/k8s-controller-manager.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "manifests/k8s-controller-manager.yaml" }}"

    - path: /etc/kubernetes/manifests/k8s-scheduler.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "manifests/k8s-scheduler.yaml" }}"

    - path: /etc/ssh/sshd_config
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/sshd_config" }}"

    - path: /etc/sysctl.d/hardening.conf
      filesystem: root
      mode: 0600
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/hardening.conf" }}"

    - path: /etc/audit/rules.d/10-docker.rules
      filesystem: root
      mode: 0600
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/10-docker.rules" }}"

    - path: /etc/docker/daemon.json
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/docker-daemon.json" }}"

    - path: /root/.docker/config.json
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/kubelet-docker-config.json" }}"

    - path: /etc/modules-load.d/ip_vs.conf
      filesystem: root
      mode: 0600
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/ip_vs.conf" }}"

    - path : /etc/containerd/config.toml
      filesystem: root
      mode: 420
      user:
        id: 0
      group:
        id: 0
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{ index .Files "conf/containerd-config.toml" }}"

    - path : /etc/systemd/system/containerd.service.d/10-use-custom-config.conf
      filesystem: root
      mode: 420
      user:
        id: 0
      group:
        id: 0
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{ index .Files "conf/10-use-custom-config.conf" }}"

    - path: /opt/install-debug-tools
      filesystem: root
      mode: 0544
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/install-debug-tools" }}"

    - path: /etc/calico/calicoctl.cfg
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/calicoctl.cfg" }}"

    - path: /etc/crictl.yaml
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/crictl" }}"

    - path: /etc/profile.d/setup-etcdctl.sh
      filesystem: root
      mode: 0444
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/etcd-alias" }}"

    - path: /srv/kube-proxy-vpa.yaml
      filesystem: root
      mode: 0444
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/kube-proxy-vpa" }}"

    - path : /etc/audit/rules.d/99-default.rules
      overwrite: true
      filesystem: root
      mode: 420
      user:
        id: 0
      group:
        id: 0
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{ index .Files "conf/99-default.rules" }}"

    {{ range .Extension.Files -}}
    - path: {{ .Metadata.Path }}
      filesystem: root
      user:
      {{- if .Metadata.Owner.User.ID }}
        id: {{ .Metadata.Owner.User.ID }}
      {{- else }}
        name: {{ .Metadata.Owner.User.Name }}
      {{- end }}
      group:
      {{- if .Metadata.Owner.Group.ID }}
        id: {{ .Metadata.Owner.Group.ID }}
      {{- else }}
        name: {{ .Metadata.Owner.Group.Name }}
      {{- end }}
      mode: {{printf "%#o" .Metadata.Permissions}}
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{ .Content }}"
        {{ if .Metadata.Compression }}
        compression: gzip
        {{end}}
    {{ end -}}

{{ range .Extension.VerbatimSections }}
{{ .Content }}
{{ end }}
`
