# alf-k8s

Alfresco Community Kubernetes Deployment CLI using [Helm](https://helm.sh) and [Docker Desktop](https://docs.docker.com/desktop/) or [Kubernetes In Docker](https://kind.sigs.k8s.io/) (KinD) cluster.

Additional details are available in [ACS Deployment](https://github.com/Alfresco/acs-deployment/blob/master/docs/helm/desktop-deployment.md).

Requires separate install of [kubectl](https://kubernetes.io/docs/reference/kubectl/) and [Helm](https://helm.sh).

## Enabling Kubernetes in Docker Desktop

Apply following configurations to Docker Desktop settings:

- `Settings > Resources > Advanced > Memory: 16 GB`
- `Settings > Kubernetes > Enable Kubernetes: ON`

After changing the necessary settings `Apply and restart` the docker desktop.

## Setting up kind

Take a look to the [KinD quickstart](https://kind.sigs.k8s.io/docs/user/quick-start/) to learn how to install the binary cli on your machine and to learn briefly the main commands that you can run.

Apply following configurations to Docker Desktop settings:

- `Settings > Resources > Advanced > Memory: 16 GB`
- `Settings > Kubernetes > Enable Kubernetes: OFF`

After changing the necessary settings `Apply and restart` the docker desktop.

## Usage

Download the binary compiled for your architecture (Linux, Windows or Mac OS) from [**Releases**](https://github.com/aborroy/alf-k8s/releases).

>> You may rename the binary to `alf-k8s`, all the following samples are using this command name by default.

Using `-h` flag provides detail on the use of the different commands available.

**Create**

`Create` command produces required assets to deploy Alfresco Community in Kubernetes.

```bash
$ ./alf-k8s create -h
Create assets to deploy Alfresco in Kubernetes

Usage:
  alf-k8s create [flags]

Flags:
  -h, --help                help for create
  -k, --kubernetes string   Kubernetes cluster: docker-desktop (default) or kind
  -o, --output string       Local Directory to write produced assets, 'output' by default
  -t, --tls string          Enable TLS protocol for ingress
  -v, --version string      Version of ACS to be deployed (23.1 or 23.2)
```

### Creating a sample deployment

Run the command selecting the Alfresco Community version to be deployed.

```bash
$ ./alf-k8s create -v 23.2
```

>> The previous command uses Docker Desktop as Kubernetes cluster, add `-k kind` to use [kind](https://kind.sigs.k8s.io) instead.

Kubernetes assets will be produced by default in `output` folder:

```bash
$ tree output
output
├── start.sh
├── stop.sh
└── values
    ├── community_values.yaml
    ├── resources_values.yaml
    └── version_values.yaml
```

Alfresco can be deployed to Kubernetes (only in Mac OS or Linux) using provided shell script:

```bash
$ cd output
$ ./start.sh
...
You can access all components of Alfresco Content Services using the same root address, but different paths as follows:

  Content: http://localhost/alfresco
  Share: http://localhost/share
  API-Explorer: http://localhost/api-explorer
```

Once the deployment has been tested, resources can be released using the following shell script:

```bash
$ ./stop.sh
release "acs" uninstalled
namespace "alfresco" deleted
release "ingress-nginx" uninstalled
namespace "ingress-nginx" deleted
```

## Troubleshooting

### Lens

The easiest way to troubleshoot issues on a Kubernetes deployment is to use the [Lens](https://k8slens.dev) desktop application, which is available for Mac, Windows and Linux. Follow the [getting started guide](https://docs.k8slens.dev/v4.0.3/getting-started) to configure your environment.

### Kubernetes Dashboard

Alternatively, the traditional Kubernetes dashboard can also be used. Presuming you have deployed the dashboard in the cluster you can use the following steps to explore your deployment:

1. Retrieve the service account token with the following command:

    ```bash
    kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')
    ```

2. Run the kubectl proxy:

    ```bash
    kubectl proxy &
    ```

3. Open a browser and navigate to: `http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login`

4. Select "Token", enter the token retrieved in step 1 and press the "Sign in" button

5. Select "alfresco" from the "Namespace" drop-down menu, click the "Pods" link and click on a pod name
