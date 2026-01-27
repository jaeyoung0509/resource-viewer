package hub

import (
	"context"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (m *metricsAPI) collectPodMetrics(ctx context.Context) ([]podMetric, error) {
	podMetrics, err := m.metricsClient.MetricsV1beta1().PodMetricses(m.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	items := make([]podMetric, 0, len(podMetrics.Items))
	for _, pod := range podMetrics.Items {
		cpuMilli, memBytes := sumContainers(pod.Containers)
		items = append(items, podMetric{
			Name:        pod.Name,
			Namespace:   pod.Namespace,
			CPUMilli:    cpuMilli,
			MemoryBytes: memBytes,
			Timestamp:   pod.Timestamp.Unix(),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func (m *metricsAPI) collectDeploymentMetrics(ctx context.Context) ([]deploymentMetric, error) {
	deployments, err := m.coreClient.AppsV1().Deployments(m.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := m.coreClient.CoreV1().Pods(m.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podMetrics, err := m.metricsClient.MetricsV1beta1().PodMetricses(m.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	metricsByPod := make(map[string]metricsv1beta1.PodMetrics, len(podMetrics.Items))
	for _, pod := range podMetrics.Items {
		metricsByPod[pod.Name] = pod
	}

	items := make([]deploymentMetric, 0, len(deployments.Items))
	for _, deploy := range deployments.Items {
		if deploy.Spec.Selector == nil {
			continue
		}
		selector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
		if err != nil {
			return nil, err
		}

		var cpuMilli int64
		var memBytes int64
		var podCount int
		var missing int

		for _, pod := range pods.Items {
			if !selector.Matches(labels.Set(pod.Labels)) {
				continue
			}
			podCount++
			metric, ok := metricsByPod[pod.Name]
			if !ok {
				missing++
				continue
			}
			cpu, mem := sumContainers(metric.Containers)
			cpuMilli += cpu
			memBytes += mem
		}

		replicas := int32(0)
		if deploy.Spec.Replicas != nil {
			replicas = *deploy.Spec.Replicas
		}

		items = append(items, deploymentMetric{
			Name:           deploy.Name,
			Replicas:       replicas,
			Available:      deploy.Status.AvailableReplicas,
			CPUMilli:       cpuMilli,
			MemoryBytes:    memBytes,
			PodCount:       podCount,
			MissingMetrics: missing,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}

func sumContainers(containers []metricsv1beta1.ContainerMetrics) (int64, int64) {
	var cpuMilli int64
	var memBytes int64
	for _, container := range containers {
		cpuMilli += container.Usage.Cpu().MilliValue()
		memBytes += container.Usage.Memory().Value()
	}
	return cpuMilli, memBytes
}
