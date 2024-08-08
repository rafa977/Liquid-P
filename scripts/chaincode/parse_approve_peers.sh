#!/bin/bash

# Declare a string array

peers=()

# nOrgs=$(./getAllOrgs.sh)
# echo $nOrgs;
size=$(cat ./network/docker/orgsForApprove.yaml | shyaml get-value Size) > /dev/null

for (( c=0; c<$size; c++ ))
do
   name=$(cat ./network/docker/orgsForApprove.yaml | shyaml get-value Organizations.$c.Name) > /dev/null
   mspId=$(cat ./network/docker/orgsForApprove.yaml | shyaml get-value Organizations.$c.ID) > /dev/null
   port=$(cat ./network/docker/orgsForApprove.yaml | shyaml get-value Organizations.$c.Port) > /dev/null
   peer=$(cat ./network/docker/orgsForApprove.yaml | shyaml get-value Organizations.$c.Peer) > /dev/null

   status=$?
    if [ $status -ne 1 ]; then
        peers+=($mspId','$name','$port','$peer)
    fi


done

# for value in "${peers[@]}"
# do
#      echo "Peers: " $value
#      IFS=',' read -r -a peerInfo <<< "$value"
#      msp=$(echo ${peerInfo[0]} | tr '[:upper:]' '[:lower:]')
#      echo $msp
# done

echo ${peers[@]}

