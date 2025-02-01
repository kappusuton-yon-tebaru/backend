package enum

type EnvType string

const (
	Config EnvType = "config"
	Secret EnvType = "secret"
)
