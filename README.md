# juv
juv is a code generator for using [go-validator](https://github.com/go-playground/validator) in JSON Unmarshall.

## Installation

Install or upgrade Scaneo with this command.

```bash
go get -u github.com/laqiiz/juv
```

## Usage

```bash
juv [options] paths...
```

### Example

```bash
juv -o model_gen.go example
```

### Options

```console
$ juv -h
Usage of juv:
  -o string
        -o is output file name (default "juv_gen.go")
  -output string
        -output is output file name (default "juv_gen.go")
  -p string
        -p is package name (default "current directory")
  -package string
        -package is package name (default "current directory")
```

## Go Generate

If you want to use juv with go generate, then just add this comment to the top of go file.

```go
//go:generate juv $GOFILE
package models

type Login struct {
	// some fields
}
```

Now you can call go generate in package models and model_gen.go will be created.


## Motivation

go-validator is a great library, but it becomes redundant when building HTTP server as follows.

```go
func main() {
  http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
  	var login Login
    if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = w.Write([]byte(err.Error()))
        return
    }
    
    if err := validator.New().Struct(login); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, _ = w.Write([]byte(err.Error()))
        return
    }
      
      // ... some logic ...
  })
  
  log.Fatal(http.ListenAndServe(":8080", nil))
}


  // Any logic
```

It is possible to simplify by performing validation at the stage of json unmarshall.

```go
func (r *Login) UnmarshalJSON(b []byte) error {
	type Alias Login // avoid stack over flow error
	var a Alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	if err := validator.New().Struct(a); err != nil {
		return err
	}

	r.ID = a.ID
	r.Pass = a.Pass

	return nil
}
```

So 1st code can be simple.

```go
func main() {
  http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
  	  var login Login
      if err := json.NewDecoder(r.Body).Decode(&login); err != nil { // JSON Unmarshall and Validate
        w.WriteHeader(http.StatusBadRequest)
        _, _ = w.Write([]byte(err.Error()))
        return
      }  
      // ... some logic ...
  })
  
  log.Fatal(http.ListenAndServe(":8080", nil))
}
````

## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
