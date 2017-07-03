FROM alpine:3.4

RUN apk add --no-cache ca-certificates

ADD ./kubeadm-token-operator /kubeadm-token-operator

ENTRYPOINT ["/kubeadm-token-operator"]
