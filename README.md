# juv
juv is JSON Validator in Unmarshall phase

## Motivation

go-validator is a great library, but it becomes redundant when building HTTP server as follows.

```go
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
  var signUp SignUp
  if err := json.NewDecoder(r.Body).Decode(&signUp); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    _, _ = w.Write([]byte(err.Error()))
    return
  }

  if err := validator.New().Struct(signUp); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    _, _ = w.Write([]byte(err.Error()))
    return
  }

  // Any logic
```

It is possible to simplify by performing validation at the stage of json unmarshall.

```go
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
  var signUp SignUp
  if err := json.NewDecoder(r.Body).Decode(&signUp); err != nil { // Unmarshall and Validation
    w.WriteHeader(http.StatusBadRequest)
    _, _ = w.Write([]byte(err.Error()))
    return
  }
  // Any logic
```


## Spec

TODO


## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
