# Untitled

This repository presents a modified version of the Hyperledger Fabric lab, MiniFabric (SOURCE).

MiniFabric is a script based tool written with ansible playbooks for the automated management of hyperledger network nodes. Featuring a custom network configuration, modified paramteres and a number of tool integrations this version creates a local supplychain based fabric network. 

# Features

---

In addition to Caliper (SOURCE) and Explorer (SOURCE), this implementation features new integrations for the purpose of security testing.

- Fabric Connector (SOURCE)
- Splunk (SOURCE)
- Prometheous (SOURCE)
- CA Advisor (SOURCE)
- Collector-for-Docker (SOURCE).

# Perquisites

---

- Docker (18.4<)
- 10GB of disk space
- Unix based OS

# Getting Started

---

1. Clone the repository.
2. From within the repository, run the below:

```bash
./supplyflow up -o distillery.supply.com -c supplychain  -s couchdb -n supplyflow -r true
```

1. From within the monitoring directory, run the below:

```bash
docker-compose up -d
```

This is launch the monitoring tools where Splunk can be accesses at [https://localhost:8000](https://localhost:8000)

Diagrams for network configuration:

Screeshots of Tools: