package main

const (
    COM_SLEEP byte = iota
    COM_QUIT
    COM_INIT_DB
    COM_QUERY
    COM_FIELD_LIST
    COM_CREATE_DB
    COM_DROP_DB
    COM_REFRESH
    COM_SHUTDOWN
    COM_STATISTICS
    COM_PROCESS_INFO
    COM_CONNECT
    COM_PROCESS_KILL
    COM_DEBUG
    COM_PING
    COM_TIME
    COM_DELAYED_INSERT
    COM_CHANGE_USER
    COM_BINLOG_DUMP
    COM_TABLE_DUMP
    COM_CONNECT_OUT
    COM_REGISTER_SLAVE
    COM_STMT_PREPARE
    COM_STMT_EXECUTE
    COM_STMT_SEND_LONG_DATA
    COM_STMT_CLOSE
    COM_STMT_RESET
    COM_SET_OPTION
    COM_STMT_FETCH
    COM_DAEMON
    COM_BINLOG_DUMP_GTID
    COM_RESET_CONNECTION
)

const (
    MYSQL_TYPE_DECIMAL byte = iota
    MYSQL_TYPE_TINY
    MYSQL_TYPE_SHORT
    MYSQL_TYPE_LONG
    MYSQL_TYPE_FLOAT
    MYSQL_TYPE_DOUBLE
    MYSQL_TYPE_NULL
    MYSQL_TYPE_TIMESTAMP
    MYSQL_TYPE_LONGLONG
    MYSQL_TYPE_INT24
    MYSQL_TYPE_DATE
    MYSQL_TYPE_TIME
    MYSQL_TYPE_DATETIME
    MYSQL_TYPE_YEAR
    MYSQL_TYPE_NEWDATE
    MYSQL_TYPE_VARCHAR
    MYSQL_TYPE_BIT
    MYSQL_TYPE_TIMESTAMP2
    MYSQL_TYPE_DATETIME2
    MYSQL_TYPE_TIME2
)

const (
    MYSQL_TYPE_JSON byte = iota + 0xf5
    MYSQL_TYPE_NEWDECIMAL
    MYSQL_TYPE_ENUM
    MYSQL_TYPE_SET
    MYSQL_TYPE_TINY_BLOB
    MYSQL_TYPE_MEDIUM_BLOB
    MYSQL_TYPE_LONG_BLOB
    MYSQL_TYPE_BLOB
    MYSQL_TYPE_VAR_STRING
    MYSQL_TYPE_STRING
    MYSQL_TYPE_GEOMETRY
)

var FillerOneByte   = []byte{0}
var FillerTenBytes  = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var Filler23Bytes   = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var ColDefHeader    = []byte{3, 0x64, 0x65, 0x66}

var CharSetMap = map[uint64][2]string {
    1: {"big5", "big5_chinese_ci"},
    3: {"dec8", "dec8_swedish_ci"},
    4: {"cp850", "cp850_general_ci"},
    6: {"hp8", "hp8_english_ci"},
    7: {"koi8r", "koi8r_general_ci"},
    8: {"latin1", "latin1_swedish_ci"},
    9: {"latin2", "latin2_general_ci"},
    10: {"swe7", "swe7_swedish_ci"},
    11: {"ascii", "ascii_general_ci"},
    12: {"ujis", "ujis_japanese_ci"},
    13: {"sjis", "sjis_japanese_ci"},
    16: {"hebrew", "hebrew_general_ci"},
    18: {"tis620", "tis620_thai_ci"},
    19: {"euckr", "euckr_korean_ci"},
    22: {"koi8u", "koi8u_general_ci"},
    24: {"gb2312", "gb2312_chinese_ci"},
    25: {"greek", "greek_general_ci"},
    26: {"cp1250", "cp1250_general_ci"},
    28: {"gbk", "gbk_chinese_ci"},
    30: {"latin5", "latin5_turkish_ci"},
    32: {"armscii8", "armscii8_general_ci"},
    33: {"utf8", "utf8_general_ci"},
    35: {"ucs2", "ucs2_general_ci"},
    36: {"cp866", "cp866_general_ci"},
    37: {"keybcs2", "keybcs2_general_ci"},
    38: {"macce", "macce_general_ci"},
    39: {"macroman","        macroman_general_ci"},
    40: {"cp852", "cp852_general_ci"},
    41: {"latin7", "latin7_general_ci"},
    45: {"utf8mb4", "utf8mb4_general_ci"},
    51: {"cp1251", "cp1251_general_ci"},
    54: {"utf16", "utf16_general_ci"},
    56: {"utf16le", "utf16le_general_ci"},
    57: {"cp1256", "cp1256_general_ci"},
    59: {"cp1257", "cp1257_general_ci"},
    60: {"utf32", "utf32_general_ci"},
    63: {"binary", "binary"},
    92: {"geostd8", "geostd8_general_ci"},
    95: {"cp932", "cp932_japanese_ci"},
    97: {"eucjpms", "eucjpms_japanese_ci"},
    248: {"gb18030", "gb18030_chinese_ci"},
}

