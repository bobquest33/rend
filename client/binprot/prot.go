package binprot

import "encoding/binary"
import "io"
import "io/ioutil"
import "math/rand"
import "time"

import "../common"

const MAX_TTL = 3600

func init() {
    rand.Seed(time.Now().UnixNano())
}

type BinProt struct {}

// get a random expiration
func exp() uint32 {
    return uint32(rand.Intn(MAX_TTL))
}

func consumeResponse(r io.Reader) error {
    res, err := readRes(r)
    if err != nil {
        return err
    }

    apperr := statusToError(res.Status)

    // read body in regardless of the error in the header
    lr := io.LimitReader(r, int64(res.BodyLen))
    io.Copy(ioutil.Discard, lr)

    if apperr != nil && srsErr(apperr) {
        return apperr
    }
    if err == io.EOF {
        return nil
    }
    return err
}

func (b BinProt) Set(rw io.ReadWriter, key, value []byte) error {
    // set packet contains the req header, flags, and expiration
    // flags are irrelevant, and are thus zero.
    // expiration could be important, so hammer with random values from 1 sec up to 1 hour

    // Header
    bodylen := 8 + len(key) + len(value)
    writeReq(rw, common.SET, len(key), 8, bodylen)
    // Extras
    binary.Write(rw, binary.BigEndian, uint32(0))
    binary.Write(rw, binary.BigEndian, exp())
    // Body / data
    rw.Write(key)
    rw.Write(value)

    // consume all of the response and discard
    return consumeResponse(rw)
}

func (b BinProt) Get(rw io.ReadWriter, key []byte) error {
    // Header
    writeReq(rw, common.GET, len(key), 0, len(key))
    // Body
    rw.Write(key)

    // consume all of the response and discard
    return consumeResponse(rw)
}

func (b BinProt) Delete(rw io.ReadWriter, key []byte) error {
    // Header
    writeReq(rw, common.DELETE, len(key), 0, len(key))
    // Body
    rw.Write(key)

    // consume all of the response and discard
    return consumeResponse(rw)
}

func (b BinProt) Touch(rw io.ReadWriter, key []byte) error {
    // Header
    writeReq(rw, common.TOUCH, len(key), 4, len(key)+4)
    // Extras
    binary.Write(rw, binary.BigEndian, exp())
    // Body
    rw.Write(key)

    // consume all of the response and discard
    return consumeResponse(rw)
}
