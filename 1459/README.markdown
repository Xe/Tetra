# r1459
--
    import "github.com/Xe/Tetra/1459"

Package r1459 implements a base structure to scrape out and utilize an RFC 1459
frame in high level Go code.

## Usage

#### type RawLine

```go
type RawLine struct {
	Source string            `json: "source"`
	Verb   string            `json:"verb"`
	Args   []string          `json:"args"`
	Tags   map[string]string `json:"tags"`
	Raw    string            `json:"-"` // Deprecated
}
```

IRC line

#### func  NewRawLine

```go
func NewRawLine(input string) (line *RawLine)
```
Create a new line and split out an RFC 1459 frame to a RawLine. This will not
return an error if it fails. TODO: fix this.

#### func (*RawLine) String

```go
func (r *RawLine) String() (res string)
```
String returns the serialized form of a RawLine as an RFC 1459 frame.
