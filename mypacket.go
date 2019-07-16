package main

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "reflect"
    "strconv"
    "strings"
    //"compress/zlib"
)

var ErrReadStream =  errors.New("Read Stream Error")
var ErrCliCommand =  errors.New("Read Client Packet Command Type Error")
var ErrParseData  =  errors.New("Parse Data Error")
var ErrUnknownType =  errors.New("Unknown Payload Type")


type Packet struct {
    length      uint64
    sequence    uint64
    payload     []byte
}


type MyConn struct {
    src_ip      string
    dst_ip      string
    src_port    string
    dst_port    string
    client      bool
    //TODO: support compressed packets and parse fields depends on capacity flags
    //compressed  bool
    capacity    uint64
    packets     chan *Packet
}

type Payload interface {
    String() string
}

func (mc MyConn) String() string {
    return fmt.Sprintf("%s:%s => %s:%s", mc.src_ip, mc.src_port, mc.dst_ip, mc.dst_port)
}


func parseFixUint(data []byte) (num uint64) {
    for i, b := range data {
        num |= uint64(b) << (8 * uint(i))
    }

    return
}


func parseLenEncUint(data []byte) (num uint64, byte_cnt int) {
    switch data[0] {

    default:
        num = parseFixUint(data[0:1])
        byte_cnt = 1
        return

    case 0xfb:
        num = 256
        byte_cnt = 1
        return

    case 0xfc:
        byte_cnt = 2
    case 0xfd:
        byte_cnt = 3
    case 0xfe:
        byte_cnt = 8
    }

    num = parseFixUint(data[1:byte_cnt+1])
    return

}

func parseFixString(data []byte) string {
    return string(data)
}

func parseStringNUL(data []byte) (string, int) {
    for i, b := range data {
        if b == 0x00 {
            return string(data[:i]), i+1
        }
    }

    return "", -1
}

func parseStringEOF(data []byte) (string, int) {
    var s []byte
    var last_ok bool
    for _, b := range data {
        // only printable bytes returned
        if 0x20 < b && b < 0x7f {
            s = append(s, b)
            last_ok = true
            continue
        } else if last_ok {
            // add space between continuous characters
            s = append(s, 0x20)
        }
        last_ok = false
    }
    return string(s), len(data)
}

func parseLenEncStringEOF(data []byte) (string, int) {

    var offset, next int
    var s []string

    total_len := len(data)
    for offset < total_len {
        v, dlen := parseLenEncUint(data[offset:])
        next = offset + dlen
        offset = next

        if v != 0 && dlen == 0 || next > total_len {
            log.Errorf("Parse parseLenEncStringEOF Error: [%s]", data[offset:])
            return "", 0
        } else if v == 256 && dlen == 1 {
            s = append(s, "NULL")
        } else {
            next = offset + int(v)
            s = append(s, parseFixString(data[offset:next]))
            offset = next
        }
    }
    return strings.Join(s, ", "), offset

}

func parseStringNULEOF(data []byte) (string, int) {

    var offset, next int
    var s []string

    total_len := len(data)
    for offset < total_len {
        v, dlen := parseStringNUL(data[offset:])
        if offset == 0 && dlen == -1 {
            return "", total_len
        } else if dlen == -1 {
            return strings.Join(s, ", "), total_len
        }

        next = offset + dlen
        s = append(s, v)
        offset = next

    }
    return strings.Join(s, ", "), total_len

}



