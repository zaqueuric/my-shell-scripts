#!/bin/bash                                                                                        
                                                                                                   
# Autor: Zaqueu Ricardo                                                                            
# Date: 22-November-2020                                                                           
# Description: LINUXweek: Descomplicando o shell no linux - Guia prático                           
                                                                                                   
# /usr/local/sbin : Put this file in this directory.                                               
                                                                                                   
# Check the user permission                                                                        
if [ "$UID" != "0" ]; then                                                                         
        echo "You need the root user for execute this system. Out..."                              
        exit 1                                                                                     
fi                                                                                                 
                                                                                                   
# (!) If not exist                                                                                 
if [ ! -d "/tmp/backup" ]; then                                                                    
        echo "Creating the directory backup that not exist."                                       
        mkdir /tmp/backup                                                                          
fi                                                                                                 
                                                                                                   
tar -czpf /tmp/backup/etc.tar.gz /etc # Compact the directory etc in the tmp/backup folder         
                                                                                                   
                                       # For discompact the folder use: tar -zxpvf tar-file-name   
                                                                                                   
                                       # z: ZIP x:DESCOMPACT p:PRESERVE PERMISSION                 
                                                                                                   
                                       # v: VERBOSE DETAIL OF THE FILES OF THE FOLDER DESCOMPACTED 
                                                                                                   
                                       # f: SPECIFY THE FOLDER THAT YOU NEED TO DESCOMPACT         
                                       # you can use tar -tzvf, t: TO TEST                         
                                       # YOU CAN REMOVE AN PARAMETER eg: v to not display the conte
                                                                                                   
# (-f) If exist                                                                                    
if [ -f "/tmp/backup/etc.tar.gz" ]; then                                                           
                                                                                                   
        tar -tzf /tmp/backup/etc.tar.gz >/dev/null                                                 
                                        # >/dev/null direciona (>) para um dispositivo do sistema  
                                                                                                   
                                        # em lugar nenhum (null)                                   
        if [ "$?" != "0" ]; then                                                                   
                echo "Faild verification of backup."                                               
                exit 1 # return code 1 that determine that exit fail                               
        fi                                                                                         
fi                                                                                                 
                                                                                                   
echo "End of the program, all is fine. Congratulations!!!"                                         
:                                                                                                  
                                      # (:) is equal (exit 0) # return code 0 if is okay           
                                                                                                   
                                                                                                   
                                      # $?: Verify the output of the last comman                   
                                                                                                   
                                      # If 0 okay                                                  
                                      # If 1 Not okay                                              
                                                                                                   
                                                                                                   
                                      # Vim editor                                                 
                                      # :set nu -> For set the number of the lines                 
                                                                                                   
                                      # :q -> Out of the editor                                    
                                      # ESC and i: to insert the content in your file              
                                                                                                   
                                      # :wq -> Out and save the file                               
                                      # :ESC -> To insert other command                            
                                      # :Other command here!                                       