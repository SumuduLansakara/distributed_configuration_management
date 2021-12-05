start_cluster:
	HostIP="192.168.12.50" \
	docker run \
		-d \
		-v /usr/share/ca-certificates/:/etc/ssl/certs \
		-p 14001:14001 \
		-p 12380:12380 \
		-p 12379:12379 \
 		--name etcd quay.io/coreos/etcd:v2.3.8 \
 		-name etcd0 \
 		-advertise-client-urls http://${HostIP}:12379,http://${HostIP}:14001 \
 		-listen-client-urls http://0.0.0.0:12379,http://0.0.0.0:14001 \
 		-initial-advertise-peer-urls http://${HostIP}:12380 \
 		-listen-peer-urls http://0.0.0.0:12380 \
 		-initial-cluster-token etcd-cluster-1 \
 		-initial-cluster etcd0=http://${HostIP}:12380 \
 		-initial-cluster-state new

cleanup:
	docker container rm -f etcd

test_cluster:
	# etcdctl --endpoints http://192.168.12.50:12379 member list
	etcdctl --endpoints http://localhost:12379 member list

.EXPORT_ALL_VARIABLES:
HostIP="192.168.12.50"
