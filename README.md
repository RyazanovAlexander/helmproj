# Helmproj

The Helmproj is a preprocessor that allows to render Helm charts from multiple repositories before deploying them to a Kubernetes cluster.

A typical structure for a [Multi-repo](https://hackernoon.com/mono-repo-vs-multi-repo-vs-hybrid-whats-the-right-approach-dv1a3ugn) project on a development machine might have the following directory structure:
```
project (git clone https://github.com/project-name/project)
├── skaffold.yaml
├── project.yaml
└── Makefile (commands applied to all microservices)

frontend (git clone https://github.com/project-name/frontend)
├── Makefile
├── src
│    └── ...
└── helm
    ├── ...
    ├── Chart.yaml
    ├── values.yaml
    ├── values.dev.yaml
    └── values.prod.yaml

backend (git clone https://github.com/project-name/backend)
├── Makefile
├── src
│    └── ...
└── helm
    ├── ...
    ├── Chart.yaml
    ├── values.yaml
    ├── values.dev.yaml
    └── values.prod.yaml
```

The content of `project.yaml` might look like this:

```yaml
values:
  scalar: 1
  tree:
    scalar: tree-scalar
  array:
    - node1
    - node2
  map:
    - node1: val1
      prop11: pv11
    - node2: val2
      prop21: pv21
  frontend:
    prod: frontend-prod-value
  backend:
    prod: backend-prod-value
charts:
  - name: rendered-frontend
    path: ./examples/charts/frontend
    additionlValues:
      - values.prod.yaml
    appVersion: 1.0.1
  - name: rendered-backend
    path: ./examples/charts/backend
    additionlValues:
      - values.prod.yaml
    appVersion: 1.0.2
outputFolder: ./tmp/rendered
```

Common variables for rendering and the charts list are taken from the configuration file. The file name `project.yaml` is used by default. Rendering occurs only for files `values.yaml`, `values.[mode].yaml`. In addition, the program can change the `appVersion` in the `Chart.yaml` file. The rendered charts are copied to the `outputFolder` directory, specified in configuration file.

The values section contains variables that can be substituted in `values.yaml` as follows:

```yaml
scalar: value # {{ .Project.scalar }}

tree:
  scalar: value # {{ .Project.tree.scalar }}

array: # {{ .Project.array }}
  - value1
  - value2

map: # {{ .Project.map }}
  - key1: val1
  - key2: val2
```

Running the `make deploy` command, the Helm charts are rendered using Helmproj and then deployed to Minikube or any other Kubernetes cluster.
After rendering, `values.yaml` will have the following content:
```yaml
scalar: 1

tree:
  scalar: tree-scalar

array:
  - node1
  - node2

map:
  - node1: val1
    prop11: pv11
  - node2: val2
    prop21: pv21
```

This approach allows to substitute common variables in all charts of the project, leaving the Helm charts unchanged. Most often this is required when writing a CI/CD. It can also come in handy for local development via [Skaffold](https://skaffold.dev).

## Alternative solutions

It is possible to write an umbrella chart ([Subcharts and Global Values](https://helm.sh/docs/chart_template_guide/subcharts_and_globals)), but in this case the ability to upgrade as separate charts are lost. The [chart hooks](https://helm.sh/docs/topics/charts_hooks) mechanism has a number of limitations and is practically not applicable to complex microservice updates.

Another alternative could be [Helmfile](https://github.com/roboll/helmfile), but this tool does not only render, but also deploy. This is not always convenient, especially when you want to use another tool to deploy, such as [Skaffold](https://skaffold.dev) or [Draft](https://draft.sh). Helmfile can render to yaml file without deployment, but in this case you lose the ability to use helm.

## Installation

- download one of [releases](https://github.com/RyazanovAlexander/helmproj/releases)
- run as a [container](#running-as-a-container)

## Examples

An example is located in the folder [examples](https://github.com/RyazanovAlexander/helmproj/tree/main/examples). For starting the example with rendering helm chart run the command `make render`. The result of the program can be viewed in the `examples/tmp/rendered` folder.

To run the example with deploy, it’s necessary to install [Skaffold](https://skaffold.dev/docs/install), [Minikube](https://v1-18.docs.kubernetes.io/docs/tasks/tools/install-minikube) and execute the `make example` command in the root of the repository.

## Running as a container

TODO

## CLI Reference

```
This application is a tool for preprocessing values.yaml and Chart.yaml files.

Common actions for Helmproj:

- helmproj:                           preprocessing helm charts based on information from the ./project.yaml file
- helmproj -f ./path/to/project.yaml: preprocessing helm charts based on information from the project.yaml file located at the specified path

Environment variables:

| Name                               | Description                                                                       |
|------------------------------------|-----------------------------------------------------------------------------------|
| $HELMPROJ_DEBUG                    | indicate whether or not Helmprog is running in Debug mode                         |

Usage:
  helmproj [flags]
  helmproj [command]

Available Commands:
  help        Help about any command
  version     Show the version for Helmproj.

Flags:
      --debug         enable verbose output
  -d, --dry-run       test run. No files are generated. The output is carried out to the shell
  -f, --file string   path to project file (default "./project.yaml")
  -h, --help          help for helmproj

Use "helmproj [command] --help" for more information about a command.
```
