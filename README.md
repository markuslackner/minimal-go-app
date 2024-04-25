# Cloud Native Software Delivery - UIBK/Dynatrace Tech Visit

[![O
pen in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/dt-uibk-workshop/minimal-go-app-demo)

<!-- TOC -->

* [Cloud Native Software Delivery - UIBK/Dynatrace Tech Visit](#cloud-native-software-delivery---uibkdynatrace-tech-visit)
* [Workshop Exercises](#workshop-exercises)
    * [1) Continuous Integration](#1-continuous-integration)
        * [1a) Introduction](#1a-introduction)
            * [Git Flow](#git-flow)
            * [Semantic Versioning](#semantic-versioning)
            * [Conventional Commits](#conventional-commits)
        * [1b) Create App Update on Feature-Branches](#1b-create-app-update-on-feature-branches)
        * [1c) Merge to Main to Trigger a Release](#1c-merge-to-main-to-trigger-a-release)
        * [1d) Access your minimal-go-app via exposed Ingress / Swagger UI](#1d-access-your-minimal-go-app-via-exposed-ingress--swagger-ui)
    * [2) Continuous Delivery](#2-continuous-delivery)
        * [2a) Introduction](#2a-introduction)
        * [2b) Create a new Release and watch CI / CD path to K8s-Cluster](#2b-create-a-new-release-and-watch-ci--cd-path-to-k8s-cluster)
        * [2c) Observe your Deployment on ArgoCD (Logs, Events, etc)](#2c-observe-your-deployment-on-argocd-logs-events-etc)
        * [2d) (advanced/optional) Scale up your Deployment to be Highly-Available](#2d-advancedoptional-scale-up-your-deployment-to-be-highly-available)
    * [3) Observability: Operations and Monitoring with Dynatrace](#3-observability-operations-and-monitoring-with-dynatrace)
        * [K8s Workshop Cluster Overview](#k8s-workshop-cluster-overview)
        * [Logs with Dynatrace Query Language](#logs-with-dynatrace-query-language)

<!-- TOC -->

# Workshop Exercises

![overview.png](docs%2Fimg%2Foverview.png)

## 1) Continuous Integration

![ci-focus.png](docs%2Fimg%2Fci-focus.png)

### 1a) Introduction

![ci.png](docs%2Fimg%2Fci.png)

#### Git Flow

![gitflow.png](docs%2Fimg%2Fgitflow.png)

https://docs.github.com/en/get-started/using-github/github-flow

#### Semantic Versioning

![semver.png](docs%2Fimg%2Fsemver.png)

https://semver.org

#### Conventional Commits

![conventional-commits.png](docs%2Fimg%2Fconventional-commits.png)

https://www.conventionalcommits.org

### 1b) Create App Update on Feature-Branches

![ci-feature-branch.png](docs%2Fimg%2Fci-feature-branch.png)

```bash
# checkout new feature-branch
git checkout -b feat/test-release

# change any file and commit with
git commit -m "feat: new release"

# push branch and create PR
git push
```

![ci-pr.png](docs%2Fimg%2Fci-pr.png)

* check build pipelines on PR and branch
* check release notes of dry-run

### 1c) Merge to Main to Trigger a Release

![ci-build.png](docs%2Fimg%2Fci-build.png)

* merge your PR
* follow build pipeline on main and check released artefacts (Helm-Chart and Container Image)
* find you new released GitHub Release

### 1d) Access your minimal-go-app via exposed Ingress / Swagger UI

![swagger-ui.png](docs%2Fimg%2Fswagger-ui.png)

https://minimal-go-app-demo.dt-uibk-workshop.com/

```bash
curl https://minimal-go-app-demo.dt-uibk-workshop.com/ # or open via browser
```

## 2) Continuous Delivery

![cd-focus.png](docs%2Fimg%2Fcd-focus.png)

### 2a) Introduction

![argocd.png](docs%2Fimg%2Fargocd.png)

* Open Argo CD https://argocd.dt-uibk-workshop.com/
* and navigate to your
  Application: [ArgoCD minimal-go-app-demo](https://argocd.dt-uibk-workshop.com/applications/argocd/minimal-go-app-demo.minimal-go-app-demo?view=tree&resource=)

### 2b) Create a new Release and watch CI / CD path to K8s-Cluster

![complete-ci-cd.png](docs%2Fimg%2Fcomplete-ci-cd.png)

* GitOps-State Repo: https://github.com/dt-uibk-workshop/gitops-state
* Argo CD https://argocd.dt-uibk-workshop.com/

### 2c) Observe your Deployment on ArgoCD (Logs, Events, etc)

![argocd-logs.png](docs%2Fimg%2Fargocd-logs.png)

* Watch Release Updates
* Delete Pod manually via ArgoCD UI
* Trigger OOM via Swagger UI
* Check Pod-Logs / Events via ArgoCD UI

### 2d) (advanced/optional) Scale up your Deployment to be Highly-Available

![replicas.png](docs%2Fimg%2Freplicas.png)

* create new release with updated replicas of 3, release again and wait for roll-out:

## 3) Observability: Operations and Monitoring with Dynatrace

### K8s Workshop Cluster Overview

![dt-workshop-cluster.png](docs%2Fimg%2Fdt-workshop-cluster.png)

### Logs with Dynatrace Query Language

![dt-logs.png](docs%2Fimg%2Fdt-logs.png)