var RefreshCommandMap = map[uint64]string {
    0x01: "FLUSH PRIVILEGES",
    0x02: "FLUSH LOGS",
    0x04: "FLUSH TABLES",
    0x08: "FLUSH HOSTS",
    0x10: "FLUSH STATUS",
    0x20: "FLUSH THREADS",
    0x40: "RESET SLAVE",
    0x80: "RESET MASTER",
    }

var ShutdownTypeMap = map[uint64]string {
    0x00: "SHUTDOWN_WAIT_ALL_BUFFERS",
    0x01: "SHUTDOWN_WAIT_CONNECTIONS",
    0x02: "SHUTDOWN_WAIT_TRANSACTIONS",
    0x08: "SHUTDOWN_WAIT_UPDATES",
    0x10: "SHUTDOWN_WAIT_ALL_BUFFERS",
    0x11: "SHUTDOWN_WAIT_NON_INNODB_BUFFERS",
    0xfe: "KILL_QUERY",
    0xff: "KILL_CONNECTION",
    }

var CliComTypes = map[byte]func() Payload {
    COM_SLEEP:                  func() Payload{return &ComSleep{}},
    COM_QUIT:                   func() Payload{return &ComQuit{}},
    COM_INIT_DB:                func() Payload{return &ComInitDB{}},
    COM_QUERY:                  func() Payload{return &ComQuery{}},
    COM_FIELD_LIST:             func() Payload{return &ComFieldList{}},
    COM_CREATE_DB:              func() Payload{return &ComCreateDB{}},
    COM_DROP_DB:                func() Payload{return &ComDropDB{}},
    COM_REFRESH:                func() Payload{return &ComRefresh{}},
    COM_SHUTDOWN:               func() Payload{return &ComShutdown{}},
    COM_STATISTICS:             func() Payload{return &ComStatistics{}},
    COM_PROCESS_INFO:           func() Payload{return &ComProcessInfo{}},
    COM_CONNECT:                func() Payload{return &ComConnect{}},
    COM_PROCESS_KILL:           func() Payload{return &ComProcessKill{}},
    COM_DEBUG:                  func() Payload{return &ComDebug{}},
    COM_PING:                   func() Payload{return &ComPing{}},
    COM_TIME:                   func() Payload{return &ComTime{}},
    COM_DELAYED_INSERT:         func() Payload{return &ComDelayedInsert{}},
    COM_CHANGE_USER:            func() Payload{return &ComChangeUser{}},
    //COM_BINLOG_DUMP:            func() Payload{return &ComBinlogDump{}},
    //COM_TABLE_DUMP:             func() Payload{return &ComTableDump{}},
    //COM_CONNECT_OUT:            func() Payload{return &ComConnectOut{}},
    //COM_REGISTER_SLAVE:         func() Payload{return &ComRegisterSlave{}},
    //COM_STMT_PREPARE:           func() Payload{return &ComStmtPrepare{}},
    //COM_STMT_EXECUTE:           func() Payload{return &ComStmtExecute{}},
    //COM_STMT_SEND_LONG_DATA:    func() Payload{return &ComStmtSendLongData{}},
    //COM_STMT_CLOSE:             func() Payload{return &ComStmtClose{}},
    //COM_STMT_RESET:             func() Payload{return &ComStmtReset{}},
    //COM_SET_OPTION:             func() Payload{return &ComSetOption{}},
    //COM_STMT_FETCH:             func() Payload{return &ComStmtFetch{}},
    COM_DAEMON:                 func() Payload{return &ComDaemon{}},
    //COM_BINLOG_DUMP_GTID:       func() Payload{return &ComBinlogDumpGTID{}},
    //COM_RESET_CONNECTION:       func() Payload{return &ComResetConnection{}},
}


