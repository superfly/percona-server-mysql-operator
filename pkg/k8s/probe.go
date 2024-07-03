package k8s

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ExecProbe(probe corev1.Probe, cmd []string) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{Command: cmd},
		},
		InitialDelaySeconds:           probe.InitialDelaySeconds,
		TimeoutSeconds:                probe.TimeoutSeconds,
		PeriodSeconds:                 probe.PeriodSeconds,
		FailureThreshold:              probe.FailureThreshold,
		SuccessThreshold:              probe.SuccessThreshold,
		TerminationGracePeriodSeconds: probe.TerminationGracePeriodSeconds,
	}
}

func HTTPCheckProbe(probe corev1.Probe, path string, port int) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: path,
				Port: intstr.FromInt(port),
			},
		},
		InitialDelaySeconds:           probe.InitialDelaySeconds,
		TimeoutSeconds:                probe.TimeoutSeconds,
		PeriodSeconds:                 probe.PeriodSeconds,
		FailureThreshold:              probe.FailureThreshold,
		SuccessThreshold:              probe.SuccessThreshold,
		TerminationGracePeriodSeconds: probe.TerminationGracePeriodSeconds,
	}
}
