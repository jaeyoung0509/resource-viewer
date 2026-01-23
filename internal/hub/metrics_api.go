package hub

import (
	"github.com/jaeyoung050/resource-checker/internal/config"
	"k8s.io/client-go/kubernetes"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type metricsAPI struct {
	coreClient    *kubernetes.Clientset
	metricsClient *metricsclient.Clientset
	namespace     string
}

func newMetricsAPI(cfg config.HubConfig) (*metricsAPI, error) {
	restCfg, err := loadKubeConfig()
	if err != nil {
		return nil, err
	}

	coreClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metricsclient.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	return &metricsAPI{
		coreClient:    coreClient,
		metricsClient: metricsClient,
		namespace:     cfg.Namespace,
	}, nil
}