/* TODO: support compressed packets
// assume first packet is compressed, get 3 bytes more,
// set compressed to true, send packet(s) to channel,
// if not compressed then send first 2 packets since 3 bytes already read
func (conn *MyConn) guessCompressed(r io.Reader) error {
    // 3 compressed length bytes
    // 1 sequence byte
    // 3 original length bytes
    header := make([]byte, 7)
    if n, err := io.ReadFull(r, header); err != nil {
        logger.Error("Read Stream %s", err)
        return err
    }

    length_bytes, sequence_byte, orig_length_bytes := header[:3], header[3:4], header[4:]
    body_length := parseFixUint(length_bytes)
    sequence_id := parseFixUint(sequence_byte)
    orig_length := parseFixUint(orig_length_bytes)

    var body_buffer bytes.Buffer
    if n, err := io.CopyN(body_buffer, r, body_length); err != nil {
        logger.Error("Read Stream %s", err)
        return err
    } else if n != body_length {
        logger.Error("Read Stream Packet Length Error")
        return errors.New("Packet Length Error")
    }

    zr, err := zlib.NewReader(body_buffer)
    // uncompressed payload, but we read 3 more bytes...
    if err != nil {
        bytes_read := body_buffer.Bytes()
        p1_first_3bytes := header[4:]
        p1_surplus_bytes := bytes_read[:body_length-3]
        p1 := &payload{
            length = body_length,
            sequence = sequence_id,
            body := append(p1_first_3bytes, p1_surplus_bytes...),
        }
        conn.payloads <- &p1

        p2_length_bytes := bytes_read[body_length-3:]
        p2_body_length := parseFixUint(p2_length_bytes)
        var p2_seq_body bytes.Buffer

        // read sequence & body without length bytes
        if n, err := io.CopyN(p2_seq_body, r, body_length+1); err != nil {
            logger.Error("Read Stream %s", err)
            return err
        } else if n != body_length+1 {
            logger.Error("Read Stream Packet Length Error")
            return errors.New("Packet Length Error")
        }
        p2_seq_body_bytes := p2_seq_body.Bytes()
        p2_sequence_byte, p2_body_bytes := p2_seq_body_bytes[0:1], p2_seq_body_bytes[1:]
        p2_sequence_id := parseFixUint(p2_sequence_byte)

        p2 := &payload{
            length: p2_body_length,
            sequence: p2_sequence_id,
            body: p2_body_bytes,
        }
        conn.payloads <- &p2
    }

    // maybe more than 1 packet were compressed
    defer zr.Close()
    conn.compressed = true
    conn.sendPayload(zr)
}

// decompress r to plain reader
func decompress(r io.Reader) (io.ReadCloser, error) {
    header := make([]byte, 7)
    if n, err := io.ReadFull(r, header); err != nil {
        logger.Error("Read Stream %s", err)
        return nil, err
    }

    length_bytes, sequence_byte, orig_length_bytes := header[:3], header[3:4], header[4:]
    body_length := parseFixUint(length_bytes)

    var body_buffer bytes.Buffer
    if n, err := io.CopyN(body_buffer, r, body_length); err != nil {
        logger.Error("Decompress Stream %s", err)
        return nil, err
    } else if n != body_length {
        logger.Error("Decompress Stream Packet Length Error")
        return nil, errors.New("Decompress Packet Length Error")
    }

    zr, err := zlib.NewReader(body_buffer)
    if err != nil {
        // maybe payload less than MIN_COMPRESS_LENGTH, 
        // this payload is uncompressed, we don't wanna guess again,
        bytes_read := body_buffer.Bytes()
        p2_length_bytes := bytes_read[body_length-3:]
        p2_body_length := parseFixUint(p2_length_bytes)

        // discard the 2nd packet
        if n, err := io.CopyN(ioutil.Discard, r, p2_body_length+1); err != nil {
            logger.Error("Decompress Stream Discard %s", err)
            return nil, err
        }
        return nil, err
    }

    var r bytes.Buffer
    r := io.Copy()
    return &zr, nil

}
*/

