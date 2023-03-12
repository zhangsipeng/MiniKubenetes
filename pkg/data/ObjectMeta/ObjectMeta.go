package ObjectMeta

type ObjectMeta struct {
	Name         string
	GenerateName string
	Namespace    string
	Labels       map[string]string
	Annotations  map[string]string
}
