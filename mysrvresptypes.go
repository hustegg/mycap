package main

import (
    "fmt"
)

// [00] || [fe]
type RespOK struct {
    RespType        []byte  `datatype:"FixBytes" length:"1"`
    AffectedRows    uint64  `datatype:"LenEncUint"`
    LastInsertID    uint64  `datatype:"LenEncUint"`
    StatusFlag      uint64  `datatype:"FixInt" length:"2"`
    Warnings        uint64  `datatype:"FixInt" length:"2"`
    Info            string  `datatype:"StringEOF"`
}

func (r *RespOK) String() string {
    return fmt.Sprintf("Server: OK: AffectedRows: [%d], Warnings: [%d], Info: [%s]", r.AffectedRows, r.LastInsertID, r.Info)
}

// [ff]
type RespErr struct {
    RespType        []byte  `datatype:"FixBytes" length:"1"`
    ErrCode         uint64  `datatype:"FixInt" length:"2"`
    SQLStateMarker  string  `datatype:"FixString" length:"1"`
    SQLState        string  `datatype:"FixString" length:"5"`
    ErrMsg          string  `datatype:"StringEOF"`
}

func (r *RespErr) String() string {
    return fmt.Sprintf("Server: Error: ErrCode: [%d], ErrMsg: [%s], SQLState: [%s]", r.ErrCode, r.ErrMsg, r.SQLState)
}


// [fe]
type RespEOF struct {
    RespType        []byte  `datatype:"FixBytes" length:"1"`
    Warnings        uint64  `datatype:"FixInt" length:"2"`
    StatusFlags     uint64  `datatype:"FixInt" length:"2"`
}

func (r *RespEOF) String() string {
    return fmt.Sprintf("Server: EOF: [deprecated]")
}


// [0a] && fillers
type RespHandShake struct {
    RespType        []byte  `datatype:"FixBytes" length:"1"`
    ServerVersion   string  `datatype:"StringNUL"`
    ConnectionID    uint64  `datatype:"FixInt" length:"4"`
    AuthDataP1      []byte  `datatype:"FixBytes" length:"8"`
    FillerP1        []byte  `datatype:"FixBytes" length:"1"`
    CapFlagsP1      []byte  `datatype:"FixBytes" length:"2"`
    CharSet         uint64  `datatype:"FixInt" length:"1"`
    StatusFlags     []byte  `datatype:"FixBytes" length:"2"`
    CapFlagsP2      []byte  `datatype:"FixBytes" length:"2"`
    AuthDataLen     []byte  `datatype:"FixBytes" length:"1"`
    FillerP2        []byte  `datatype:"FixBytes" length:"10"`
    AuthDataP2      []byte  `datatype:"FixBytes" length:"13"`
    AuthPluginName  string  `datatype:"StringNUL"`
}

func (r *RespHandShake) String() string {
    scramble := append(r.AuthDataP1, r.AuthDataP2...)
    charset := CharSetMap[r.CharSet]
    return fmt.Sprintf("Server: HandShake: Version: [%s], ConnectionID: [%d], Scramble: [%v], Charset: [%s], AuthPlugin: [%s]",
                            r.ServerVersion, r.ConnectionID, scramble, charset, r.AuthPluginName)
}

type RespColCnt struct {
    ColCnt          uint64  `datatype:"LenEncUint"`
}

func (r *RespColCnt) String() string {
    return fmt.Sprintf("Server: Column Count: [%d]", r.ColCnt)
}

type RespColDef struct {
    Catalog         string  `datatype:"LenEncString"`
    Schema          string  `datatype:"LenEncString"`
    Table           string  `datatype:"LenEncString"`
    OrgTable        string  `datatype:"LenEncString"`
    Name            string  `datatype:"LenEncString"`
    OrgName         string  `datatype:"LenEncString"`
    NextLen         []byte  `datatype:"FixBytes" length:"1"`
    CharSet         uint64  `datatype:"FixInt" length:"2"`
    ColLen          uint64  `datatype:"FixInt" length:"4"`
    ColType         uint64  `datatype:"FixInt" length:"1"`
    Flags           []byte  `datatype:"FixBytes" length:"2"`
    DecimalShown    []byte  `datatype:"FixBytes" length:"1"`
    FillerP1        []byte  `datatype:"FixBytes" length:"2"`
    Ignore          string  `datatype:"StringEOF"`
}

func (r *RespColDef) String() string {

    charset := CharSetMap[r.CharSet]
    return fmt.Sprintf("Server: Column Definition: Catalog: [%s], Schema: [%s], Table: [%s], RealTable:[%s], Column: [%s], RealColumn: [%s], CharSet: [%s]",
                        r.Catalog, r.Schema, r.Table, r.OrgTable, r.Name, r.OrgName, charset)
}

type RespStatistics struct {
    Info            string  `datatype:"StringEOF"`
}

func (r *RespStatistics) String() string {
    return fmt.Sprintf("Server: Statistics: [%s]", r.Info)
}



type RespResultSet struct {
    Results          string  `datatype:"LenEncStringEOF"`
}

func (r *RespResultSet) String() string {
    return fmt.Sprintf("Server: Result Set: [%s]", r.Results)
}


