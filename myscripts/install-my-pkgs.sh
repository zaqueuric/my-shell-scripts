#!/bin/bash

# Author: Zaqueu Ricardo
# Date: 18-10-20

# Instalar meus principais pacotes (Programas)
# - pacotes necessarios para imagens docker...
# - pode ser usado para atualizar tambem os pacotes

echo

echo "Update the system"
apt-get update
echo

echo "Upgrate the system"
apt-get upgrade
echo

echo "Install systemd"
apt-get install systemd
echo 

echo "Install nano editor"
apt-get install nano
echo 

echo "Install vim editor"
apt-get install vim
echo

echo "Install man pages-info"
apt-get install man
echo 

echo "Install neofetch"
apt-get install neofetch
echo 

echo "Install htop"
apt-get install htop
echo

echo "Install git" 
apt-get install git
echo
