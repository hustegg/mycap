package main

import (
    "fmt"
    "strings"
)

type ComSleep struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComSleep) String() string {
    return fmt.Sprintf("Client: Sleep: [Internal Command]")
}

type ComQuit struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComQuit) String() string {
    return fmt.Sprintf("Client: Quit: MySQL Client Quit")
}

type ComInitDB struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    SchemaName  string  `datatype:"StringEOF"`
}

func (c *ComInitDB) String() string {
    return fmt.Sprintf("Client: Init DB: USE %s", c.SchemaName)
}

type ComQuery struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    QueryStr    string  `datatype:"StringEOF"`
}

func (c *ComQuery) String() string {
    return fmt.Sprintf("Client: Query: %s", c.QueryStr)
}

type ComFieldList struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    TableName   string  `datatype:"StringNUL"`
    FieldWc     string  `datatype:"StringEOF"`
}

func (c *ComFieldList) String() string {
    return fmt.Sprintf("Client: Field List: DESC %s %s", c.TableName, c.FieldWc)
}

type ComCreateDB struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    SchemaName  string  `datatype:"StringEOF"`
}

func (c *ComCreateDB) String() string {
    return fmt.Sprintf("Client: CREATE DATABASE: %s", c.SchemaName)
}

type ComDropDB struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    SchemaName  string  `datatype:"StringEOF"`
}

func (c *ComDropDB) String() string {
    return fmt.Sprintf("Client: DROP DATABASE: %s", c.SchemaName)
}

type ComRefresh struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    SubCommand  uint64  `datatype:"FixUint" length:"1"`
}

func (c *ComRefresh) String() string {
    var s []string
    for k, v := range RefreshCommandMap {
        if c.SubCommand & k != 0 {
            s = append(s, v)
        }
    }
    return fmt.Sprintf("Client: Refresh: %s", strings.Join(s, ", "))
}

type ComShutdown struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    //ShutdownType    uint64  `datatype:"FixUint" length:"1"`
}

func (c *ComShutdown) String() string {
    //shutdown_type := ShutdownTypeMap[c.ShutdownType]
    //return fmt.Sprintf("Client: SHUTDOWN WITH %s", shutdown_type)
    return fmt.Sprintf("Client: Shutdown")
}

type ComStatistics struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComStatistics) String() string {
    return fmt.Sprintf("Client: Statistics")
}

type ComProcessInfo struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComProcessInfo) String() string {
    return fmt.Sprintf("Client: SHOW PROCESSLIST: [deprecated]")
}

type ComConnect struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComConnect) String() string {
    return fmt.Sprintf("Client: Connect: [Internal Command]")
}

type ComProcessKill struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    ConnectionID    uint64  `datatype:"FixInt" length:"4"`
}

func (c *ComProcessKill) String() string {
    return fmt.Sprintf("Client: KILL %d [deprecated]", c.ConnectionID)
}

type ComDebug struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComDebug) String() string {
    return fmt.Sprintf("Client: Debug")
}

type ComPing struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComPing) String() string {
    return fmt.Sprintf("Client: Ping")
}

type ComTime struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComTime) String() string {
    return fmt.Sprintf("Client: Time: [Internal Command]")
}

type ComDelayedInsert struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComDelayedInsert) String() string {
    return fmt.Sprintf("Client: Delayed Insert: [Internal Command]")
}

type ComChangeUser struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    User        string  `datatype:"StringNUL"`
    AuthResp    string  `datatype:"LenEncString"`
    SchemaName  string  `datatype:"StringNUL"`
    Info            string  `datatype:"StringNULEOF"`
}

func (c *ComChangeUser) String() string {
    return fmt.Sprintf("Client: Change User: User: [%s], SchemaName: [%s], AuthResp: [%s], Info: [%s]",
                                                c.User, c.SchemaName, c.AuthResp, c.Info)
}

type ComResetConn struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComResetConn) String() string {
    return fmt.Sprintf("Client: Reset Connection")
}

type ComDaemon struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
}

func (c *ComDaemon) String() string {
    return fmt.Sprintf("Client: Daemon [Internal Command]")
}

type ComBinlogDump struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    BinlogPos   uint64  `datatype:"FixUint" length:"4"`
    Flags       uint64  `datatype:"FixUint" length:"2"`
    ServerID    uint64  `datatype:"FixUint" length:"4"`
    BinlogName  string  `datatype:"StringEOF"`
}

func (c *ComBinlogDump) String() string {
    return fmt.Sprintf("Client: Binlog Dump: SlaveServerID: [%d], MasterLogFile: [%s], MasterLogPos: [%d]",
        c.ServerID, c.BinlogName, c.BinlogPos)
}

type ComStmtPrepare struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    Query       string  `datatype:"StringEOF"`
}

func (c *ComStmtPrepare) String() string {
    return fmt.Sprintf("Client: Statement Prepare: [%s]", c.Query)
}

type ComStmtExecute struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    StmtID      uint64  `datatype:"FixUint" length:"4"`
    Flags       uint64  `datatype:"FixUint" length:"1"`
    IterCnt     uint64  `datatype:"FixUint" length:"4"`
    Info        string  `datatype:"StringEOF"`
}

func (c *ComStmtExecute) String() string {
    return fmt.Sprintf("Client: Statement Execute: StatementID: [%d], Info: [%s]", c.StmtID, c.Info)
}

type ComStmtClose struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    StmtID      uint64  `datatype:"FixUint" length:"4"`
}

func (c *ComStmtClose) String() string {
    return fmt.Sprintf("Client: Statement Close: [%d]", c.StmtID)
}

type ComStmtReset struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    StmtID      uint64  `datatype:"FixUint" length:"4"`
}

func (c *ComStmtReset) String() string {
    return fmt.Sprintf("Client: Statement Reset: [%d]", c.StmtID)
}



type ComBinlogDumpGTID struct {
    ComType     []byte  `datatype:"FixBytes" length:"1"`
    Flags       uint64  `datatype:"FixUint" length:"2"`
    ServerID    uint64  `datatype:"FixUint" length:"4"`
    BinlogName  string  `datatype:"FixLenEncString" length:"4"`
    BinlogPos   uint64  `datatype:"FixUint" length:"8"`
    Info        string  `datatype:"StringEOF"`

}

func (c *ComBinlogDumpGTID) String() string {
    return fmt.Sprintf("Client: Binlog Dump: SlaveServerID: [%d], MasterLogFile: [%s], MasterLogPos: [%d], Info: [%s]",
        c.ServerID, c.BinlogName, c.BinlogPos, c.Info)
}






type ComHandShake struct {
    CapFlags        []byte  `datatype:"FixBytes" length:"4"`
    MaxPacketSize   uint64  `datatype:"FixInt" length:"4"`
    CharSet         uint64  `datatype:"FixInt" length:"1"`
    FillerP1        []byte  `datatype:"FixBytes" length:"23"`
    UserName        string  `datatype:"StringNUL"`
    AuthData        []byte  `datatype:"LenEncBytes"`
    Info            string  `datatype:"StringNULEOF"`
}

func (c *ComHandShake) String() string {
    charset := CharSetMap[c.CharSet]
    return fmt.Sprintf("Client: HandShake: UserName:[%s], AuthData: [%v], CharSet: [%v], Info: [%s]",
                        c.UserName, c.AuthData, charset, c.Info)
}

