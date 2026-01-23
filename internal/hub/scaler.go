package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/jaeyoung050/resource-checker/internal/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type scaler struct {
	client      *kubernetes.Clientset
	namespace   string
	maxReplicas int32
	allowed     map[string]struct{}
}

type scaleTarget struct {
	Name              string `json:"name"`
	Replicas          int32  `json:"replicas"`
	AvailableReplicas int32  `json:"availableReplicas"`
	Error             string `json:"error,omitempty"`
}

type scaleResponse struct {
	MaxReplicas int32         `json:"maxReplicas"`
	Targets     []scaleTarget `json:"targets"`
}

type scaleRequest struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replicas"`
}

func newScaler(cfg config.HubConfig) (*scaler, error) {
	if !cfg.ScaleEnabled {
		return nil, nil
	}
	if len(cfg.ScaleTargets) == 0 {
		return nil, nil
	}

	restCfg, err := loadKubeConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, err
	}

	allowed := make(map[string]struct{}, len(cfg.ScaleTargets))
	for _, name := range cfg.ScaleTargets {
		allowed[name] = struct{}{}
	}

	return &scaler{
		client:      client,
		namespace:   cfg.Namespace,
		maxReplicas: cfg.MaxReplicas,
		allowed:     allowed,
	}, nil
}

func loadKubeConfig() (*rest.Config, error) {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func (s *scaler) handleTargets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	names := make([]string, 0, len(s.allowed))
	for name := range s.allowed {
		names = append(names, name)
	}
	sort.Strings(names)

	targets := make([]scaleTarget, 0, len(names))
	for _, name := range names {
		deployment, err := s.client.AppsV1().Deployments(s.namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			targets = append(targets, scaleTarget{
				Name:  name,
				Error: err.Error(),
			})
			continue
		}
		replicas := int32(0)
		if deployment.Spec.Replicas != nil {
			replicas = *deployment.Spec.Replicas
		}
		targets = append(targets, scaleTarget{
			Name:              name,
			Replicas:          replicas,
			AvailableReplicas: deployment.Status.AvailableReplicas,
		})
	}

	writeJSON(w, scaleResponse{MaxReplicas: s.maxReplicas, Targets: targets})
}

func (s *scaler) handleScale(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req scaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	if _, ok := s.allowed[req.Name]; !ok {
		writeJSONError(w, http.StatusForbidden, errors.New("target not allowed"))
		return
	}

	if req.Replicas < 0 || req.Replicas > s.maxReplicas {
		writeJSONError(w, http.StatusBadRequest, errors.New("replicas out of range"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	patch := fmt.Sprintf(`{"spec":{"replicas":%d}}`, req.Replicas)
	updated, err := s.client.AppsV1().Deployments(s.namespace).Patch(
		ctx,
		req.Name,
		types.StrategicMergePatchType,
		[]byte(patch),
		metav1.PatchOptions{},
	)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, err)
		return
	}

	replicas := int32(0)
	if updated.Spec.Replicas != nil {
		replicas = *updated.Spec.Replicas
	}

	writeJSON(w, scaleTarget{
		Name:              req.Name,
		Replicas:          replicas,
		AvailableReplicas: updated.Status.AvailableReplicas,
	})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
