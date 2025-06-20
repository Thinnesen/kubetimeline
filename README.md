# kubetimeline 🚀⏳

A modern, Kubernetes-native timeline and event-tracking operator! 

## ✨ Overview

**kubetimeline** brings visibility and traceability to your Kubernetes resources by providing a timeline of key events, changes, and custom milestones. Whether you're debugging, auditing, or just curious about the lifecycle of your workloads, kubetimeline helps you answer: _"What happened, when, and why?"_ 🕵️‍♂️📅

---

## 🛠️ Features (WIP)

- 📜 **Resource Timeline**: Visualize the history of any Kubernetes resource (Pods, Deployments, CRDs, etc.)
- 🔔 **Event Aggregation**: Collects and correlates events from multiple sources (Kubernetes events, controllers, custom hooks)
- 🧩 **Custom Milestones**: Define your own events or milestones for resources
- 🕰️ **Historical Audit**: See what changed, who changed it, and when
- 🖼️ **Web UI (Planned)**: Beautiful, interactive timeline dashboard (coming soon!)
- 🔗 **API Integration (Planned)**: Query timelines programmatically for automation and reporting
- 🛡️ **RBAC-Aware**: Secure by default, respects Kubernetes RBAC
- 🏗️ **Extensible**: Easily add new event sources or output formats

> **Note:** Many features are in active development! Contributions and feedback are welcome. See [Contributing](#contributing) below. 🚧

---

## 💡 Use Cases

- 🐛 **Debugging**: Quickly see the sequence of events leading to a failure
- 🔍 **Auditing**: Track who did what, and when, across your cluster
- 📈 **Change Tracking**: Visualize deployments, rollouts, and config changes over time
- 🧪 **Testing**: Validate that your controllers and webhooks emit the right events
- 🛠️ **Custom Workflows**: Integrate with CI/CD or incident response pipelines

---

## 🚀 Getting Started

### Prerequisites
- Go v1.24.0+
- Docker 17.03+
- kubectl v1.11.3+
- Access to a Kubernetes v1.11.3+ cluster

### Deploy to Your Cluster

1. **Build and push your image:**

   ```sh
   make docker-build docker-push IMG=<some-registry>/kubetimeline:tag
   ```

2. **Install the CRDs:**

   ```sh
   make install
   ```

3. **Deploy the Manager:**

   ```sh
   make deploy IMG=<some-registry>/kubetimeline:tag
   ```

4. **Create sample instances:**

   ```sh
   kubectl apply -k config/samples/
   ```

> ⚠️ If you encounter RBAC errors, ensure you have cluster-admin privileges.

---

### 🧹 Uninstall

- **Delete sample instances:**
  ```sh
  kubectl delete -k config/samples/
  ```
- **Delete CRDs:**
  ```sh
  make uninstall
  ```
- **Undeploy the controller:**
  ```sh
  make undeploy
  ```

---

## 📦 Project Distribution

### Option 1: YAML Bundle

1. **Build the installer:**
   ```sh
   make build-installer IMG=<some-registry>/kubetimeline:tag
   ```
   > Generates `dist/install.yaml` for easy installation.

2. **Install with kubectl:**
   ```sh
   kubectl apply -f https://raw.githubusercontent.com/<org>/kubetimeline/<tag or branch>/dist/install.yaml
   ```

### Option 2: Helm Chart

1. **Build the chart:**
   ```sh
   kubebuilder edit --plugins=helm/v1-alpha
   ```
2. **Find the chart in `dist/chart` and install as usual.**

> **Note:** Update the chart after changes. For webhooks, use `--force` and manually re-apply custom config.

---

## 🤝 Contributing

We welcome PRs, issues, and feature requests! 
- See [Kubebuilder Docs](https://book.kubebuilder.io/introduction.html) for operator development tips
- Run `make help` for all available targets
- Check the TODOs in this README for areas needing help

---

## 📄 License

Apache License 2.0. See [LICENSE](LICENSE) for details.

---

> _kubetimeline is a work in progress. Star ⭐ the repo to follow updates!_

