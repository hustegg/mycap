# MyCap
获取MySQL网络包并解析内容，支持指定IP黑白名单、最大抓包数量，后台运行时输出JSON格式日志便于分析


## MySQL version recommend
建议MySQL-5.7.5以上，较低版本的包可能导致解析异常，未严格测试低版本协议


## Packets type supported
暂不支持压缩协议、Prepare语句、复制协议，后续版本会支持


## Make & Install
>git clone https://github.com/hustegg/mycap  
>cd mycap  
>go build  


## Usage

>Usage of ./mycap:  
>  -b value  
>        Packets white ip list separated by comma  
>  -c int  
>        Packets number captured before exit (default 1024)  
>  -d string  
>        Capture MySQL Packet direction [client|server|both] (default "client")  
>  -f string  
>        Captured packets filename  
>  -i string  
>        Network interface name (default "eth0")  
>  -j    Logging with JSON formatter  
>  -m    Capture with promisc mode  
>  -p int  
>        MySQL server port capture (default 3306)  
>  -s int  
>        Snap length for pcap packet capture (default 1600)  
>  -v    Logging in detail  
>  -vv  
>        Logging in verbose  
>  -w value  
>        Packets white ip list separated by comma  

# Example
sudo ./mycap -i eth1 -w 192.168.0.1 -d both
>Start capture MySQL packets, device:eth1, max-cap-num:1024, packet-filter:tcp and (port 3306) and (host 192.168.0.1)  
>WARN[2019-07-15 12:42:01.423059] [192.168.0.1:3306 => 192.168.0.2:35029] Server: HandShake: Version: [5.7.17-log], ConnectionID: [2053427], Scramble: [[74 27 83 56 27 83 1 106 17 91 26 98 62 77 106 32 3 24 83 1 0]], Charset: [[latin1 latin1_swedish_ci]], AuthPlugin: [mysql_native_password]  
>WARN[2019-07-15 12:42:01.423146] [192.168.0.2:35029 => 192.168.0.1:3306] Client: HandShake: UserName:[test], AuthData: [[71 235 187 29 224 25 104 20 143 3 168 29 69 30 158 93 219 21 217 145]], CharSet: [[latin1 latin1_swedish_ci]], Info: [mysql_native_password]  
>WARN[2019-07-15 12:42:01.455108] [192.168.0.1:3306 => 192.168.0.2:35029] Server: OK: AffectedRows: [0], Warnings: [0], Info: []  
>WARN[2019-07-15 12:42:01.455115] [192.168.0.2:35029 => 192.168.0.1:3306] Client: Query: select @@version_comment limit 1  
>WARN[2019-07-15 12:42:01.488202] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Count: [1]  
>WARN[2019-07-15 12:42:01.488234] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Definition: Catalog: [def], Schema: [], Table: [], Column:  [@@version_comment], CharSet: [[ ]]  
>WARN[2019-07-15 12:42:01.488243] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Result Set: [Source distribution]  
>WARN[2019-07-15 12:42:01.488252] [192.168.0.1:3306 => 192.168.0.2:35029] Server: OK: AffectedRows: [0], Warnings: [0], Info: [_ character_set_client latin1 character_set_connection latin1 character_set_results latin1]  
>WARN[2019-07-15 12:42:10.385813] [192.168.0.2:35029 => 192.168.0.1:3306] Client: Query: select now()  
>WARN[2019-07-15 12:42:10.419304] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Count: [1]  
>WARN[2019-07-15 12:42:10.419347] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Definition: Catalog: [def], Schema: [], Table: [], Column: [now()], CharSet: [[ ]]  
>WARN[2019-07-15 12:42:10.419356] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Result Set: [2019-07-15 12:42:10]  
>WARN[2019-07-15 12:42:10.419367] [192.168.0.1:3306 => 192.168.0.2:35029] Server: OK: AffectedRows: [0], Warnings: [0], Info: []  
>WARN[2019-07-15 12:42:18.449785] [192.168.0.2:35029 => 192.168.0.1:3306] Client: Query: select * from dba.delay_monitor  
>WARN[2019-07-15 12:42:18.482213] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Count: [3]  
>WARN[2019-07-15 12:42:18.482229] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Definition: Catalog: [def], Schema: [dba], Table: [delay_monitor], Column: [id], CharSet: [[ ]]  
>WARN[2019-07-15 12:42:18.482238] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Definition: Catalog: [def], Schema: [dba], Table: [delay_monitor], Column: [Ftime], CharSet: [[ ]]  
>WARN[2019-07-15 12:42:18.482245] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Column Definition: Catalog: [def], Schema: [dba], Table: [delay_monitor], Column: [Fgtid], CharSet: [[ ]]  
>WARN[2019-07-15 12:42:18.482252] [192.168.0.1:3306 => 192.168.0.2:35029] Server: Result Set: [1, 2019-06-12 15:14:17, TestHisDB_QP_20190612151417_b0b166d7-8ce1-11e9-a0a8-c81fbecfd710]  
>WARN[2019-07-15 12:42:18.482259] [192.168.0.1:3306 => 192.168.0.2:35029] Server: OK: AffectedRows: [0], Warnings: [0], Info: []  
>ERRO[2019-07-15 12:42:20.137657] [192.168.0.2:35029 => 192.168.0.1:3306] Read Stream Error, [EOF], Read bytes [0]  
>ERRO[2019-07-15 12:42:20.137733] [192.168.0.2:35029 => 192.168.0.1:3306] Connection Closed  
>WARN[2019-07-15 12:42:20.13767] [192.168.0.2:35029 => 192.168.0.1:3306] Client: Quit: MySQL Client Quit  
>ERRO[2019-07-15 12:42:20.1706] [192.168.0.1:3306 => 192.168.0.2:35029] Read Stream Error, [EOF], Read bytes [0]  
>ERRO[2019-07-15 12:42:20.170619] [192.168.0.1:3306 => 192.168.0.2:35029] Connection Closed  
