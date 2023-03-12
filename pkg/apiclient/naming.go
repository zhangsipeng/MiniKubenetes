package apiclient

import "fmt"

func MangleName(nameParts ...string) string {
	mangled := ""
	for _, namePart := range nameParts {
		mangled += fmt.Sprintf("%x-%s", len(namePart), namePart)
	}
	return mangled
}
