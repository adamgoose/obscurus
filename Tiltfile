# Not needed for docker for mac
default_registry("registry.local:5000")

# Vault
docker_build('vault-init', context='./hack/vault/init')
k8s_yaml(kustomize('./hack/vault'))

# Obscurus
docker_build('obscurus', context='.', live_update=[
    sync('./public', '/public')
])
k8s_yaml(kustomize('./hack/obscurus'))