package main

import (
    "errors"
    "fmt"
    "io"
    "os"
    "strconv"
    //"sync"
    "time"

    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/tcpassembly"
    "github.com/google/gopacket/tcpassembly/tcpreader"
)

// TODO: support prepare statements needs prepare id
//var PrepareStatements sync.Map

var StopCapture bool
var DisplayClient = make(chan [2]Payload, 8192)
var DisplayServer = make(chan [2]Payload, 8192)

type StreamFactory struct {}

func (sf *StreamFactory) New(net, tcp gopacket.Flow) tcpassembly.Stream {
    mstrm := &MyStream{
        net: net,
        tcp: tcp,
        rdr: tcpreader.NewReaderStream(),
    }

    go mstrm.run()

    return &mstrm.rdr
}

type MyStream struct {
    net, tcp    gopacket.Flow
    rdr         tcpreader.ReaderStream
}

func (mstrm *MyStream) run() {

    src_ip, dst_ip := mstrm.net.Endpoints()
    src_port, dst_port := mstrm.tcp.Endpoints()
    conn := MyConn{
        src_ip:     src_ip.String(),
        dst_ip:     dst_ip.String(),
        src_port:   src_port.String(),
        dst_port:   dst_port.String(),
        client:     dst_port.String() == strconv.Itoa(*arg_myport),
        packets:    make(chan *Packet, 4096),
    }

    if conn.client {
        conn.ParseClient(&mstrm.rdr)
    } else {
        conn.ParseServer(&mstrm.rdr)
    }
}

func Capture() error {


    var handle *pcap.Handle
    var err error

    if *arg_capfile != "" {
        fmt.Fprintf(os.Stderr, "Read MySQL packets from pcap file\n")
        handle, err = pcap.OpenOffline(*arg_capfile)
    } else {
        fmt.Fprintf(os.Stderr, "Start capture MySQL packets, device:%s, max-cap-num:%d, packet-filter:%s\n", *arg_dev, *arg_pcnt, BPF)
        handle, err = pcap.OpenLive(*arg_dev, int32(*arg_snaplen), *arg_promisc, pcap.BlockForever)
    }
    if err != nil {
        log.Fatal(err)
    }

    if err := handle.SetBPFFilter(BPF); err != nil {
        log.Fatal(err)
    }

    stream_factory := StreamFactory{}
    stream_pool := tcpassembly.NewStreamPool(&stream_factory)
    assembler := tcpassembly.NewAssembler(stream_pool)

    packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
    raw_packets := packetSource.Packets()
    ticker := time.Tick(time.Minute)

    var pcap_cnt int64
    for pcap_cnt <= *arg_pcnt {
        select {
        case raw_packet := <-raw_packets:
            pcap_cnt++
            if raw_packet == nil {
                return io.EOF
            }

            log.Tracef("Captured %d packets", pcap_cnt)
            log.Trace(raw_packet)

            if raw_packet.NetworkLayer() == nil ||
                raw_packet.TransportLayer() == nil ||
                raw_packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
                log.Errorf("Packet layer error: [%v]", raw_packet)
            }

            tcp := raw_packet.TransportLayer().(*layers.TCP)
            assembler.AssembleWithTimestamp(raw_packet.NetworkLayer().NetworkFlow(), tcp, raw_packet.Metadata().Timestamp)

        case <-ticker:
            assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
            //assembler.FlushOlderThan(time.Now())
        }
    }

    StopCapture = true
    close(raw_packets)
    return errors.New("Capture End")
}

func DisplayPayload(Display chan [2]Payload) {
    for !StopCapture {
        payloads, ok := <-Display
        if !ok {
            log.Warn("Display End")
            return
        }
        // conn, payload
        display.Infof("[%s] %s", payloads[0], payloads[1])
    }

}

