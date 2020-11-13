#!/bin/bash

# *******************************************************************************
# Name: Zaqueu Ricardo **********************************************************
# Ocupation: Engenheiro DevOps **************************************************
# Company: Grupo Luminae Energia ************************************************
# Date: 28-10-2020 **************************************************************
# *******************************************************************************

# Variables
origin_bkp_files="/data/db/" # Folders or files that you need to make backoup
destination_bkp="/backup/testebkp/" # Folder to destiny your content

dia=$(date +%d-%m-%y) # System data
hostname=$(hostname -s) # Machine name

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
ls -lh $destination_bk # Show the archive files capacity in MB
echo ""
