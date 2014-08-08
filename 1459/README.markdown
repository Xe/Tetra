# r1459
--
    import "github.com/Xe/Tetra/1459"

Package r1459 implements a base structure to scrape out and utilize an RFC 1459
frame in high level Go code.

## Usage

#### type RawLine

```go
type RawLine struct {
	Source    string
	Verb      string
	Args      []string
	Processed bool
	Raw       string
}
```

IRC line

#### func  NewRawLine

```go
func NewRawLine(input string) (line *RawLine)
```
Create a new line and split out an RFC 1459 frame to a RawLine. This will not
return an error if it fails. TODO: fix this.
