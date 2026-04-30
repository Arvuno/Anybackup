# anybackup-agent Helm Chart

这个目录现在是发布包入口 Chart。

和上次产出保持一致的部分：

- `charts/v9-infra`
- `charts/v9-services`
- `charts/core-agent-service`

本次重打包没有把 3 个子 Chart 强行拍平成单一 `templates/` 目录，而是保留为可维护的本地子 Chart 结构。这样我们既符合交付物顶层目录口径，也不需要把已经验证过的模板逻辑再拆一次。

安装方式：

```bash
helm dependency build ./helm-chart/anybackup-agent
helm upgrade --install anybackup-agent ./helm-chart/anybackup-agent -n anybackup --create-namespace
```
