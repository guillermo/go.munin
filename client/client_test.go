// Copyright 2014 Guillermo Álvarez.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license.

package client

import (
	"testing"
)

const (
	Hostname     = "blanquito.cientifico.net"
	MuninVersion = "2.0.19-3"
)

func TestClient(t *testing.T) {
	client, err := Dial("tcp4", "127.0.0.1:4949")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Test List
	metrics, err := client.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 42 {
		t.Error(metrics)
	}

	// Test ListNode
	metrics, err = client.ListNode(Hostname)
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 42 {
		t.Error(metrics)
	}

	// Test Nodes
	nodes, err := client.Nodes()
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Error(len(nodes))
	}
	if nodes[0] != Hostname {
		t.Error(nodes[0])
	}

	// Test Config
	config, err := client.Config("load")
	if err != nil {
		t.Fatal(err)
	}
	if len(config) != 8 {
		t.Fatal(len(config))
	}

	config, err = client.Config("tiopepe")
	if err != ErrMetricNotFound {
		t.Fatal(err, config)
	}

	// Test Fetch
	data, err := client.Fetch("load")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1 {
		t.Fatal(len(data))
	}

	config, err = client.Fetch("tiopepe")
	if err != ErrMetricNotFound {
		t.Fatal(err, config)
	}

	// Test Version
	version, err := client.Version()
	if err != nil {
		t.Fatal(err)
	}
	if version != MuninVersion {
		t.Error(err)
	}

}
