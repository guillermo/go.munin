// Copyright 2014 Guillermo Álvarez.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license.

/*
Package client provides primitives to talk with a munin server.

It is based in net/textproto so it should support concurrent access.
*/
package client

import (
	"errors"
	"net/textproto"
	"strings"
)

var (
	// ErrMetricNotFound will be returned as an error for the Fetch and Config methods
	ErrMetricNotFound = errors.New("metric not found")
)

// A Client represents a connection to a munin Server
type Client struct {
	textproto.Conn
	Host string
}

// Dial connects to the given address on the given network using textproto.Dial
// and returns the connection
func Dial(network, addr string) (*Client, error) {
	mainConn, err := textproto.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	host, err := mainConn.ReadLine()
	if err != nil {
		return nil, err
	}
	host = strings.TrimPrefix(host, "# munin node at ")

	return &Client{*mainConn, host}, nil
}

// ListNode list the metrics for a given node. See Node on how to get a list of Nodes
func (c *Client) ListNode(node string) (metrics []string, err error) {
	id, err := c.Cmd("list " + node)
	if err != nil {
		return metrics, err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	line, err := c.ReadLine()
	if err != nil {
		return metrics, err
	}
	return strings.Split(line, " "), nil
}

// List items available for gathering
func (c *Client) List() (metrics []string, err error) {
	id, err := c.Cmd("list")
	if err != nil {
		return metrics, err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	line, err := c.ReadLine()
	if err != nil {
		return metrics, err
	}
	return strings.Split(line, " "), nil
}

// Version return the server version
func (c *Client) Version() (version string, err error) {
	id, err := c.Cmd("version")
	if err != nil {
		return "", err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	versionString, err := c.ReadLine()
	if err != nil {
		return "", err
	}
	tokens := strings.Split(versionString, " ")

	return tokens[len(tokens)-1], nil
}

// Nodes return a list of nodes
func (c *Client) Nodes() (nodes []string, err error) {
	id, err := c.Cmd("nodes")
	if err != nil {
		return nodes, err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	nodes, err = c.ReadDotLines()
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}

func (c *Client) makeFetchOrConfigCommand(cmd string) (kv map[string]string, err error) {
	id, err := c.Cmd(cmd)
	if err != nil {
		return kv, err
	}
	c.StartResponse(id)
	defer c.EndResponse(id)
	lines, err := c.ReadDotLines()
	if err != nil {
		return kv, err
	}

	if len(lines) == 1 && lines[0] == "# Unknown service" {
		return kv, ErrMetricNotFound
	}

	kv = make(map[string]string)
	for _, line := range lines {
		splitedLine := strings.SplitN(line, " ", 2)
		if len(splitedLine) == 2 {
			kv[strings.Trim(splitedLine[0], " ")] = strings.Trim(splitedLine[1], " ")
		} else if len(splitedLine) == 1 {
			kv[strings.Trim(splitedLine[0], " ")] = ""
		}
	}
	return kv, nil
}

// Config return the configuration for a specific metric
// If the metric is not found, ErrMetricNotFound is returned as error
func (c *Client) Config(metric string) (data map[string]string, err error) {
	return c.makeFetchOrConfigCommand("config " + metric)
}

// Fetch return the data asociated with one item.
// If the metric is not found, ErrMetricNotFound is returned as error
func (c *Client) Fetch(metric string) (data map[string]string, err error) {
	return c.makeFetchOrConfigCommand("fetch " + metric)
}
