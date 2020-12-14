package himongo

type ServerVersion struct {
	version string
}

var ServerVersion1 ServerVersion = ServerVersion{"1"}

func (s ServerVersion) String() string {
	return s.version
}
