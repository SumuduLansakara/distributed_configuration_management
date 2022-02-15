start_cluster:
	HostIP="192.168.12.50" \
	docker run \
 		--name etcd \
		-d \
		--publish ${PORT1}:2379 \
		--publish ${PORT2}:2380 \
		--env ALLOW_NONE_AUTHENTICATION=yes \
		--env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
		bitnami/etcd:latest

cleanup:
	docker container rm -f etcd

test_cluster:
	etcdctl --endpoints http://localhost:${PORT1} member list

.EXPORT_ALL_VARIABLES:
PORT1=12379
PORT2=12380
