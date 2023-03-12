package WorkloadResources

import "example/Minik8s/pkg/data/ObjectMeta"

type ReplicaSet struct {
	ApiVersion string
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       ReplicaSetSpec
	Status     ReplicaSetStatus
}

type ReplicaSetSpec struct {
	Replicas int32
}

type ReplicaSetStatus struct {
	Replicas          int32
	AvailableReplicas int32
	ReadyReplicas     int32
}
