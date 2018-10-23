# kubectl-evict
kubectl plugin to evict pods

This requires a k8s cluster that supports the eviction API.
Build this binary and place it on your PATH somewhere.
Use it via `kubectl evict PODNAME` or `kubectl evict -namespace NAMESPACE PODNAME`.
