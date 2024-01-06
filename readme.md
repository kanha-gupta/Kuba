# Kuba CLI

Kuba is a command-line tool designed to facilitate interaction with Kubernetes clusters. It provides capabilities to obtain cluster information, manage resource quotas, and execute various other cluster-level operations.

## Prerequisites

Before installing Kuba, ensure you have the following prerequisites installed on your system:

- **Minikube**: A tool that allows you to run Kubernetes locally.
- **Golang**: The Go programming language environment, required for building and installing the application.

## Installation Instructions

Follow these steps to install Kuba CLI:

1. **Start Minikube**:

   Begin by starting Minikube. This will create a local Kubernetes cluster on your machine. Run the following command in your terminal:

   ```bash
   minikube start
   ```

2. **Clone the Kuba Repository**:

   Clone the Kuba repository from GitHub to your local machine using the following command:

   ```bash
   git clone https://github.com/kanha-gupta/Kuba
   ```

   Navigate to the cloned repository's directory:

   ```bash
   cd Kuba
   ```

3. **Build the Application**:

   Compile the application by running the build command in the root of the Kuba directory:

   ```bash
   go build
   ```

4. **Install the Application**:

   Install the Kuba CLI tool on your system by executing:

   ```bash
   go install
   ```

This completes the installation process. You are now ready to use Kuba CLI to interact with your Kubernetes clusters.

# Usage

**Kuba** provides the main command `kuba` for managing your Kubernetes minikube cluster. Here are some of the key functionalities it offers:

## Creating Kubernetes Resources

You can use the `create` subcommand to create Kubernetes resources from a YAML file.

```bash
kuba create --fp=<yaml_file_path> --ns=<namespace>
```

- `--fp`: Path to the YAML file containing the resource definition.
- `--ns`: namespace name

## Deleting Kubernetes Resources

To delete a Kubernetes resource, use the `delete` subcommand.

```bash
kuba delete --k=<resource_kind> --rn=<resource_name> --ns=<namespace>
```

- `--k`: Kind of the resource you want to delete (e.g., Deployment, Service, Pod, etc.).
- `--rn`: Name of the resource to be deleted.
- `--ns`: Name of the namespace

## Getting Resource Details
 
Kuba allows you to view specific details of a resource. For example, you can view details of a Deployment, Service, or Pod.

```bash
kuba details deployment -d=<deployment_name> --ns=<namespace>
kuba details service -s=<service_name> --ns=<namespace>
kuba details pod -p=<pod_name> --ns=<namespace>
```

## Namespace Details

You can obtain details about a specific namespace using the `namespace` subcommand.

```bash
kuba details namespace --ns=<namespace_name>
```


## Listing All Resources

You can list all resources or specific types of resources using the `show` subcommand.

```bash
kuba show all --ns=<namespace>
kuba show deploy --ns=<namespace>
kuba show services --ns=<namespace>
kuba show pods --ns=<namespace>
kuba show namespaces
```

- `--ns`: (Optional) Filter resources by namespace. If not provided, it will show resources from all namespaces.

Remember to use the `--ns=<namespace>` flag at the root command level to specify the namespace for subsequent commands. This flag will apply to all commands unless explicitly overridden in the subcommands.


