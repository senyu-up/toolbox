# jaeger本地部署
- 注意启动顺序，先启动es，再启动其他部件

# jaeger生产环境部署
- es需要是集群
- 每个物理机部署一个jaeger agent 或者 k8s 一个node一个jaeger agent
- jaeger部件崩掉不会影响正常程序执行
- 定期删除es过时数据