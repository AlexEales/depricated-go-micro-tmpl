#!/usr/bin/env bash

# TODO: In the future would be good to upgrade this to detect if the latest version is installed and prompt for upgrade if not

# Install go (if required and desired)
install_go () {
    local go_version=$(curl 'https://golang.org/VERSION?m=text')
    echo "Installing go ${go_version}"
    echo "Downloading binary"
    wget "https://dl.google.com/go/$go_version.linux-amd64.tar.gz"
    echo "Unzipping to /usr/local"
    sudo tar -C /usr/local -xzf $go_version.linux-amd64.tar.gz
    echo "Cleaning up"
    rm $go_version.linux-amd64.tar.gz
}

uninstall_go () {
    echo "Removing go installation"
    sudo rm -rf /usr/local/go
}

if ! command -v go &> /dev/null; then
    while true; do
        read -p "go not installed, do you want me to install it (required)? [y/n] " yn
        case $yn in
            [Yy]* ) install_go; break;;
            [Nn]* ) echo "Skipped installing go (not advised)"; break;;
            * ) echo "Please answer yes or no.";;
        esac
    done
else
    while true; do
        read -p "go already installed, do you want to re-install to get the latest version? [y/n] " yn
        case $yn in
            [Yy]* ) uninstall_go; install_go; break;;
            [Nn]* ) echo "Skipped re-installing go"; break;;
            * ) echo "Please answer yes or no.";;
        esac
    done
fi

# Gomock
if command -v go &> /dev/null; then
    if ! command -v mockgen &> /dev/null; then
        echo "Installing gomock"
        go install github.com/golang/mock/mockgen@v1.6.0
    else
        echo "gomock already installed, skipping"
    fi
else
    echo "go not installed, skipping gomock installation, once go is installed run: 'go install github.com/golang/mock/mockgen@v1.6.0'"
fi

# Buf
install_buf () {
    BIN="/usr/local/bin" && \
    VERSION="1.0.0-rc5" && \
    BINARY_NAME="buf" && \
        sudo curl -sSL \
            "https://github.com/bufbuild/buf/releases/download/v${VERSION}/${BINARY_NAME}-$(uname -s)-$(uname -m)" \
            -o "${BIN}/${BINARY_NAME}" && \
        sudo chmod +x "${BIN}/${BINARY_NAME}"
}

if ! command -v buf &> /dev/null; then
    while true; do
        read -p "buf not installed, do you want me to install it (required)? [y/n] " yn
        case $yn in
            [Yy]* ) install_buf; break;;
            [Nn]* ) echo "Skipped installing buf (not advised)"; break;;
            * ) echo "Please answer yes or no.";;
        esac
    done
else
    echo "Buf already installed, skipping"
fi

# Docker
if ! command -v docker &> /dev/null; then
    echo "Docker not installed, it is essential to have docker installed for your platform of choice: https://docs.docker.com/get-docker/"
    exit
fi
echo "Docker already installed, skipping"

# K8s
if ! command -v kubectl &> /dev/null; then
    echo "Kubernetes not installed, it is essential to have kubernetes installed for your platform of choice: https://kubernetes.io/docs/tasks/tools/"
    exit
fi
echo "Kubernetes already installed, skipping"

# Helm
install_helm () {
    echo "Installing helm"
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
    chmod 700 get_helm.sh
    ./get_helm.sh
    rm get_helm.sh
    echo "Helm installed"
}

if ! command -v helm &> /dev/null; then
    echo "Helm not installed, do you want me to install it (required)? [y/n]" yn
    case $yn in
            [Yy]* ) install_helm; break;;
            [Nn]* ) echo "Skipped installing helm (not advised)"; break;;
            * ) echo "Please answer yes or no.";;
        esac
else
    echo "Helm already installed, skipping"
fi

# Minikube
install_minikube () {
    echo "Installing minikube"
    curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
    rm minikube-linux-amd64
    echo "Minikube installed"
}

if ! command -v helm &> /dev/null; then
    echo "Minikube not installed, do you want me to install it (required)? [y/n]" yn
    case $yn in
            [Yy]* ) install_minikube; break;;
            [Nn]* ) echo "Skipped installing minikube, please install manually or some other local k8s cluser"; break;;
            * ) echo "Please answer yes or no.";;
        esac
else
    echo "Minikube already installed, skipping"
fi
