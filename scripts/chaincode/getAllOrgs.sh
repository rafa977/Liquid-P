#!/bin/bash


mapfile -t my_array < <( find ./network/organizations/peerOrganizations -maxdepth 1 -type d  -printf '%f\n' )

numberOfOrgs=0
numberOfOrderer=0

for user in "${my_array[@]}"
do
    if [[ "$user" != "peerOrganizations" ]]
    then
        #echo "Peer Organization: " $user
        ((numberOfOrgs=numberOfOrgs+1))
    fi
done


mapfile -t my_array < <( find ./network/organizations/ordererOrganizations -maxdepth 1 -type d  -printf '%f\n' )

for user in "${my_array[@]}"
do
    if [[ "$user" != "ordererOrganizations" ]]
    then
        #echo "Orderer Organization: " $user
        ((numberOfOrderer=numberOfOrderer+1))
    fi
done

echo $((numberOfOrderer+numberOfOrgs))

