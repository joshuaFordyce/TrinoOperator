Trino Operator for Kubernetes

This project provides a Kubernetes Operator to automate the deployment and management of production-grade Trino clusters. It simplifies the operational complexities of running Trino on Kubernetes, allowing users to manage their clusters using a simple, declarative API.

Features

- Declarative Deployment: Deploy a complete Trino cluster (coordinator and workers) with a single Custom Resource.

- Automated Lifecycle Management: The operator handles scaling, configuration updates, and rolling upgrades of the Trino cluster.

- Production-Grade Configuration: Deploys with best-practice configurations for Trino components, including a Hive Metastore.

- Self-Service: Enables developers to deploy and manage their own Trino clusters with kubectl without needing deep operational knowledge.

Getting Started

* Prerequisites:

  * A Kubernetes cluster (version 1.25 or newer).

  * kubectl configured to connect to your cluster.

  * kustomize (v5.0.0 or newer).

  * operator-sdk (v1.31.0 or newer).

* Installation:

  * Clone the Repository:
  
  * git clone https://github.com/<your-username>/trino-operator.git
  * cd trino-operator


* Deploy the Operator:

  * make deploy


* Create an Operator Instance:

  * kubectl apply -f config/samples/trino_v1alpha1_trinocluster.yaml

  * Design Document & Architecture

Problem Statement

- Deploying and managing a production-ready Trino cluster on Kubernetes is a complex, multi-step process that involves orchestrating several components (coordinator, workers, Hive Metastore, etc.). Manually handling upgrades, scaling, and configuration changes is prone to errors and creates significant operational overhead.

Goals

- The goal of this project is to create a Kubernetes Operator that codifies this operational knowledge, providing a unified, declarative API that:

- Automates the complete lifecycle of a Trino cluster.

- Ensures consistent and reproducible deployments.

- Reduces manual intervention and human error.

System Architecture

- The Trino Operator follows the standard Kubernetes Operator pattern, extending the API to introduce a new custom resource. 

The system's core components are:

- Custom Resource Definition (CRD): 

  - The TrinoCluster CRD defines a new API object that serves as the single source of truth for the desired state of a Trino cluster.

- Trino Operator (The Controller): 

  - A Go application that runs inside the Kubernetes cluster. Its core is a reconciliation loop that constantly monitors for changes to TrinoCluster CRs.

- Kubernetes API Server: 

  - The central component where the TrinoCluster CRs are stored. The operator interacts with the API server to read the desired state and report the observed status.

- Managed Resources: 

  - The operator creates and manages standard Kubernetes resources like Deployments, StatefulSets, ConfigMaps, and Services to deploy and configure the Trino cluster components.

- Architecture Flow
  - A user creates a TrinoCluster Custom Resource (CR) manifest and applies it via kubectl apply -f ....

  - The Operator detects the new CR and its reconciliation loop begins.

  - The Operator reads the .spec (the desired state) from the CR.

  - It then creates a Trino coordinator, a service for it, and the specified number of worker nodes by creating standard Deployment and Service resources.

 - The Operator continually updates the .status of the TrinoCluster CR to reflect the real-time state of the deployed components (e.g., number of ready workers, cluster health).

- Contributing
  - This is an open-source project, and contributions are highly welcome. To get started:

- Fork this repository.

  - Clone your forked repository to your local machine.

  - Create a new branch for your feature or bug fix.

  - Submit a pull request to the main branch with a clear description of your changes.

License
 - This project is licensed under the Apache 2.0 License.
