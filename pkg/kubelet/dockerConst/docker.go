package dockerConst

const NetworkName = "minik8s-docker-bridge"
const PauseImage = "minik8s/pause"

const (
	LabelMinik8sKey = "k8s"
	LabelMinik8sVal = "minik8s"
	LabelMinik8s    = LabelMinik8sKey + "=" + LabelMinik8sVal
)
