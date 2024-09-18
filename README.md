# Kubernetes Operator 新手开发一文入门

- [Kubernetes Operator 新手开发一文入门](#kubernetes-operator-新手开发一文入门)
  - [概述](#概述)
    - [Kubernetes Operator 简介](#kubernetes-operator-简介)
    - [为什么使用 Operator](#为什么使用-operator)
    - [Operator 与 Controller 的关系](#operator-与-controller-的关系)
  - [核心概念](#核心概念)
    - [自定义资源 (CR) 和自定义资源定义 (CRD)](#自定义资源-cr-和自定义资源定义-crd)
      - [什么是 CR](#什么是-cr)
      - [什么是 CRD](#什么是-crd)
      - [CRD 的结构与定义详解](#crd-的结构与定义详解)
    - [控制器 (Controller)](#控制器-controller)
      - [Controller 的工作原理](#controller-的工作原理)
      - [Controller 的生命周期管理](#controller-的生命周期管理)
      - [核心控制循环解释 (Control Loop)](#核心控制循环解释-control-loop)
    - [Informer](#informer)
      - [Informer 的原理](#informer-的原理)
      - [SharedInformer 与 Controller 的配合](#sharedinformer-与-controller-的配合)
      - [Informer 的缓存机制及事件处理流程](#informer-的缓存机制及事件处理流程)
      - [核心组件介绍：ClientSet、Indexer、Lister](#核心组件介绍clientsetindexerlister)
    - [RBAC (Role-Based Access Control)](#rbac-role-based-access-control)
      - [什么是 RBAC](#什么是-rbac)
      - [创建 RBAC 规则与权限控制](#创建-rbac-规则与权限控制)
      - [Operator 中的 RBAC 规则应用](#operator-中的-rbac-规则应用)
  - [Kubebuilder 快速入门](#kubebuilder-快速入门)
    - [什么是 Kubebuilder](#什么是-kubebuilder)
      - [Kubebuilder 的功能概述](#kubebuilder-的功能概述)
      - [与 Operator SDK 的区别](#与-operator-sdk-的区别)
    - [安装 Kubebuilder](#安装-kubebuilder)
    - [Kubebuilder 的项目结构](#kubebuilder-的项目结构)
      - [初始化 Kubebuilder 项目](#初始化-kubebuilder-项目)
      - [项目文件结构详解](#项目文件结构详解)
      - [API、Controller 和配置文件的作用](#apicontroller-和配置文件的作用)
  - [MySQL Operator](#mysql-operator)
    - [1. 使用 Kubebuilder 初始化项目](#1-使用-kubebuilder-初始化项目)
    - [2.CRD 的定义与生成](#2crd-的定义与生成)
    - [3. 详细的 Controller 部分开发文档](#3-详细的-controller-部分开发文档)
    - [4. 手动部署 Operator 到 Kubernetes 集群并查看效果](#4-手动部署-operator-到-kubernetes-集群并查看效果)
      - [**测试 Operator**](#测试-operator)
  - [Kubebuilder 原理详解](#kubebuilder-原理详解)
    - [Kubebuilder 的核心组件](#kubebuilder-的核心组件)
    - [Kubebuilder 的工作流程](#kubebuilder-的工作流程)
    - [与 Kubernetes API 的集成](#与-kubernetes-api-的集成)
    - [总结：开发者如何使用 Kubebuilder 开发 Operator](#总结开发者如何使用-kubebuilder-开发-operator)


## 概述
### Kubernetes Operator 简介
Kubernetes Operator 是一类 Kubernetes 控制器，它能够自动化管理复杂的应用程序和其生命周期，通常被用来管理有状态应用（如数据库、缓存等）。通过扩展 Kubernetes API，Operator 可以将日常操作流程（如安装、升级、扩展、备份等）转换为 Kubernetes 原生对象，从而实现自动化和声明式管理。

### 为什么使用 Operator
Operator 通过将 DevOps 团队日常管理应用的运维知识和流程编码化，使复杂的应用程序管理变得简单和自动化。在 Kubernetes 中，Operator 可以持续监控自定义资源，并自动进行相应操作，确保应用程序的状态与用户期望一致。

### Operator 与 Controller 的关系
Operator 实际上是一个高级 Controller，它不仅负责监控和管理 Kubernetes 中的自定义资源 (CR)，还可以执行特定的业务逻辑。Controller 是 Kubernetes 架构中管理资源状态的核心组件，Operator 是对 Controller 的封装和扩展，专门用于复杂应用的生命周期管理。

## 核心概念

### 自定义资源 (CR) 和自定义资源定义 (CRD)

#### 什么是 CR
自定义资源 (Custom Resource, CR) 是 Kubernetes 用户可以定义的扩展对象，用于描述某个具体的应用或资源的期望状态。每个 CR 对象的结构基于其相应的 CRD (Custom Resource Definition)，通过 CR，用户可以声明他们希望 Kubernetes 管理的特定应用或服务。

#### 什么是 CRD
CRD（自定义资源定义）是 Kubernetes 的一种扩展机制，允许用户向 Kubernetes API 添加新的对象类型。通过 CRD，用户可以定义新的资源种类（类似于内置的 `Pod`、`Service` 等），并指定这些资源的结构和行为。

#### CRD 的结构与定义详解
- **`apiVersion`**: 指定 API 组和版本，例如 `apps/v1`。
- **`kind`**: 定义资源的类型，比如 `MySQL`、`Redis` 等。
- **`metadata`**: 描述 CR 对象的元数据信息，如名称、命名空间等。
- **`spec`**: 用于描述资源的期望状态，包含资源的配置项，如副本数、存储大小、版本等。
- **`status`**: 系统生成，用于记录资源的当前状态，如运行中的副本数、最后备份时间等。

### 控制器 (Controller)

#### Controller 的工作原理
Controller 是 Kubernetes 中的核心组件之一，用于确保集群中的资源状态与用户的期望状态一致。Controller 通过监听资源对象的变化事件（如创建、更新、删除等），并根据这些事件采取行动来调整实际状态。例如，当用户期望创建一个 MySQL 实例时，Controller 监听到 MySQL CR 的创建事件，并根据该 CR 的 `spec` 定义自动创建相应的 Kubernetes 资源（如 `Deployment`、`Service`）。

#### Controller 的生命周期管理
Controller 通过一个无限循环的 "控制循环"（Control Loop）来工作，它会不断地获取资源的当前状态，并与期望状态进行对比。如果发现不一致，Controller 会采取相应的操作来修正状态。Controller 的生命周期管理包括以下几个阶段：
- **监听事件**：Controller 通过 Informer 监听自定义资源的增删改等事件。
- **执行 Reconcile**：每当资源的状态发生变化时，Controller 会调用 `Reconcile` 函数来同步状态。
- **更新状态**：Controller 操作 Kubernetes 资源（如创建 Pod、删除 Service 等），确保资源状态符合期望，并更新状态信息。

#### 核心控制循环解释 (Control Loop)
控制循环（Control Loop）是 Controller 实现资源管理的核心机制。它的工作原理是：
1. **获取资源的实际状态**：通过 Kubernetes API 监听或查询资源的当前状态。
2. **对比期望状态和实际状态**：根据 CR 中定义的 `spec` 与资源的当前状态 (`status`) 进行对比。
3. **采取行动**：如果发现状态不一致，Controller 会采取相应的操作（如创建、删除、更新资源），确保资源的实际状态与用户期望一致。
4. **重复此过程**：控制循环是持续运行的，确保资源状态始终与期望一致。

### Informer

#### Informer 的原理

在 Kubernetes 中，**Informer** 是一个核心组件，它负责监听 Kubernetes API 资源的变化事件（如增、删、改等），并将这些事件通知给相应的 Controller。Informer 是基于 **缓存（Cache）** 的机制，通过减少直接与 API Server 的交互来提升系统的性能和效率。

Informer 的主要工作流程包括：
1. **List**：启动时，Informer 会从 API Server 中获取资源的当前状态列表，并将其缓存。
2. **Watch**：Informer 监听资源的变化（创建、更新、删除等），并将变化事件发送给对应的 Controller。
3. **同步数据**：Informer 将变化的资源同步到本地缓存，避免每次都向 API Server 请求资源，减少了对 API Server 的压力。

这种机制保证了 Kubernetes 系统的高可用性和高性能。

#### SharedInformer 与 Controller 的配合

**SharedInformer** 是 Informer 的高级版本，允许多个 Controller 共享同一个资源的缓存数据。由于 Kubernetes 集群中的资源可能会被多个 Controller 监控，如果每个 Controller 都独立与 API Server 交互，这会增加系统的负载。通过 SharedInformer，不同的 Controller 可以共享同一个数据源，避免重复查询，提升效率。

**SharedInformer 的主要特点**：
- **共享缓存**：多个 Controller 可以通过一个 SharedInformer 来访问相同的缓存数据，避免每个 Controller 独立维护缓存。
- **事件广播**：SharedInformer 会将监听到的事件广播给所有监听该资源的 Controller，每个 Controller 可以根据业务逻辑处理相应事件。

SharedInformer 的典型工作流程如下：
1. 启动时，SharedInformer 获取资源的列表并缓存。
2. 监听资源的变化，并更新缓存。
3. 将变化事件通知给所有订阅该资源的 Controller。

#### Informer 的缓存机制及事件处理流程

Informer 依赖本地缓存来加速数据访问。每当 API Server 中的资源发生变化时，Informer 会将变化事件存储在本地缓存中，并通过事件处理机制通知 Controller。缓存的机制允许 Controller 可以快速访问已经监听的资源，而无需频繁与 API Server 通信。

**Informer 的事件处理流程**：
1. **List 阶段**：Informer 启动时，通过 `List` 操作获取当前所有资源的完整状态，并将这些资源存储到本地缓存中。
2. **Watch 阶段**：Informer 通过 `Watch` 机制持续监听资源的变化事件，如资源的 `Add`（增加）、`Update`（更新）和 `Delete`（删除）。
3. **本地缓存更新**：每当有资源变化时，Informer 会更新本地缓存，并根据不同的事件类型（添加、更新、删除）触发不同的事件处理函数。
4. **事件通知**：Informer 会将资源变化事件传递给 Controller，Controller 再根据业务逻辑对事件进行处理。

通过这种机制，Informer 可以有效地减少 API Server 的压力，并加速 Controller 的事件处理过程。

下图展示了 **Informer 的事件处理机制**：

```
+-------------------+           +-----------------------+
|   API Server       |           |     Controller        |
|                   |           |                       |
|  (Add/Update/Delete)-----------> (Reconcile function)  |
|   (List/Watch)     |           |   Handles resource    |
+-------------------+           +-----------------------+
          ^                              ^
          |                              |
          |    +------------------+      |
          +----+ SharedInformer    +------+
               | (Cache + Event)   |
               +------------------+
```

#### 核心组件介绍：ClientSet、Indexer、Lister

Informer 的工作依赖于以下几个关键组件：

1. **ClientSet**：
   - **ClientSet** 是与 Kubernetes API Server 交互的客户端工具，负责发出请求来获取资源的列表（`List`）并监听资源的变化（`Watch`）。Informer 依赖 ClientSet 来与 API Server 通信。
   - 每种资源类型都有一个对应的 ClientSet，例如 `PodsClient`、`ServicesClient` 等。

2. **Indexer**：
   - **Indexer** 是 Kubernetes 缓存中的一种数据结构，用来根据特定的键（如对象的名称或标签）索引和检索资源对象。它可以高效地从缓存中获取指定资源的信息。
   - Indexer 提供了一种高效的资源查找方式，特别是在需要从大规模资源列表中查找特定对象时。

3. **Lister**：
   - **Lister** 是一个用于从本地缓存中快速获取资源的工具，通常结合 Indexer 使用。与直接查询 API Server 不同，Lister 可以从缓存中快速读取资源，提升查询效率。
   - Lister 允许控制器以类似于直接调用 Kubernetes API 的方式访问缓存数据。

### RBAC (Role-Based Access Control)

#### 什么是 RBAC

RBAC（基于角色的访问控制）是 Kubernetes 中用于管理用户和服务对集群中资源的访问权限的机制。通过 RBAC，集群管理员可以定义哪些用户或服务账户有权访问哪些资源，以及能够执行哪些操作。

RBAC 的四个核心概念：
1. **Role**：定义一组对资源的访问权限，如对 `Pods` 的读取权限或对 `Deployments` 的修改权限。
2. **RoleBinding**：将 `Role` 分配给一个或多个用户或服务账户，授权他们执行 `Role` 中定义的操作。
3. **ClusterRole**：类似于 `Role`，但 `ClusterRole` 可以跨命名空间作用，通常用于管理全局资源或集群级别的访问权限。
4. **ClusterRoleBinding**：将 `ClusterRole` 绑定到用户或服务账户，使其在整个集群范围内具备相应的权限。

#### 创建 RBAC 规则与权限控制

当部署 Kubernetes Operator 时，通常需要为 Operator 定义一系列 RBAC 规则，确保 Operator 能够访问和管理特定资源。为了确保 Operator 拥有必要的权限，必须创建相关的 `Role` 和 `RoleBinding`，或者 `ClusterRole` 和 `ClusterRoleBinding`。

创建一个简单的 `ClusterRole`，例如允许 Operator 访问和管理 `Deployment` 资源：

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: operator-role
rules:
  - apiGroups: [""]
    resources: ["pods", "services"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

然后，创建 `ClusterRoleBinding`，将该 `ClusterRole` 绑定到 Operator 的 `ServiceAccount`：

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: operator-role-binding
subjects:
  - kind: ServiceAccount
    name: operator-sa
    namespace: default
roleRef:
  kind: ClusterRole
  name: operator-role
  apiGroup: rbac.authorization.k8s.io
```

#### Operator 中的 RBAC 规则应用

Operator 作为 Kubernetes 控制器的一部分，需要管理集群中的各种资源。因此，正确配置 RBAC 对 Operator 的安全性和功能至关重要。如果 Operator 需要管理多种资源（如 CRD、Deployment、Service 等），则需要为 Operator 创建适当的 `ClusterRole` 和 `ClusterRoleBinding`，确保它具备足够的权限去执行任务。

在 Operator 项目中，RBAC 通常通过注解的方式生成。例如，Kubebuilder 会自动根据控制器中的注解生成 RBAC 清单：

```go
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
```

这些注解会在 `make manifests` 或 `make install` 时生成对应的 RBAC 文件，确保 Operator 有足够的权限管理 Kubernetes 资源。

## Kubebuilder 快速入门

### 什么是 Kubebuilder

**Kubebuilder** 是一个用于快速构建 Kubernetes Operator 的开发框架，它简化了 Operator 的开发流程，并自动生成所需的代码和配置文件。Kubebuilder 基于 **controller-runtime**，为开发者提供了一套完整的工具链，帮助他们轻松构建、测试和部署 Kubernetes 控制器和自定义资源。

#### Kubebuilder 的功能概述

1. **项目结构初始化**：Kubebuilder 提供了项目初始化工具，能够自动生成符合最佳实践的项目结构。开发者可以专注于业务逻辑，而不需要手动设置复杂的项目配置。

2. **自动生成 CRD**：开发者可以通过简单的命令定义自定义资源 (CRD) 的 API 结构，Kubebuilder 会自动生成对应的 CRD 定义文件以及与 Kubernetes API 交互的 Go 代码。

3. **自动生成控制器**：通过 Kubebuilder，开发者可以快速生成控制器 (Controller) 的基础代码，控制器负责管理自定义资源的生命周期和状态变化。

4. **RBAC 配置管理**：Kubebuilder 支持通过代码注解自动生成 Kubernetes RBAC (基于角色的访问控制) 配置文件，确保 Operator 拥有正确的权限。

5. **测试支持**：Kubebuilder 提供了内置的测试框架，支持单元测试、集成测试和端到端测试，帮助开发者在本地和 CI 管道中验证 Operator 的行为。

6. **Kustomize 集成**：Kubebuilder 使用 Kustomize 进行配置管理，简化了 Kubernetes 清单的定制和部署。开发者可以轻松管理 CRD、RBAC 和 Operator 的部署配置。

#### 与 Operator SDK 的区别

虽然 **Kubebuilder** 和 **Operator SDK** 都是用于开发 Kubernetes Operator 的框架，但它们有一些区别：

1. **框架基础**：
   - **Kubebuilder** 是基于 `controller-runtime` 构建的，从底层开始就提供了对 Kubernetes 控制器的更细粒度控制。
   - **Operator SDK** 最初是基于 Kubebuilder 的，也采用了 `controller-runtime`，但它集成了更多高级工具，如 Ansible、Helm，用于简化非 Go 语言开发者的 Operator 开发。

2. **开发语言支持**：
   - **Kubebuilder** 主要支持 Go 语言开发，专注于 Go 原生的 Operator 开发体验。
   - **Operator SDK** 不仅支持 Go，还支持通过 Ansible 和 Helm 开发 Operator，适合不熟悉 Go 语言的开发者。

3. **项目结构和生成工具**：
   - **Kubebuilder** 注重生成符合最佳实践的 Go 代码，项目结构清晰且与 Kubernetes 社区的标准紧密一致。
   - **Operator SDK** 提供更广泛的工具和命令，允许开发者通过不同的方式生成 Operator，但在生成 Go 项目时与 Kubebuilder 的结构较为相似。

4. **社区和支持**：
   - **Kubebuilder** 是 Kubernetes 官方提供的 Operator 开发框架，紧密与 Kubernetes 社区保持一致，跟随 Kubernetes 的更新而更新。
   - **Operator SDK** 起源于 Red Hat，并在 Ansible 和 Helm Operator 生态系统中有着强大的支持，尤其适合 Red Hat 的 OpenShift 平台。

### 安装 Kubebuilder

github地址：https://github.com/kubernetes-sigs/kubebuilder

- 安装步骤

```go
wget -c https://github.com/kubernetes-sigs/kubebuilder/releases/download/v4.2.0/kubebuilder_darwin_amd64
mv kubebuilder_darwin_amd64 /usr/local/bin/kubebuilder
chmod +x /usr/local/bin/kubebuilder

kubebuilder version
Version: main.version{KubeBuilderVersion:"4.2.0", KubernetesVendor:"1.31.0", GitCommit:"c7cde5172dc8271267dbf2899e65ef6f9d30f91e", BuildDate:"2024-08-17T09:41:45Z", GoOs:"darwin", GoArch:"amd64"}
```

### Kubebuilder 的项目结构

#### 初始化 Kubebuilder 项目

```go
kubebuilder init --domain siming.com --repo github.com/ZhangSIming-blyq/mysql-operator
INFO Writing kustomize manifests for you to edit...
INFO Writing scaffold for you to edit...
INFO Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.19.0
...
...
INFO Update dependencies:
$ go mod tidy
...
...
Next: define a resource with:
$ kubebuilder create api
```

生成自定义资源 (CRD) 的 API 和控制器代码, 这里因为我下面要创建的是MySQL资源，所以这里的group是apps，version是v1，kind是MySQL。

```go
kubebuilder create api --group apps --version v1 --kind MySQL
```

#### 项目文件结构详解

```
├── Dockerfile                     # 用于构建 Operator 的容器镜像
├── Makefile                       # 定义了常用的构建、测试、部署命令，简化操作流程
├── PROJECT                        # Kubebuilder 项目元数据文件，记录项目配置和版本信息
├── README.md                      # 项目介绍文件，记录项目的目标、安装步骤、功能等信息
├── api
│   └── v1
│       ├── groupversion_info.go   # 定义了 API 版本信息和资源组信息
│       ├── mysql_types.go         # 自定义资源 (CRD) 的结构定义，包含 CRD 的字段和序列化逻辑
│       └── zz_generated.deepcopy.go  # 自动生成的代码，用于深度拷贝自定义资源对象
├── bin
│   ├── controller-gen             # Kubebuilder 用于生成控制器的工具，提供 CRD、RBAC 等生成功能的二进制文件
│   └── controller-gen-v0.16.1     # 特定版本的 controller-gen 工具
├── cmd
│   └── main.go                    # Operator 入口点，初始化 Manager 并启动控制器
├── config
│   ├── crd
│   │   ├── kustomization.yaml     # CRD 的 kustomize 配置，用于定制 CRD 的生成
│   │   └── kustomizeconfig.yaml   # 自定义资源定义 (CRD) 的额外配置文件
│   ├── default
│   │   ├── kustomization.yaml     # 默认的 kustomize 配置文件，用于管理 Operator 的部署
│   │   ├── manager_metrics_patch.yaml  # 配置 Manager 的指标导出补丁
│   │   └── metrics_service.yaml   # 用于暴露 Operator 监控指标的服务配置
│   ├── manager
│   │   ├── kustomization.yaml     # 用于部署 Manager 的 kustomize 配置
│   │   └── manager.yaml           # Manager 的 Kubernetes 部署清单
│   ├── network-policy
│   │   ├── allow-metrics-traffic.yaml  # 网络策略，允许访问指标服务
│   │   └── kustomization.yaml     # 网络策略的 kustomize 配置
│   ├── prometheus
│   │   ├── kustomization.yaml     # Prometheus 监控的 kustomize 配置
│   │   └── monitor.yaml           # Prometheus 对 Operator 进行监控的规则
│   ├── rbac
│   │   ├── kustomization.yaml     # 用于生成 RBAC 配置的 kustomize 配置
│   │   ├── leader_election_role.yaml   # Leader 选举的角色权限配置
│   │   ├── leader_election_role_binding.yaml  # Leader 选举的角色绑定
│   │   ├── metrics_auth_role.yaml # 监控指标授权的 RBAC 配置
│   │   ├── metrics_auth_role_binding.yaml # 监控指标授权的角色绑定
│   │   ├── metrics_reader_role.yaml  # 用于读取指标的角色
│   │   ├── mysql_editor_role.yaml    # MySQL 资源编辑者角色
│   │   ├── mysql_viewer_role.yaml    # MySQL 资源查看者角色
│   │   ├── role.yaml                 # 默认的 Operator 角色
│   │   ├── role_binding.yaml         # 角色绑定，将角色分配给 ServiceAccount
│   │   └── service_account.yaml      # 定义 Operator 的 ServiceAccount
│   └── samples
│       ├── apps_v1_mysql.yaml       # 示例自定义资源，定义 MySQL 资源的 YAML 文件
│       └── kustomization.yaml       # 样例资源的 kustomize 配置
├── go.mod                           # Go 模块文件，记录依赖关系
├── go.sum                           # Go 依赖的版本锁定文件
├── hack
│   └── boilerplate.go.txt           # 代码文件的版权声明模板
├── internal
│   └── controller
│       ├── mysql_controller.go      # MySQL 控制器的核心逻辑，处理 CR 的状态同步和管理
│       ├── mysql_controller_test.go # 控制器的单元测试
│       └── suite_test.go            # 控制器的测试套件配置
└── test
    ├── e2e
    │   ├── e2e_suite_test.go        # 端到端测试的套件配置
    │   └── e2e_test.go              # 端到端测试的逻辑
    └── utils
        └── utils.go                 # 测试过程中使用的工具函数
```

**文件和目录详细说明**：

- **`Dockerfile`**：用于构建 Operator 的容器镜像。部署到 Kubernetes 集群之前，Operator 会被打包为 Docker 镜像。

- **`Makefile`**：包含常见的构建命令，如生成 CRD、安装 CRD、编译控制器代码、运行测试、打包 Operator 等。

- **`api/v1/`**：该目录包含自定义资源的定义文件。
  - **`mysql_types.go`**：定义了 MySQL CRD 的 API 结构体，包括 spec 和 status 字段。
  - **`zz_generated.deepcopy.go`**：通过代码生成工具自动生成的代码，用于深度拷贝自定义资源对象。

- **`controllers/`**：这个目录存放控制器逻辑。
  - **`mysql_controller.go`**：实现了核心的控制器逻辑，监听 MySQL CR 的变化并执行相应的操作，如创建、删除、更新 MySQL 实例。
  
- **`config/`**：存放所有与 Kubernetes 相关的配置文件，如 CRD、RBAC、部署和监控配置。
  - **`crd/`**：定义 CRD 的生成与管理。
  - **`rbac/`**：定义 Operator 所需的 RBAC 权限，包括 ServiceAccount 和角色绑定。
  - **`samples/`**：提供了示例自定义资源文件，用于创建 MySQL 实例。

- **`cmd/`**：存放 `main.go` 文件，作为 Operator 的入口点，负责启动控制器并与 Kubernetes API 交互。

- **`test/`**：存放测试代码，分为端到端测试（e2e）和辅助工具文件。

#### API、Controller 和配置文件的作用

1. api/ 目录：该目录用于定义自定义资源 (CRD) 的 API，包括资源的结构和字段。开发者可以在这里定义 Go 结构体，Kubebuilder 会根据这些结构体自动生成 CRD。
2. controllers/ 目录：该目录用于编写控制器逻辑。控制器负责监听自定义资源的状态变化，并根据需要采取相应的行动。Kubebuilder 会生成控制器的基础代码，开发者只需填充业务逻辑即可。
3. config/ 目录：包含了与 Operator 相关的配置文件，包括：CRD 定义：在 config/crd 目录下生成的 CRD 清单文件。RBAC 配置：在 config/rbac 目录下定义了 Operator 所需的 RBAC 权限。 Manager 配置：在 config/manager 目录下定义了控制器管理器的部署配置。

## MySQL Operator

### 1. 使用 Kubebuilder 初始化项目

在已经初始化的项目中（你已经运行过 `kubebuilder init`），你可以定义新的 API 和控制器逻辑，而不直接应用到集群。首先，我们生成 API 和控制器的代码：

```bash
kubebuilder create api --group apps --version v1 --kind MySQL --resource --controller
```

- **`--group apps`**：定义 API 组为 `apps`，通常与 Kubernetes 中已有的 `apps` 组保持一致。
- **`--version v1`**：定义 API 版本为 `v1`，表示资源版本为 `v1`。
- **`--kind MySQL`**：定义新自定义资源的种类（Kind）为 `MySQL`。
- **`--resource`**：生成与自定义资源 (CRD) 相关的代码。
- **`--controller`**：生成控制器相关的代码，控制器将用于管理 MySQL 实例的生命周期。

**生成文件内容：**

- `api/v1/mysql_types.go`：此文件将包含自定义资源的定义，包括 `Spec` 和 `Status` 的结构。
- `controllers/mysql_controller.go`：此文件将包含初始的控制器代码，用于后续管理 MySQL 资源的生命周期。
- `config/crd/`：此目录将包含生成的 Kubernetes CRD 定义清单文件。

### 2.CRD 的定义与生成

**流程综述：**
1. **生成 API 和控制器文件**：使用 `kubebuilder create api` 命令生成自定义资源相关的文件，但不应用。
2. **编辑 `mysql_types.go`**：定义 MySQL 资源的 `Spec` 和 `Status` 字段。
3. **生成代码**：使用 `make generate` 生成深度拷贝函数等辅助代码。
4. **生成 CRD 清单文件**：使用 `make manifests` 生成 Kubernetes CRD YAML 文件，但不将其应用到集群。
5. **手动查看或修改生成的文件**：在 `config/crd/bases/` 中找到生成的 CRD 文件，进一步检查或修改。

我们将编写自定义资源 `MySQL` 的 `Spec` 和 `Status` 字段，并使用 Kubebuilder 的标注 (annotations) 自动生成深度拷贝函数、CRD 定义等。

修改 `mysql_types.go`: `mysql_types.go` 文件是定义自定义资源类型的主要文件。在这里，我们需要定义 MySQL 的 `Spec`（期望状态）和 `Status`（当前状态）。

**文件路径：`api/v1/mysql_types.go`**

```go
package v1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MySQLSpec 定义了 MySQL 资源的期望状态
type MySQLSpec struct {
    // MySQL 用户名
    Username string `json:"username"`

    // MySQL 用户密码
    Password string `json:"password"`

    // 数据库名称
    Database string `json:"database"`

    // MySQL 实例的副本数（用于扩展）
    Size int32 `json:"size"`

    // 定期备份的时间表（Cron 表达式）
    BackupSchedule string `json:"backupSchedule"`

    // 备份存储路径
    BackupPath string `json:"backupPath"`
}

// MySQLStatus 定义了 MySQL 资源的当前状态
type MySQLStatus struct {
    // 当前可用的 MySQL 副本数
    ReadyReplicas int32 `json:"readyReplicas"`

    // 最近备份的时间
    LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`

    // 资源状态条件（如是否可用、是否需要备份等）
    Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MySQL 是 MySQL 资源的 Schema
type MySQL struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MySQLSpec   `json:"spec,omitempty"`
    Status MySQLStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MySQLList 包含 MySQL 资源的列表
type MySQLList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []MySQL `json:"items"`
}

func init() {
    SchemeBuilder.Register(&MySQL{}, &MySQLList{})
}
```

**代码解释**

- **`MySQLSpec`**：定义期望的 MySQL 实例的配置信息（如用户名、密码、副本数、备份计划等）。
- **`MySQLStatus`**：定义 MySQL 实例的运行状态，包括副本数、最后备份时间和状态条件。
- **`+kubebuilder:object:root=true`**：告诉 Kubebuilder 这是自定义资源的根对象。
- **`+kubebuilder:subresource:status`**：让 Kubernetes 自动创建一个子资源来管理 `status` 字段。

在完成 `mysql_types.go` 的定义后，运行以下命令来生成相应的辅助代码：

```bash
# `make generate` 的作用：
# 生成深度拷贝函数 `zz_generated.deepcopy.go`，确保在 Kubernetes 控制器中可以安全地拷贝自定义资源对象。
# 确保自定义资源的代码与 Kubernetes API 兼容。
make generate
```

你将在 `api/v1/` 目录下看到生成的 `zz_generated.deepcopy.go` 文件。

**生成 Kubernetes CRD 清单文件**

为了将定义的 CRD 注册到 Kubernetes，我们需要生成相应的 CRD 定义 YAML 文件。此时，**生成代码但不应用**。你可以使用 `make manifests` 命令生成清单文件。

```bash
# **`make manifests` 的作用**：
# - 该命令会基于 `mysql_types.go` 中的定义生成相应的 Kubernetes CRD 定义文件。
# - 文件会被生成到 `config/crd/bases/` 目录下。
make manifests
```

**生成的文件**

`config/crd/bases/apps.example.com_mysqls.yaml`：这个 YAML 文件包含了 MySQL 自定义资源的定义，你可以查看文件，里面会包含以下内容：
- `spec` 和 `status` 字段的定义。
- API 组名、版本信息和其他元数据。

### 3. 详细的 Controller 部分开发文档

**流程综述**

1. **生成控制器文件**：使用 `kubebuilder create api --controller` 生成控制器文件。
2. **编写控制器逻辑**：在 `mysql_controller.go` 中编写 `Reconcile` 函数，处理自定义资源的状态变化。
3. **生成代码**：运行 `make generate` 生成辅助代码。
4. **手动查看和修改文件**：生成的文件位于 `config/` 目录下，可手动查看或修改生成的 Kubernetes 清单。

当你运行 `kubebuilder create api --controller` 时，已经生成了一个初始的控制器文件。接下来，我们将编写控制器的业务逻辑。

**生成的文件路径**：`controllers/mysql_controller.go`

**1. 编辑 `mysql_controller.go`**

`mysql_controller.go` 文件中包含了控制器的主要逻辑。我们将在这里编写 `Reconcile` 函数，它将根据 MySQL CR 的状态采取相应操作。

```go
package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/ZhangSIming-blyq/mysql-operator/api/v1"
)

// MySQLReconciler 是 MySQL 控制器的结构体
type MySQLReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile 是控制器的核心逻辑，用于处理 MySQL CR 的状态变化
func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	/*
		是的，目前的 Reconcile 逻辑主要涵盖以下两点：

		创建 MySQL 的 Deployment：如果 MySQL CR 对应的 Deployment 不存在，会根据 CR 中定义的 spec 创建一个新的 Deployment。
		控制副本数：如果 MySQL CR 的 spec.size 与现有 Deployment 的副本数不匹配，控制器会更新 Deployment，以确保实际副本数与期望的副本数一致。
		除此之外，Reconcile 函数还会更新 MySQL CR 的 status，将 Deployment 中实际的副本数同步到 MySQL CR 的 ReadyReplicas 字段。
	*/
	log := log.FromContext(ctx)

	// 获取 MySQL 实例: 在 Kubernetes 的控制器中，r.Get 函数只会返回一个具体的对象，而不会返回多个对象。这是因为 req.NamespacedName 包含了特定的 namespace 和 name，代表了一个唯一的资源实例。因此，不会出现 Get 返回多个对象的情况。
	var mysql appsv1alpha1.MySQL
	if err := r.Get(ctx, req.NamespacedName, &mysql); err != nil {
		if errors.IsNotFound(err) {
			// 如果没有找到 MySQL 实例，可能已经被删除，不采取行动
			log.Info("MySQL resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MySQL")
		return ctrl.Result{}, err
	}

	// 检查 MySQL Secret 是否存在，如果不存在则创建
	var secret corev1.Secret
	secretName := mysql.Name // 使用 MySQL CR 的名称作为 Secret 名称
	err := r.Get(ctx, types.NamespacedName{Name: secretName, Namespace: mysql.Namespace}, &secret)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Secret for MySQL", "Secret.Namespace", mysql.Namespace, "Secret.Name", secretName)

		// 从 MySQL CR 中获取 username 和 password
		username := mysql.Spec.Username
		password := mysql.Spec.Password

		// 创建 Secret
		secret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: mysql.Namespace,
			},
			StringData: map[string]string{
				"MYSQL_ROOT_PASSWORD": password, // 使用 CR 中的 password
				"MYSQL_USER":          username, // 使用 CR 中的 username
			},
			Type: corev1.SecretTypeOpaque,
		}

		if err := r.Create(ctx, &secret); err != nil {
			log.Error(err, "Failed to create new Secret", "Secret.Namespace", mysql.Namespace, "Secret.Name", secretName)
			return ctrl.Result{}, err
		}
		// 重新排队 Reconcile
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	// 检查是否存在 MySQL Deployment，如果不存在则创建
	var deployment appsv1.Deployment
	err = r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, &deployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment for MySQL", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
		// 根据CRD的定义创建 Deployment
		dep := r.mysqlDeployment(&mysql)
		if err := r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
			return ctrl.Result{}, err
		}
		// 重新排队 Reconcile
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// 确保副本数与期望一致
	size := mysql.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", mysql.Namespace, "Deployment.Name", mysql.Name)
			return ctrl.Result{}, err
		}
		// 创建完secret没有必要立刻重建
		return ctrl.Result{}, nil
	}

	// 更新 MySQL 状态
	mysql.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	if err := r.Status().Update(ctx, &mysql); err != nil {
		log.Error(err, "Failed to update MySQL status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// mysqlDeployment 返回定义的 MySQL Deployment
func (r *MySQLReconciler) mysqlDeployment(mysql *appsv1alpha1.MySQL) *appsv1.Deployment {
	labels := map[string]string{"app": "mysql", "mysql_cr": mysql.Name}
	replicas := mysql.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysql.Name,
			Namespace: mysql.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "mysql",
						Image: "mysql:5.7",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name: "MYSQL_ROOT_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										Key: "MYSQL_ROOT_PASSWORD",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: mysql.Name,
										},
									},
								},
							},
						},
					}},
				},
			},
		},
	}
	// Set the owner reference for garbage collection
	ctrl.SetControllerReference(mysql, dep, r.Scheme)
	return dep
}

// SetupWithManager 将控制器注册到 Manager 中
func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.MySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
```

**代码逻辑解释**

1. **获取 MySQL CR 对象**：
   - 首先，通过 `r.Get` 函数获取当前的 MySQL 实例。如果实例不存在，控制器会跳过当前操作。

2. **检查 MySQL Deployment(还有前置的secret) 是否存在**：
   - 如果 MySQL 实例的 `Deployment` 不存在，控制器将根据 MySQL CR 的 `Spec` 字段创建新的 `Deployment`，用于启动 MySQL 容器。

3. **同步副本数**：
   - 如果 MySQL 的实际副本数与 `spec` 中定义的不一致，控制器会更新 `Deployment`，确保副本数与期望保持一致。

4. **更新 MySQL 的状态**：
   - 将 `Deployment` 中的副本数写入 MySQL CR 的 `Status` 字段，以反映当前 MySQL 实例的运行状态。

5. **SetupWithManager**：
   - 控制器通过 `SetupWithManager` 注册到管理器 (Manager) 中，监听 `MySQL` 资源以及 `Deployment` 的变化。

**六、生成 Controller 代码**

在编写完控制器逻辑之后，执行以下命令生成并更新相应的代码：

```bash
make generate
```

### 4. 手动部署 Operator 到 Kubernetes 集群并查看效果

完成代码生成后，接下来需要将生成的 Operator 部署到 Kubernetes 集群并验证其功能。以下步骤将指导你如何构建 Operator 镜像、部署 Operator、并测试其运行情况。

**构建并推送 Operator 镜像**

首先，我们需要将 Operator 打包为容器镜像。执行以下命令构建 Docker 镜像：

```bash
make docker-build IMG=<your-operator-image>:<tag>
```

将 `<your-operator-image>` 替换为你镜像的名称（例如 `myregistry/mysql-operator`），`<tag>` 替换为版本号（如 `v1.0.0`）。例如：

```bash
make docker-build IMG=myregistry/mysql-operator:v1.0.0
```

接着，将构建好的 Docker 镜像推送到镜像仓库：

```bash
make docker-push IMG=myregistry/mysql-operator:v1.0.0
```

**更新 Kubernetes 部署文件**

在 `config/manager/manager.yaml` 文件中，找到 `image:` 字段，将其更新为刚刚推送的 Docker 镜像：
```yaml
spec:
  containers:
    - name: manager
      image: myregistry/mysql-operator:v1.0.0
      command:
      - /manager
```

**生成 Kubernetes 清单文件**

使用以下命令生成 Kubernetes 所需的 CRD、RBAC 规则和部署清单文件：

```bash
make manifests
```

生成的文件将位于 `config/crd/`、`config/rbac/` 和 `config/manager/` 目录中。

**部署 Operator到kubernetes集群**

```bash
# 上传镜像到kubernetes所在机器
docker image save docker.io/myregistry/mysql-operator:v1.0.0 > dockerimage                                                                                      

du -sh dockerimage 
 81M	dockerimage

scp dockerimage siming-dev:~
dockerimage

# 部署各种yaml
k apply -f crd/bases/apps.siming.com_mysqls.yaml
k apply -f manager/manager.yaml
k apply -f rbac/service_account.yaml
k apply -f rbac/role_binding.yaml
k apply -f rbac/role.yaml

# 查看状态: 如果rbac权限不够，比如代码逻辑需要，记得手动更新
k logs -f controller-manager-6c784ddb46-2b7lb                                                                                                         
2024-09-17T12:46:43Z	INFO	setup	starting manager
2024-09-17T12:46:43Z	INFO	starting server	{"name": "health probe", "addr": "[::]:8081"}
I0917 12:46:43.256807       1 leaderelection.go:254] attempting to acquire leader lease system/b3982890.siming.com...
I0917 12:47:00.943439       1 leaderelection.go:268] successfully acquired lease system/b3982890.siming.com
2024-09-17T12:47:00Z	DEBUG	events	controller-manager-6c784ddb46-2b7lb_8dccc7f2-caa7-4e2e-8268-8726aab5a9cf became leader	{"type": "Normal", "object": {"kind":"Lease","namespace":"system","name":"b3982890.siming.com","uid":"bd617fbe-a002-4545-8a23-ed4dd77d7f82","apiVersion":"coordination.k8s.io/v1","resourceVersion":"1470723"}, "reason": "LeaderElection"}
2024-09-17T12:47:00Z	INFO	Starting EventSource	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "source": "kind source: *v1.MySQL"}
2024-09-17T12:47:00Z	INFO	Starting EventSource	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "source": "kind source: *v1.Deployment"}
2024-09-17T12:47:00Z	INFO	Starting Controller	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL"}
2024-09-17T12:47:01Z	INFO	Starting workers	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "worker count": 1}
```

#### **测试 Operator**

接下来，通过应用 MySQL 自定义资源来测试 Operator 的功能。你可以修改 `config/samples/apps_v1_mysql.yaml` 文件来创建一个 MySQL 实例。例如：
```yaml
apiVersion: apps.example.com/v1
kind: MySQL
metadata:
  name: example-mysql
spec:
  username: root
  password: password123
  database: mydatabase
  size: 3
  backupSchedule: "*/5 * * * *"
  backupPath: /backups
```

创建mysql，查看operator日志

```bash
k logs -f controller-manager-578b69d9d4-t228b                                                                                                              ok | base py | with ubuntu@VM-24-14-ubuntu | at 22:45:08
2024-09-17T14:41:49Z	INFO	setup	starting manager
2024-09-17T14:41:49Z	INFO	starting server	{"name": "health probe", "addr": "[::]:8081"}
I0917 14:41:49.676064       1 leaderelection.go:254] attempting to acquire leader lease system/b3982890.siming.com...
I0917 14:42:17.231093       1 leaderelection.go:268] successfully acquired lease system/b3982890.siming.com
2024-09-17T14:42:17Z	DEBUG	events	controller-manager-578b69d9d4-t228b_521b43d8-0f30-44b1-a596-83d2946336b4 became leader	{"type": "Normal", "object": {"kind":"Lease","namespace":"system","name":"b3982890.siming.com","uid":"bd617fbe-a002-4545-8a23-ed4dd77d7f82","apiVersion":"coordination.k8s.io/v1","resourceVersion":"1484357"}, "reason": "LeaderElection"}
2024-09-17T14:42:17Z	INFO	Starting EventSource	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "source": "kind source: *v1.MySQL"}
2024-09-17T14:42:17Z	INFO	Starting EventSource	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "source": "kind source: *v1.Deployment"}
2024-09-17T14:42:17Z	INFO	Starting Controller	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL"}
2024-09-17T14:42:17Z	INFO	Starting workers	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "worker count": 1}
2024-09-17T14:43:01Z	INFO	Creating a new Deployment for MySQL	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "MySQL": {"name":"mysql-sample","namespace":"system"}, "namespace": "system", "name": "mysql-sample", "reconcileID": "90483a7a-4ba9-4381-a402-04a18635a0b9", "Deployment.Namespace": "system", "Deployment.Name": "mysql-sample"}
2024-09-17T14:45:00Z	INFO	MySQL resource not found. Ignoring since object must be deleted	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "MySQL": {"name":"mysql-sample","namespace":"system"}, "namespace": "system", "name": "mysql-sample", "reconcileID": "a93dd691-e78f-46fa-9672-fde7cf421d6d"}
2024-09-17T14:45:00Z	INFO	MySQL resource not found. Ignoring since object must be deleted	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "MySQL": {"name":"mysql-sample","namespace":"system"}, "namespace": "system", "name": "mysql-sample", "reconcileID": "10467c4f-e29d-4aca-b60d-878a3ceac98a"}
2024-09-17T14:45:07Z	INFO	Creating a new Secret for MySQL	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "MySQL": {"name":"mysql-sample","namespace":"system"}, "namespace": "system", "name": "mysql-sample", "reconcileID": "f4e005b4-b4db-49fd-a838-e9031958e9e8", "Secret.Namespace": "system", "Secret.Name": "mysql-sample"}
2024-09-17T14:45:07Z	INFO	Creating a new Deployment for MySQL	{"controller": "mysql", "controllerGroup": "apps.siming.com", "controllerKind": "MySQL", "MySQL": {"name":"mysql-sample","namespace":"system"}, "namespace": "system", "name": "mysql-sample", "reconcileID": "748ccf1b-e0ca-4544-8f11-8cd40b526458", "Deployment.Namespace": "system", "Deployment.Name": "mysql-sample"}
...

kp   
controller-manager-578b69d9d4-t228b   1/1     Running   0          3m39s
mysql-sample-64b88f58f7-9lkr9         1/1     Running   0          20s

k get mysql mysql-sample -o yaml                                                                                                      
apiVersion: apps.siming.com/v1
kind: MySQL
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps.siming.com/v1","kind":"MySQL","metadata":{"annotations":{},"labels":{"app.kubernetes.io/managed-by":"kustomize","app.kubernetes.io/name":"mysql-operator"},"name":"mysql-sample","namespace":"system"},"spec":{"backupPath":"/backups","backupSchedule":"*/5 * * * *","database":"mydatabase","password":"password123","size":1,"username":"root"}}
  creationTimestamp: "2024-09-17T14:45:07Z"
  generation: 1
  labels:
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/name: mysql-operator
  name: mysql-sample
  namespace: system
  resourceVersion: "1484746"
  uid: 3e799684-8646-4c20-a1dd-6724ab78e56b
spec:
  backupPath: /backups
  backupSchedule: '*/5 * * * *'
  database: mydatabase
  password: password123
  size: 1
  username: root
status:
  readyReplicas: 1

# 把副本改成2之后会通过operator自动创建
k apply -f samples/apps_v1_mysql.yaml
mysql.apps.siming.com/mysql-sample configured

kp -w 
NAME                                  READY   STATUS              RESTARTS   AGE
controller-manager-578b69d9d4-t228b   1/1     Running             0          5m27s
mysql-sample-64b88f58f7-9lkr9         1/1     Running             0          2m8s
mysql-sample-64b88f58f7-trwrc         0/1     ContainerCreating   0          1s
mysql-sample-64b88f58f7-trwrc         1/1     Running             0          2s
```

## Kubebuilder 原理详解

### Kubebuilder 的核心组件

**Controller-runtime 的作用**

**Controller-runtime** 是 Kubebuilder 构建 Kubernetes Operator 的核心库，它为开发者提供了一套标准化的控制器开发工具，简化了 Kubernetes 资源的管理工作。`controller-runtime` 的主要作用包括：

1. **控制器管理**：
   - `controller-runtime` 提供了一个 **Manager**，负责启动和管理所有的控制器。Manager 会管理控制器的生命周期，并确保它们持续运行。
   - 通过 `Manager`，多个控制器可以并行运行，监控和处理不同资源的状态变化。

2. **事件监听与处理**：
   - **Informer** 监听 Kubernetes 资源的变化事件（如创建、更新、删除等），`controller-runtime` 会通过这些事件驱动控制器的 **Reconcile** 函数。
   - 当资源的状态与期望不一致时，`Reconcile` 函数被触发，执行对应的逻辑进行同步。

3. **资源的缓存机制**：
   - `controller-runtime` 使用缓存本地存储资源的状态，减少与 Kubernetes API Server 的交互。控制器可以优先从缓存中获取资源信息，这提高了系统性能并减少了对 API Server 的负载。

4. **封装 Client 交互**：
   - `controller-runtime` 提供了一个简化的 **Client**，开发者可以通过它来对 Kubernetes 资源进行增删改查操作，而无需直接与 Kubernetes API Server 进行复杂的 HTTP 请求交互。
   - 例如，通过 Client 可以轻松创建或更新 Kubernetes 资源，如 Pods、Deployments、Services 等。

> **为什么 `controller-runtime` 能做这些？**
`controller-runtime` 封装了复杂的 API 调用、缓存机制和事件监听系统，开发者只需编写核心业务逻辑，它会自动处理资源的状态同步、错误处理、重试等过程。它能够通过 Manager 和 Client 管理控制器的生命周期和资源交互，使得 Operator 能够以高效、优雅的方式与 Kubernetes 生态系统交互。

**Code Generation（代码生成）的原理**

**代码生成** 是 Kubebuilder 的一个重要特性，它通过自动生成大量样板代码，减少了开发者的重复工作量。Kubebuilder 的代码生成功能基于 Go 语言的代码注解，通过注解和命令行工具生成相应的 CRD 文件、深度拷贝函数、RBAC 权限等。

1. **API 生成**：
   - 当开发者定义自定义资源 (CRD) 的结构体时，Kubebuilder 会根据这些定义自动生成 Kubernetes 所需的 CRD 清单文件以及相应的 Go 代码。例如，你定义的资源 `MySQLSpec` 和 `MySQLStatus` 会生成相应的 `yaml` 文件，用于在集群中注册 CRD。

2. **深度拷贝函数生成**：
   - Kubernetes 需要在内部对对象进行深度拷贝操作，Kubebuilder 提供了自动生成深度拷贝函数的能力。开发者只需定义资源的结构，生成工具会自动为这些结构体生成 `DeepCopy` 函数，确保对象可以安全地在不同线程和上下文中传递。

3. **RBAC 权限生成**：
   - 在控制器代码中，开发者可以通过注解（如 `+kubebuilder:rbac`）为控制器生成相应的 RBAC 配置。注解会告诉 Kubebuilder 控制器需要对哪些资源有权限，这样 Kubebuilder 会自动生成 Kubernetes 中的 RBAC 规则清单，确保控制器可以正确操作相关资源。

> **为什么 Kubebuilder 能做到？**
Kubebuilder 通过静态分析代码注解，结合 Kubernetes API 的需求，自动生成符合 Kubernetes 标准的资源清单和代码逻辑。它的底层依赖于 Kubernetes 的工具链（如 `controller-tools`）和 Go 语言的反射机制，开发者只需编写核心业务逻辑，Kubebuilder 就能自动生成复杂的样板代码。

### Kubebuilder 的工作流程

Kubebuilder 的工作流程包括初始化项目、生成 API 和控制器代码、编写业务逻辑，以及最终生成 Kubernetes 清单文件。以下是详细的工作流程。

**1. 初始化项目**

开发者通过 `kubebuilder init` 命令初始化一个标准化的 Operator 项目。这个命令会生成项目的目录结构，包括 `api/`、`controllers/`、`config/` 等文件夹。

```bash
kubebuilder init --domain mydomain.com --repo github.com/myorg/my-operator
```

- **生成的文件结构**：
  - `api/`：存放自定义资源的定义文件。
  - `controllers/`：存放控制器的逻辑。
  - `config/`：存放 Kubernetes 的清单文件（如 CRD、RBAC 配置）。

> **开发者要做什么？**  
执行 `kubebuilder init`，然后根据生成的结构填充业务逻辑。此时，项目已经具备标准的目录结构。

**2. 创建 API 和控制器**

开发者通过 `kubebuilder create api` 命令为自定义资源生成 API 结构和控制器骨架代码。

```bash
kubebuilder create api --group apps --version v1 --kind MySQL
```

- **生成的文件**：
  - `api/v1/mysql_types.go`：自定义资源 API 结构的定义。
  - `controllers/mysql_controller.go`：控制器的骨架代码。

> **开发者要做什么？**  
在 `api/v1/mysql_types.go` 中定义自定义资源的 `Spec` 和 `Status`，在 `controllers/mysql_controller.go` 中编写业务逻辑（如如何创建、更新资源）。

**3. 生成代码和配置文件**

开发者完成业务逻辑编写后，运行以下命令生成辅助代码和 Kubernetes 清单文件。

```bash
make generate
make manifests
```

- **`make generate`**：生成深度拷贝函数、API 注册等必要的辅助代码。
- **`make manifests`**：生成 CRD、RBAC、Deployment 等 Kubernetes 资源的清单文件。

> **开发者要做什么？**  
通过 `make generate` 生成自动化的辅助代码，无需手动编写深度拷贝和序列化逻辑。通过 `make manifests` 生成可部署的 CRD 和 Operator 清单文件。

**4. 编写控制器逻辑**

控制器是 Operator 的核心，负责监听和处理资源的状态变化。开发者需要在控制器的 `Reconcile` 函数中编写核心逻辑，例如当检测到 MySQL CR 创建时，如何为它创建一个相应的 Kubernetes `Deployment`。

```go
func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 获取 MySQL 实例
    var mysql appsv1alpha1.MySQL
    if err := r.Get(ctx, req.NamespacedName, &mysql); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // 检查是否需要创建或更新 Deployment
    // 处理业务逻辑
}
```

> **开发者要做什么？**  
在 `controllers/mysql_controller.go` 文件中编写核心逻辑，如资源的创建、更新和状态同步。Kubebuilder 提供的骨架代码会自动调用 `Reconcile` 函数，开发者只需专注于业务逻辑。

### 与 Kubernetes API 的集成

Kubebuilder 通过 `controller-runtime` 进行与 Kubernetes API 的交互，使用 **ClientSet**、**Informer** 和 **Controller** 协同工作，确保资源状态与期望一致。

**如何通过 ClientSet 进行资源管理**

在 `controller-runtime` 中，**Client** 是与 Kubernetes API Server 交互的主要方式。开发者可以通过 `Client` 进行资源的 CRUD 操作。

- **获取资源**：通过 `r.Get` 获取当前的资源实例。
- **创建资源**：通过 `r.Create` 向 Kubernetes 集群中创建资源。
- **更新资源**：通过 `r.Update` 更新资源状态。
- **删除资源**：通过 `r.Delete` 删除不再需要的资源。

```go
var mysql appsv1alpha1.MySQL
if err := r.Get(ctx, req.NamespacedName, &mysql); err != nil {
    return ctrl.Result{}, client.IgnoreNotFound(err)
}
```

> **开发者要做什么？**  
在控制器中使用 `Client` 来管理 Kubernetes 资源，编写 `Get`、`Create`、`Update` 和 `Delete` 的业务逻辑。`controller-runtime` 封装了与 API Server 的交互，简化了操作。

**Informer、Controller 如何配合工作**

**Informer** 监听 Kubernetes 中资源的变化（如创建、更新、删除），并将这些事件通知给相应的控制器。控制器通过 **Reconcile Loop** 响应这些事件，执行资源状态的同步和调整。

1. **Informer 的工作流程**:

   - 启动时，Informer 从 API Server 获取资源列表并缓存。
   - 通过 Watch 机制，Informer 监听资源的变化并更新缓存。
   - 当资源发生变化时，Informer 将事件通知给 Controller。

2. **Controller 的工作流程**：
   - Controller 通过 `Reconcile` 函数响应来自 Informer 的事件。每当有新的事件发生，`Reconcile` 函数被调用，执行资源的状态更新或调整。

> **开发者要做什么？**  
Kubebuilder 自动生成了 Informer 和 Controller 的配合逻辑。开发者只需在 `Reconcile` 函数中编写具体的业务逻辑，无需手动处理事件监听或缓存管理。

### 总结：开发者如何使用 Kubebuilder 开发 Operator

1. **初始化项目**：使用 `kubebuilder init` 创建标准项目结构。
2. **创建 API 和控制器**：通过 `kubebuilder create api` 生成 CRD 和控制器骨架代码。
3. **编写 API 和控制器**：定义 CRD 的 `Spec` 和 `Status`，编写控制器的 `Reconcile` 逻辑。
4. **生成代码和清单文件**：运行 `make generate` 和 `make manifests`，生成 Kubernetes 所需的清单文件和代码。
5. **部署和测试 Operator**：通过 `make install` 部署 CRD，使用 `kubectl` 验证 Operator 的工作效果。
