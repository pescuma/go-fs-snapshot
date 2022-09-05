package cli

type serverArgs struct {
	Server               string `help:"Server to connect to, in the format ip:port"`
	ServerOnlyAsFallback bool   `help:"Use server only as fallback. This only applies if --server is used."`
}
