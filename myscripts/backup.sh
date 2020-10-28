#!/bin/bash

# *******************************************************************************
# Name: Zaqueu Ricardo **********************************************************
# Ocupation: Engenheiro DevOps **************************************************
# Company: Grupo Luminae Energia ************************************************
# Date: 28-10-2020 **************************************************************
# *******************************************************************************

# Variables
origin_bkp_files="/data/db/" # Diretorios ou arquivo que se deseja fazer bkp
destination_bkp="/backup/testebkp/" # Diretorio de destino dos arquivos de backup

dia=$(date +%d-%m-%y) # Data do sistema
hostname=$(hostname -s) # Nome da maquina

file="$hostname-$dia.tar.gz"

# Functions

sleep 1 

echo "Realizando backup: $origin_bkp_files para $destination_bkp/$file"

sleep 1

# Compact the files
echo "Compactando os arquivos!"
tar -zcvf $destination_bkp/$file $origin_bkp_files
echo ""

sleep 1

echo "Sucesso!"
echo ""

echo "Listando os arquivos do destino"
ls -lh $destination_bk # -lh para mostrar os arquivos em MB
echo ""
