All users can be created running the script : 'createAll.sh' of the scripts folder
Otherwise please execute the following commands one by one

#Peers for Liquid
./crypto-gen.sh -t peer -o liquid -u 'peer0' -p 6054 -c ca_liquid
./crypto-gen.sh -t peer -o liquid -u 'peer1' -p 6054 -c ca_liquid

#Peers for Applicant
./crypto-gen.sh -t peer -o applicant -u 'peer0' -p 7054 -c ca_applicant
./crypto-gen.sh -t peer -o applicant -u 'peer1' -p 7054 -c ca_applicant

#Peers for Financer
./crypto-gen.sh -t peer -o financer -u 'peer0' -p 8054 -c ca_financer
./crypto-gen.sh -t peer -o financer -u 'peer1' -p 8054 -c ca_financer

#Peers for AltFinancer
./crypto-gen.sh -t peer -o altfinancer -u 'peer0' -p 8154 -c ca_altfinancer
./crypto-gen.sh -t peer -o altfinancer -u 'peer1' -p 8154 -c ca_altfinancer

#Peers for Auditor
./crypto-gen.sh -t peer -o auditor -u 'peer0' -p 9054 -c ca_auditor
./crypto-gen.sh -t peer -o auditor -u 'peer1' -p 9054 -c ca_auditor


#Orderer for Authorizer 
./crypto-gen.sh -o authorizer -t orderer -p 9158 -c ca_authorizer


#Administrator for Liquid
./crypto-gen.sh -t admin -o liquid -u 'Admin' -p 6054 -c ca_liquid
./addAttribute.sh -o liquid -c ca_liquid -u Admin -a 'role=superadmin:ecert' -i 6054 -p Admin_liquid

#Administrator for Applicant
./crypto-gen.sh -t admin -o applicant -u 'Admin' -p 7054 -c ca_applicant
./addAttribute.sh -o applicant -c ca_applicant -u Admin -a 'role=superadmin:ecert' -i 7054 -p Admin_applicant

#Administrator for Financer
./crypto-gen.sh -t admin -o financer -u 'Admin' -p 8054 -c ca_financer
./addAttribute.sh -o financer -c ca_financer -u Admin -a 'role=superadmin:ecert' -i 8054 -p Admin_financer

#Administrator for AltFinancer
./crypto-gen.sh -t admin -o altfinancer -u 'Admin' -p 8154 -c ca_altfinancer
./addAttribute.sh -o altfinancer -c ca_altfinancer -u Admin -a 'role=superadmin:ecert' -i 8154 -p Admin_altfinancer

#Administrator for Auditor
./crypto-gen.sh -t admin -o auditor -u 'Admin' -p 9054 -c ca_auditor
./addAttribute.sh -o auditor -c ca_auditor -u Admin -a 'role=superadmin:ecert' -i 9054 -p Admin_auditor





#Users for Liquid
./crypto-gen.sh -t user -o liquid -u 'George' -p 6054 -c ca_liquid
./crypto-gen.sh -t user -o liquid -u 'Max' -p 6054 -c ca_liquid
./crypto-gen.sh -t user -o liquid -u 'Jason' -p 6054 -c ca_liquid
./crypto-gen.sh -t user -o liquid -u 'Rafael' -p 6054 -c ca_liquid

#Users for Applicant
./crypto-gen.sh -t user -o applicant -u 'Smith' -p 7054 -c ca_applicant
./crypto-gen.sh -t user -o applicant -u 'Jones' -p 7054 -c ca_applicant
./crypto-gen.sh -t user -o applicant -u 'Linda' -p 7054 -c ca_applicant
./crypto-gen.sh -t user -o applicant -u 'Mary' -p 7054 -c ca_applicant

#Users for Financer
./crypto-gen.sh -t user -o financer -u 'Davis' -p 8054 -c ca_financer
./crypto-gen.sh -t user -o financer -u 'Amy' -p 8054 -c ca_financer
./crypto-gen.sh -t user -o financer -u 'Anna' -p 8054 -c ca_financer
./crypto-gen.sh -t user -o financer -u 'Jerry' -p 8054 -c ca_financer

#Users for AltFinancer
./crypto-gen.sh -t user -o altfinancer -u 'Davis' -p 8154 -c ca_altfinancer
./crypto-gen.sh -t user -o altfinancer -u 'Amy' -p 8154 -c ca_altfinancer
./crypto-gen.sh -t user -o altfinancer -u 'Anna' -p 8154 -c ca_altfinancer
./crypto-gen.sh -t user -o altfinancer -u 'Jerry' -p 8154 -c ca_altfinancer

#Users for Auditor
./crypto-gen.sh -t user -o auditor -u 'Davis' -p 9054 -c ca_auditor
./crypto-gen.sh -t user -o auditor -u 'Amy' -p 9054 -c ca_auditor
./crypto-gen.sh -t user -o auditor -u 'Anna' -p 9054 -c ca_auditor
./crypto-gen.sh -t user -o auditor -u 'Jerry' -p 9054 -c ca_auditor