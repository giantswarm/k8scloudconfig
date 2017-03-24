package cloudconfig

type CompactTLSAssets struct {
	APIServerCACrt    string
	APIServerKey      string
	APIServerCrt      string
	WorkerCA          string
	WorkerKey         string
	WorkerCrt         string
	CalicoClientCACrt string
	CalicoClientKey   string
	CalicoClientCrt   string
	EtcdServerCACrt   string
	EtcdServerKey     string
	EtcdServerCrt     string
}
