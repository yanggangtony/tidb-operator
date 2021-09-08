// Copyright PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TidbInitializerLister helps list TidbInitializers.
// All objects returned here must be treated as read-only.
type TidbInitializerLister interface {
	// List lists all TidbInitializers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TidbInitializer, err error)
	// TidbInitializers returns an object that can list and get TidbInitializers.
	TidbInitializers(namespace string) TidbInitializerNamespaceLister
	TidbInitializerListerExpansion
}

// tidbInitializerLister implements the TidbInitializerLister interface.
type tidbInitializerLister struct {
	indexer cache.Indexer
}

// NewTidbInitializerLister returns a new TidbInitializerLister.
func NewTidbInitializerLister(indexer cache.Indexer) TidbInitializerLister {
	return &tidbInitializerLister{indexer: indexer}
}

// List lists all TidbInitializers in the indexer.
func (s *tidbInitializerLister) List(selector labels.Selector) (ret []*v1alpha1.TidbInitializer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TidbInitializer))
	})
	return ret, err
}

// TidbInitializers returns an object that can list and get TidbInitializers.
func (s *tidbInitializerLister) TidbInitializers(namespace string) TidbInitializerNamespaceLister {
	return tidbInitializerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TidbInitializerNamespaceLister helps list and get TidbInitializers.
// All objects returned here must be treated as read-only.
type TidbInitializerNamespaceLister interface {
	// List lists all TidbInitializers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TidbInitializer, err error)
	// Get retrieves the TidbInitializer from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.TidbInitializer, error)
	TidbInitializerNamespaceListerExpansion
}

// tidbInitializerNamespaceLister implements the TidbInitializerNamespaceLister
// interface.
type tidbInitializerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TidbInitializers in the indexer for a given namespace.
func (s tidbInitializerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TidbInitializer, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TidbInitializer))
	})
	return ret, err
}

// Get retrieves the TidbInitializer from the indexer for a given namespace and name.
func (s tidbInitializerNamespaceLister) Get(name string) (*v1alpha1.TidbInitializer, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("tidbinitializer"), name)
	}
	return obj.(*v1alpha1.TidbInitializer), nil
}
