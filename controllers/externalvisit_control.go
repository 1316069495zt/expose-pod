package controllers

import appsv1alpha1 "external/api/v1alpha1"

type commonControl struct {
	*appsv1alpha1.ExternalvisitSet
}

func (c *commonControl) IsActiveExternalvisitSet() bool {
	return true
}

func New(cs *appsv1alpha1.ExternalvisitSet) *commonControl {
	return &commonControl{ExternalvisitSet: cs}
}
