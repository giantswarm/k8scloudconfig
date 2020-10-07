package template

const BastionTemplate = `---
ignition:
  version: "2.2.0"
passwd:
  users:
    - name: giantswarm
      shell: "/bin/bash"
      uid: 1000
      groups:
        - "sudo"
        - "docker"
      sshAuthorizedKeys:
        - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCuJvxy3FKGrfJ4XB5exEdKXiqqteXEPFzPtex6dC0lHyigtO7l+NXXbs9Lga2+Ifs0Tza92MRhg/FJ+6za3oULFo7+gDyt86DIkZkMFdnSv9+YxYe+g4zqakSV+bLVf2KP6krUGJb7t4Nb+gGH62AiUx+58Onxn5rvYC0/AXOYhkAiH8PydXTDJDPhSA/qWSWEeCQistpZEDFnaVi0e7uq/k3hWJ+v9Gz0qqChHKWWOYp3W6aiIE3G6gLOXNEBdWRrjK6xmrSmo9Toqh1G7iIV0Y6o9w5gIHJxf6+8X70DCuVDx9OLHmjjMyGnd+1c3yTFMUdugtvmeiGWE0E7ZjNSNIqWlnvYJ0E1XPBiyQ7nhitOtVvPC4kpRP7nOFiCK9n8Lr3z3p4v3GO0FU3/qvLX+ECOrYK316gtwSJMd+HIouCbaJaFGvT34peaq1uluOP/JE+rFOnszZFpCYgTY2b4lWjf2krkI/a/3NDJPnRpjoE3RjmbepkZeIdOKTCTH1xYZ3O8dWKRX8X4xORvKJO+oV2UdoZlFa/WJTmq23z4pCVm0UWDYR5C2b9fHwxh/xrPT7CQ0E+E9wmeOvR4wppDMseGQCL+rSzy2AYiQ3D8iQxk0r6T+9MyiRCfuY73p63gB3m37jMQSLHvm77MkRnYcBy61Qxk+y+ls2D0xJfqxw== giantswarm"
    - name: vol
      shell: "/bin/bash"
      groups:
        - "sudo"
        - "docker"
      sshAuthorizedKeys:
        - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDFGDg/p4JWewXAs8kExJnCaNXEN1v2LZf0YWWiblHFp1+i2bp8qSmAJT3i6Yw0kHY2/6MotBCKAsFtlqxuhKaFs3jDcmdOugmWz4Qj7oerQ/ypJE/wZ9PY79gbK75aEKyOdVf7dUT6Ah+oSfETgpY/3a9pVZ/dSF3WBFIBw5k4YarFzcELQE4Bo4dcsLHsNrkI9Bk6gkGbTY+1TtfJmOu0bEXxXHdEq+JfW0MFssjh3I5n0DT09qDnztAvRAjjqjlyNKNt8reErV0LlvsDM5c+426Bz9JgM5vP3sD5ai8lpuH0iCBHoo9678XTKKTYbbz0s7kgXUb0vGS+GbOcaKBKmZ8a0xDpsft9+/LbmnuUic8b4c4/cRw5wSV1IYqyDqARp/d9PaJlYa22ISGnDbYmXUTsef0PhUenK9gtYrGsVhQmkqeLYiIYqwsl7+uouFMpQDmdZjY/B4fKcRA3oRGCFuwzT1vrtJL41dw9WyzM+3xnHTMFZdko9TlgDiEeu6gdpsTGJf4VALUWgXeyW/egte2im86kjMxzQuCw/aOmiYMqwZH2YfI0dS9jLuZbxePKTUounct66SrNXBrbu2d0BiPj6bl1dG6oZhwtArRnbiG5+cTakDvLhFgahTQFAT1De7o3Nr+BfjNQkVlQNKaIPUOdypiDNJE/6q/GOHVRQw== me@volodymyrstoiko.com"

systemd:
  units:
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

storage:
  files:
    - path: /boot/coreos/first_boot
      filesystem: root

    - path: /etc/ssh/trusted-user-ca-keys.pem
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;base64,{{ index .Files "conf/trusted-user-ca-keys.pem" }}"

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

    - path: /etc/docker/daemon.json
      filesystem: root
      mode: 0644
      contents:
        source: "data:text/plain;charset=utf-8;base64,{{  index .Files "conf/docker-daemon.json" }}"
`
