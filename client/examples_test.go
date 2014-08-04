// Copyright 2014 Guillermo Álvarez.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license.

package client

import (
	"fmt"
	"strings"
)

// Fetch a specific metric
func ExampleClient_Fetch() {
	client, err := Dial("tcp4", "127.0.0.1:4949")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	data, err := client.Fetch("load")
	if err != nil {
		panic(err)
	}

	for k, v := range data {
		fmt.Println(k, v)
	}

	// Output:
	// load.value 0.19
}

// List all the metrics supported by a host
func ExampleClient_List() {
	client, err := Dial("tcp4", "127.0.0.1:4949")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	metrics, err := client.List()
	if err != nil {
		panic(err)
	}

	fmt.Printf(strings.Join(metrics, "\n"))

	// Output:
	// cpu
	// load
	// users
	// processes
}

// Get the configuration of metric
func ExampleClient_Config() {
	client, err := Dial("tcp4", "127.0.0.1:4949")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	config, err := client.Config("load")
	if err != nil {
		panic(err)
	}

	for k, v := range config {
		fmt.Println(k, v)
	}

	// Output:
	// graph_title Load average
	// graph_args --base 1000 -l 0
	// graph_vlabel load
	// graph_scale no
	// graph_category system
	// load.label load
	// graph_info The load average of the machine describes how many processes are in the run-queue (scheduled to run "immediately").
	// load.info 5 minute load average

}
