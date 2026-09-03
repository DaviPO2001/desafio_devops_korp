#!/bin/bash

echo "Atualizando lista de pacotes..."
apt update

echo "Instalando Ansible..."
apt install -y ansible

#echo "Iniciando o playbook de automação..."

#cd /vagrant/

#ansible-playbook -i ansible/inventory.ini ansible/playbook.yml




