package main

import (
  "net/http"
  "net/http/httptest"
  "github.com/labstack/echo"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("EdgeController", func() {
  Describe("GET /", func() {
      Context("With a valid request", func() {
          It("should return successful response", func() {
              e := echo.New()
              req, _ := http.NewRequest(http.MethodGet, "/", nil)
              rec := httptest.NewRecorder()
              c := e.NewContext(req, rec)

              HelloController(c)
              Expect(rec.Code).To(Equal(http.StatusOK))
              Expect(rec.Body.String()).To(Equal("Hello, World!\n"))
          })
      })
  })
})