// TODO: parse with a state machine
func (conn *MyConn) ParseCliPayload(packet *Packet) (payload Payload, err error) {

    var offset, next int
    var com interface{}

    // MySQL Packet Sequence starts from client command
    if packet.sequence == 0 {
        com_type := packet.payload[0]
        com_fac, ok := CliComTypes[com_type]
        if !ok {
            log.Infof("[%s] Unknown command type [%v]", conn, com_type)
            return nil, ErrUnknownType
        }
        com = com_fac()
    // Auth Response Packet has no specific header
    } else if packet.sequence == 1 && packet.length > 32 && bytes.Equal(packet.payload[9:32], Filler23Bytes) {
        cap_flags := parseFixUint(packet.payload[:4])
        conn.capacity = cap_flags
        com = &ComHandShake{}
    } else {
            log.Infof("[%s] Unknown command type [%v], Packet sequence [%d]", conn, packet.sequence)
            return nil, ErrUnknownType
    }

    irt := reflect.TypeOf(com).Elem()

    // Client Command Types need pointer reciever while type assert
    rv := reflect.ValueOf(com)
    irv := rv.Elem()
    nf := irt.NumField()
    for i := 0; i < nf && uint64(offset) < packet.length; i++ {
        ft := irt.Field(i)
        f := irv.Field(i)
        dt := ft.Tag.Get("datatype")
        switch dt {
            case "FixBytes":
                dlen, _ := strconv.Atoi(ft.Tag.Get("length"))
                next = offset + dlen
                fv := packet.payload[offset:next]
                f.SetBytes(fv)
                offset = next

            case "LenEncBytes":
                v, dlen := parseLenEncUint(packet.payload[offset:])
                next = offset + dlen
                if v != 0 && dlen == 0 {
                    log.Errorf("[%s] Parse LenEncBytes Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                offset = next

                // slen here is not a struct value but following string length
                next = offset + int(v)
                fv := packet.payload[offset:next]
                f.SetBytes(fv)
                offset = next

            case "FixInt":
                dlen, _ := strconv.Atoi(ft.Tag.Get("length"))
                next = offset + dlen
                fv := parseFixUint(packet.payload[offset:next])
                f.SetUint(fv)
                offset = next

            case "LenEncUint":
                fv, dlen := parseLenEncUint(packet.payload[offset:])
                next = offset + dlen
                if fv != 0 && dlen == 0 {
                    log.Errorf("[%s] Parse LenEncInt Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetUint(fv)
                offset = next

            case "StringNUL":
                fv, dlen := parseStringNUL(packet.payload[offset:])
                next = offset + dlen
                if dlen < 0 {
                    log.Errorf("[%s] Parse StringNUL Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetString(fv)
                offset = next

            case "StringEOF":
                fv, dlen := parseStringEOF(packet.payload[offset:])
                next = offset + dlen
                f.SetString(fv)
                offset = next

            case "StringNULEOF":
                fv, dlen := parseStringNULEOF(packet.payload[offset:])
                next = offset + dlen
                if dlen == 0 {
                    log.Errorf("[%s] Parse StringNULEOF Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetString(fv)
                offset = next

            default:
                log.Errorf("[%s] Unkown Data Type: [%s]", conn, dt)
                return nil, ErrParseData
        }

    }

    if uint64(offset) != packet.length {
        log.Errorf("[%s]: %s Packet Length: [%d], Parsed Data Length: [%d]", conn, ErrParseData, packet.length, offset)
        //return nil, ErrParseData
    }

    return rv.Interface().(Payload), nil

}

func (conn *MyConn) ParseSrvPayload(packet *Packet) (payload Payload, err error) {
    var offset, next int
    var resp interface{}
    switch {
    case packet.sequence == 0 && packet.payload[0] == 0x0a:
        _, dlen := parseLenEncUint(packet.payload[1:])
        cap_bytes_p1 := packet.payload[dlen+14:dlen+16]
        cap_bytes_p2 := packet.payload[dlen+19:dlen+21]
        cap_bytes := append(cap_bytes_p1, cap_bytes_p2...)
        cap_flags := parseFixUint(cap_bytes)
        conn.capacity = cap_flags
        resp = &RespHandShake{}

    case (packet.payload[0] == 0x00 || packet.payload[0] == 0xfe) && packet.length >= 7:
        resp = &RespOK{}

    case packet.payload[0] == 0xfe && packet.length < 7:
        resp = &RespEOF{}

    case packet.payload[0] == 0xff:
        resp = &RespErr{}

    case packet.sequence == 1:
        // maybe a column count packet, whose payload contains only columns count
        _, dlen := parseLenEncUint(packet.payload)
        if int(packet.length) == dlen {
            resp = &RespColCnt{}
        } else {
            resp = &RespStatistics{}
        }

    case packet.sequence > 1 && packet.length > 4 && bytes.Equal(packet.payload[0:4], ColDefHeader):
        resp = &RespColDef{}

    default:
        resp = &RespResultSet{}
    }

    irt := reflect.TypeOf(resp).Elem()
    rv := reflect.ValueOf(resp)
    irv := rv.Elem()
    nf := irt.NumField()
    for i := 0; i < nf && uint64(offset) < packet.length; i++ {
        ft := irt.Field(i)
        f := irv.Field(i)
        dt := ft.Tag.Get("datatype")

        switch dt {
            case "FixBytes":
                dlen, _ := strconv.Atoi(ft.Tag.Get("length"))
                next = offset + dlen
                fv := packet.payload[offset:next]
                f.SetBytes(fv)
                offset = next

            case "FixInt":
                dlen, _ := strconv.Atoi(ft.Tag.Get("length"))
                next = offset + dlen
                fv := parseFixUint(packet.payload[offset:next])
                f.SetUint(fv)
                offset = next

            case "FixString":
                dlen, _ := strconv.Atoi(ft.Tag.Get("length"))
                next = offset + dlen
                fv := parseFixString(packet.payload[offset:next])
                f.SetString(fv)
                offset = next

            case "LenEncUint":
                fv, dlen := parseLenEncUint(packet.payload[offset:])
                next = offset + dlen
                if fv != 0 && dlen == 0 {
                    log.Errorf("[%s] Parse LenEncInt Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetUint(fv)
                offset = next

            case "LenEncString":
                v, dlen := parseLenEncUint(packet.payload[offset:])
                next = offset + dlen
                if v != 0 && dlen == 0 {
                    log.Errorf("[%s] Parse LenEncString Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                offset = next

                // slen here is not a struct value but following string length
                next = offset + int(v)
                fv := parseFixString(packet.payload[offset:next])
                f.SetString(fv)
                offset = next

            case "StringNUL":
                fv, dlen := parseStringNUL(packet.payload[offset:])
                next = offset + dlen
                if dlen < 0 {
                    log.Errorf("[%s] Parse StringNUL Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetString(fv)
                offset = next

            case "StringEOF":
                fv, dlen := parseStringEOF(packet.payload[offset:])
                next = offset + dlen
                f.SetString(fv)
                offset = next

            case "LenEncStringEOF":
                fv, dlen := parseLenEncStringEOF(packet.payload[offset:])
                next = offset + dlen
                if dlen == 0 {
                    log.Errorf("[%s] Parse LenEncStringEOF Error: [%s]", conn, packet.payload[offset:])
                    return nil, ErrParseData
                }
                f.SetString(fv)
                offset = next

            default:
                log.Errorf("[%s] Unkown Struct Tag Data Type: [%s]", conn, dt)
                return nil, ErrParseData
        }

    }

    if uint64(offset) != packet.length {
        log.Errorf("[%s]: %s Packet Length: [%d], Data Type Length: [%d]", conn, ErrParseData, packet.length, offset)
        return nil, ErrParseData
    }

    return rv.Interface().(Payload), nil

}


func (conn *MyConn) UnpackStream(r io.Reader) {
    for !StopCapture {
        header := make([]byte, 4)
        n, err := io.ReadFull(r, header)
        if err != nil {
            log.Errorf("[%s] %s, [%s], Read bytes [%d]", conn, ErrReadStream, err, n)
            if err == io.EOF {
                close(conn.packets)
                return
            }
            continue
        }

        length_bytes, sequence_byte := header[:3], header[3:]
        payload_length := parseFixUint(length_bytes)
        sequence_id := parseFixUint(sequence_byte)

        payload := make([]byte, payload_length)
        n, err = io.ReadFull(r, payload)
        if err != nil {
            if err == io.EOF {
                close(conn.packets)
                return
            }
            log.Errorf("[%s] %s, [%s], Read bytes [%d]", conn, ErrReadStream, err, n)
            continue
        }

        packet := &Packet{
            length:     payload_length,
            sequence:   sequence_id,
            payload:    payload,
        }

        conn.packets <- packet
    }
    close(conn.packets)
}


func (conn *MyConn) ParseClient(r io.Reader) {

    go conn.UnpackStream(r)

    for !StopCapture {
        packet, ok := <-conn.packets
        if !ok {
            log.Errorf("[%s] Connection Closed", conn)
            return
        }
        log.Infof("[%s] Client packet: [%v]", conn, packet)

        payload, err := conn.ParseCliPayload(packet)
        if err != nil {
            log.Errorf("[%s] Parse Client Payload Error [%s]", conn, err)
            continue
        }
        DisplayClient <- [2]Payload{conn, payload}
    }

}

func (conn *MyConn) ParseServer(r io.Reader) {

    go conn.UnpackStream(r)

    for !StopCapture {
        packet, ok := <-conn.packets
        if !ok {
            log.Errorf("[%s] Connection Closed", conn)
            return
        }
        log.Infof("[%s] Server packet: [%v]", conn, packet)

        payload, err := conn.ParseSrvPayload(packet)
        if err != nil {
            log.Errorf("[%s] Parse Server Payload Error [%s]", conn, err)
            continue
        }
        DisplayServer <- [2]Payload{conn, payload}
    }

}

