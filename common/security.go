package common

// SecurityOptions configures document protection and access control metadata.
type SecurityOptions struct {
	Password   string
	Author     string
	AllowPrint bool
	AllowCopy  bool
	AllowEdit  bool
}
