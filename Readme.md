This repository presents a modified version of the Hyperledger Fabric lab, MiniFabric (SOURCE).

MiniFabric is a script based tool written with ansible playbooks for the automated management of hyperledger network nodes. Featuring a custom network configuration, modified paramteres and a number of tool integrations this version creates a local supplychain based fabric network. 

In addition to Caliper (SOURCE) and Explorer (SOURCE), this implementation features integrations with Fabric Connector (SOURCE), Splunk (SOURCE), Prometheous (SOURCE), CA Advisor (SOURCE) and Collector-for-Docker (SOURCE).

Run wth the below on Unix based OS:
./supplyflow up -o distillery.supply.com -c supplychain  -s couchdb -n supplyflow -r true

Diagrams for network configuration:


Screeshots of Tools:



