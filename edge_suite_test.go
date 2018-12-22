package main

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "testing"
)

func TestEdge(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Edge Suite")
}
