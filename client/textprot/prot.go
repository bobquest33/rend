package textprot

import "fmt"
import "io"
import "strings"

const VERBOSE = false

// reads a line without needing to use a bufio.Reader
func rl(r io.Reader) (string, error) {
    retval := make([]byte, 1024)
    b := make([]byte, 1)
    cur := 0

    for b[0] != '\n' {
        n, err := r.Read(b)
        if err != nil {
            return "", err
        }
        if n < 1 {
            continue
        }

        retval[cur] = b[0]
        cur++

        if cur >= cap(retval) {
            newretval := make([]byte, 2*len(retval))
            copy(newretval, retval)
            retval = newretval
        }
    }

    return string(retval[:cur]), nil
}

type TextProt struct {}

func (t TextProt) Set(rw io.ReadWriter, key []byte, value []byte) error {
    strKey := string(key)
    if VERBOSE { fmt.Printf("Setting key %v to value of length %v\r\n", strKey, len(value)) }

    _, err := fmt.Fprintf(rw, "set %v 0 0 %v\r\n", strKey, len(value))
    if err != nil { return err }
    _, err = fmt.Fprintf(rw, "%v\r\n", string(value))
    if err != nil { return err }
    
    response, err := rl(rw)
    if err != nil { return err }

    if VERBOSE {
        fmt.Println(response)
        fmt.Printf("Set key %v\r\n", strKey)
    }

    return nil
}

func (t TextProt) Get(rw io.ReadWriter, key []byte) error {
    strKey := string(key)
    if VERBOSE { fmt.Printf("Getting key %v\r\n", strKey) }

    _, err := fmt.Fprintf(rw, "get %v\r\n", strKey)
    if err != nil { return err }
    
    // read the header line
    response, err := rl(rw)
    if err != nil { return err }
    if VERBOSE { fmt.Println(response) }
    
    if strings.TrimSpace(response) == "END" {
        if VERBOSE { fmt.Println("Empty response / cache miss") }
        return nil
    }

    // then read the value
    response, err = rl(rw)
    if err != nil { return err }
    if VERBOSE { fmt.Println(response) }

    // then read the END
    response, err = rl(rw)
    if err != nil { return err }
    if VERBOSE { fmt.Println(response) }
    
    if VERBOSE { fmt.Printf("Got key %v\r\n", key) }
    return nil
}

func (t TextProt) Delete(rw io.ReadWriter, key []byte) error {
    strKey := string(key)
    if VERBOSE { fmt.Printf("Deleting key %s\r\n", strKey) }
    
    _, err := fmt.Fprintf(rw, "delete %s\r\n", strKey)
    if err != nil { return err }
    
    response, err := rl(rw)
    if err != nil { return err }
    if VERBOSE { fmt.Println(response) }
    
    if VERBOSE { fmt.Printf("Deleted key %s\r\n", strKey) }
    return nil
}

func (t TextProt) Touch(rw io.ReadWriter, key []byte) error {
    strKey := string(key)
    if VERBOSE { fmt.Printf("Touching key %s\r\n", strKey) }
    
    _, err := fmt.Fprintf(rw, "touch %s 123456\r\n", strKey)
    if err != nil { return err }

    response, err := rl(rw)
    if err != nil { return err }
    if VERBOSE { fmt.Println(response) }
    
    if VERBOSE { fmt.Printf("Touched key %s\r\n", strKey) }
    return nil
}
