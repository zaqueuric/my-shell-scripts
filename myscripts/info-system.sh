#!/bin/bash

echo 
# Usuario
echo "Usuario: "
whoami
echo

# Data do sistema
echo "Data do sistema: "
date
echo

# Uso em disco
echo "Uso em disco: "
df -h
dh -i
echo

# Ponto de montagem
echo "Ponto de montagem dos dicos: "
lsblk
echo

# Informacoes do dispositivo
echo "Info hardware-software"
neofetch
echo
