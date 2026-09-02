#!/bin/bash

set -e

echo "=============================================="
echo "ATUALIZANDO LISTA DE PACOTES"
echo "=============================================="

apt update


echo "=============================================="
echo "INSTALANDO ANSIBLE"
echo "=============================================="

apt install -y ansible


echo "=============================================="
echo "EXECUTANDO PLAYBOOK"
echo "=============================================="

cd /vagrant

ansible-playbook \
  -i ansible/inventory.ini \
  ansible/playbook.yml