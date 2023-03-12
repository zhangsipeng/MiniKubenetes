package WorkloadResources

import (
	"example/Minik8s/pkg/data/ObjectMeta"
	"time"
)

type HorizontalPodAutoscaler struct {
	ApiVersion string
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       HorizontalPodAutoscalerSpec
	Status     HorizontalPodAutoscalerStatus
}

type HorizontalPodAutoscalerSpec struct {
	MaxReplicas                    int32
	MinReplicas                    int32
	ScaleTargetRef                 CrossVersionObjectReference
	TargetCPUUtilizationPercentage int32
}

type HorizontalPodAutoscalerStatus struct {
	CurrentReplicas                 int32
	DesiredReplicas                 int32
	CurrentCPUUtilizationPercentage int32
	LastScaleTime                   time.Time
}

type CrossVersionObjectReference struct {
	Kind       string
	Name       string
	ApiVersion string
}
