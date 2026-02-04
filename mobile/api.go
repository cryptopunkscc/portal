package mobile

type Api interface {
	Config
	Status(id int32)
	Error(message string)
	StartHtml(pkg string, args string) error
	Net() Net
}

type Config interface {
	CacheDir() string
	DataDir() string
	DbDir() string
}

type Net interface {
	Addresses() (string, error)
	Interfaces() (NetInterfaceIterator, error)
}

type NetInterfaceIterator interface{ Next() *NetInterface }

type NetInterface struct {
	Index        int    // positive integer that starts at one, zero is never used
	MTU          int    // maximum transmission unit
	Name         string // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr []byte // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags        int    // e.g., FlagUp, FlagLoopback, FlagMulticast
	Addresses    string
}

// Status ID
const (
	STOPPED = int32(iota)
	STARTING
	STARTED
	FRESH
)
