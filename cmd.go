package redis

import (
	"bytes"
	"strconv"
)

type Cmd []byte

// redis commands listed as int redis.io/commands#<category>, 201604
var (
	cmdAuth Cmd
	CmdEcho Cmd
	CmdPing Cmd
	cmdQuit Cmd
	CmdSelect Cmd

	CmdGeoadd Cmd
	CmdGeohash Cmd
	CmdGeopos Cmd
	CmdGeodist Cmd
	CmdGeoradius Cmd
	CmdGeoradiusbymember Cmd

	CmdHdel Cmd
	CmdHexists Cmd
	CmdHget Cmd
	CmdHgetall Cmd
	CmdHincrby Cmd
	CmdHincrbyfloat Cmd
	CmdHKeys Cmd
	CmdHlen Cmd
	CmdHmget Cmd
	CmdHmset Cmd
	CmdHset Cmd
	CmdHsetnx Cmd
	CmdHstrlen Cmd
	CmdHvals Cmd
	CmdHscan Cmd

	CmdPfadd Cmd
	CmdPfcount Cmd
	CmdPfmerge Cmd

	CmdDel Cmd
	CmdDump Cmd
	CmdExists Cmd
	CmdExpire Cmd
	CmdExpireat Cmd
	CmdKeys Cmd
	CmdMigrate Cmd
	CmdMove Cmd
	CmdObject Cmd
	CmdPersist Cmd
	CmdPexpire Cmd
	CmdPexpireat Cmd
	CmdPttl Cmd
	CmdRandomkey Cmd
	CmdRename Cmd
	CmdRenamenx Cmd
	CmdRestore Cmd
	CmdSort Cmd
	CmdTtl Cmd
	CmdType Cmd
	CmdWait Cmd
	CmdScan Cmd

	cmdBlpop Cmd
	cmdBrpop Cmd
	cmdBrpoplpush Cmd
	CmdLindex Cmd
	CmdLinsert Cmd
	CmdLlen Cmd
	CmdLpop Cmd
	CmdLpush Cmd
	CmdLpushx Cmd
	CmdLrange Cmd
	CmdLrem Cmd
	CmdLset Cmd
	CmdLtrim Cmd
	CmdRpop Cmd
	CmdRpoplpush Cmd
	CmdRpush Cmd
	CmdRpushx Cmd

	cmdPsubscribe Cmd
	CmdPubsub Cmd
	CmdPublish Cmd
	cmdPunsubscribe Cmd
	cmdSubscribe Cmd
	cmdUnsubscribe Cmd

	CmdEval Cmd
	CmdEvalsha Cmd
	CmdScript Cmd

	CmdBgrewriteaof Cmd
	CmdBgsave Cmd
	CmdClient Cmd
	CmdCommand Cmd
	CmdConfig Cmd
	CmdDbsize Cmd
	CmdDebug Cmd
	CmdFlushAll Cmd
	CmdFlushdb Cmd
	CmdInfo Cmd
	CmdLastsave Cmd
	cmdMonitor Cmd
	CmdRole Cmd
	CmdSave Cmd
	CmdShutdown Cmd
	CmdSlaveof Cmd
	CmdSlowlog Cmd
	CmdSync Cmd
	CmdTime Cmd

	CmdSadd Cmd
	CmdScard Cmd
	CmdSdiff Cmd
	CmdSdiffstore Cmd
	CmdSinter Cmd
	CmdSinterstore Cmd
	CmdSismember Cmd
	CmdSmembers Cmd
	CmdSmove Cmd
	CmdSpop Cmd
	CmdSrandmember Cmd
	CmdSrem Cmd
	CmdSunion Cmd
	CmdSunionstore Cmd
	CmdSscan Cmd

	CmdZadd Cmd
	CmdZcard Cmd
	CmdZcount Cmd
	CmdZincrby Cmd
	CmdZinterstore Cmd
	CmdZlexcount Cmd
	CmdZrange Cmd
	CmdZrangebylex Cmd
	CmdZrevrangebylex Cmd
	CmdZrangebyscore Cmd
	CmdZrank Cmd
	CmdZrem Cmd
	CmdZremrangebylex Cmd
	CmdZremrangebyrank Cmd
	CmdZremrangebyscore Cmd
	CmdZrevrange Cmd
	CmdZrevrangebyscore Cmd
	CmdZrevrank Cmd
	CmdZscore Cmd
	CmdZunionstore Cmd
	CmdZscan Cmd

	CmdAppend Cmd
	CmdBitcount Cmd
	CmdBitop Cmd
	CmdBitpos Cmd
	CmdDecr Cmd
	CmdDecrby Cmd
	CmdGet Cmd
	CmdGetbit Cmd
	CmdGetrange Cmd
	CmdGetset Cmd
	CmdIncr Cmd
	CmdIncrby Cmd
	CmdIncybyfloat Cmd
	CmdMget Cmd
	CmdMset Cmd
	CmdMsetnx Cmd
	CmdPsetex Cmd
	CmdSet Cmd
	CmdSetbit Cmd
	CmdSetex Cmd
	CmdSetnx Cmd
	CmdSetrange Cmd
	CmdStrlen Cmd

	cmdDiscard Cmd
	CmdExec Cmd
	CmdMulti Cmd
	cmdUnwatch Cmd
	cmdWatch Cmd
)

var cmdQuitBytes []byte



//////////////////////////////////////////////////////////////////////

func init() {
	makeCmdBytes := func(s string) []byte {
		var buf bytes.Buffer
		b := []byte(s)
		ll := len(b)
		lenBytes := []byte(strconv.Itoa(ll))
		buf.WriteByte(MarkBulkString)
		buf.Write(lenBytes)
		buf.Write(bytesCrlf)
		buf.Write(b)
		buf.Write(bytesCrlf)
		return buf.Bytes()
	}

	cmdAuth = Cmd(makeCmdBytes("AUTH"))
	CmdEcho = Cmd(makeCmdBytes("ECHO"))
	CmdPing = Cmd(makeCmdBytes("PING"))
	cmdQuit = Cmd(makeCmdBytes("QUIT"))
	CmdSelect = Cmd(makeCmdBytes("SELECT"))

	CmdGeoadd = Cmd(makeCmdBytes("GEOADD"))
	CmdGeohash = Cmd(makeCmdBytes("GEOHAS"))
	CmdGeopos = Cmd(makeCmdBytes("GEOPOS"))
	CmdGeodist = Cmd(makeCmdBytes("GEODIST"))
	CmdGeoradius = Cmd(makeCmdBytes("GEORADIUS"))
	CmdGeoradiusbymember = Cmd(makeCmdBytes("GEORADIUSBYMEMBER"))

	CmdHdel = Cmd(makeCmdBytes("HDEL"))
	CmdHexists = Cmd(makeCmdBytes("HEXISTS"))
	CmdHget = Cmd(makeCmdBytes("HGET"))
	CmdHgetall = Cmd(makeCmdBytes("HGETALL"))
	CmdHincrby = Cmd(makeCmdBytes("HINCRBY"))
	CmdHincrbyfloat = Cmd(makeCmdBytes("HINCRBYFLOAT"))
	CmdHKeys = Cmd(makeCmdBytes("HKEYS"))
	CmdHlen = Cmd(makeCmdBytes("HLEN"))
	CmdHmget = Cmd(makeCmdBytes("HMGET"))
	CmdHmset = Cmd(makeCmdBytes("HMSET"))
	CmdHset = Cmd(makeCmdBytes("HSET"))
	CmdHsetnx = Cmd(makeCmdBytes("HSETNX"))
	CmdHstrlen = Cmd(makeCmdBytes("HSTRLEN"))
	CmdHvals = Cmd(makeCmdBytes("HVALS"))
	CmdHscan = Cmd(makeCmdBytes("HSCAN"))

	CmdPfadd = Cmd(makeCmdBytes("PFADD"))
	CmdPfcount = Cmd(makeCmdBytes("PFCOUNT"))
	CmdPfmerge = Cmd(makeCmdBytes("PFMERGE"))

	CmdDel = Cmd(makeCmdBytes("DEL"))
	CmdDump = Cmd(makeCmdBytes("DUMP"))
	CmdExists = Cmd(makeCmdBytes("EXISTS"))
	CmdExpire = Cmd(makeCmdBytes("EXPIRE"))
	CmdExpireat = Cmd(makeCmdBytes("EXPIREAT"))
	CmdKeys = Cmd(makeCmdBytes("KEYS"))
	CmdMigrate = Cmd(makeCmdBytes("MIGRATE"))
	CmdMove = Cmd(makeCmdBytes("MOVE"))
	CmdObject = Cmd(makeCmdBytes("OBJECT"))
	CmdPersist = Cmd(makeCmdBytes("PERSIST"))
	CmdPexpire = Cmd(makeCmdBytes("PEXPIRE"))
	CmdPexpireat = Cmd(makeCmdBytes("PEXPIREAT"))
	CmdPttl = Cmd(makeCmdBytes("PTTL"))
	CmdRandomkey = Cmd(makeCmdBytes("RANDOMKEY"))
	CmdRename = Cmd(makeCmdBytes("RENAME"))
	CmdRenamenx = Cmd(makeCmdBytes("RENAMENX"))
	CmdRestore = Cmd(makeCmdBytes("RESTORE"))
	CmdSort = Cmd(makeCmdBytes("SORT"))
	CmdTtl = Cmd(makeCmdBytes("TTL"))
	CmdType = Cmd(makeCmdBytes("TYPE"))
	CmdWait = Cmd(makeCmdBytes("WAIT"))
	CmdScan = Cmd(makeCmdBytes("SCAN"))

	cmdBlpop = Cmd(makeCmdBytes("BLPOP"))
	cmdBrpop = Cmd(makeCmdBytes("BRPOP"))
	cmdBrpoplpush = Cmd(makeCmdBytes("BRPOPLPUSH"))
	CmdLindex = Cmd(makeCmdBytes("LINDEX"))
	CmdLinsert = Cmd(makeCmdBytes("LINSERT"))
	CmdLlen = Cmd(makeCmdBytes("LLEN"))
	CmdLpop = Cmd(makeCmdBytes("LPOP"))
	CmdLpush = Cmd(makeCmdBytes("LPUSH"))
	CmdLpushx = Cmd(makeCmdBytes("LPUSHX"))
	CmdLrange = Cmd(makeCmdBytes("LRANGE"))
	CmdLrem = Cmd(makeCmdBytes("LREM"))
	CmdLset = Cmd(makeCmdBytes("LSET"))
	CmdLtrim = Cmd(makeCmdBytes("LTRIM"))
	CmdRpop = Cmd(makeCmdBytes("RPOP"))
	CmdRpoplpush = Cmd(makeCmdBytes("RPOPLPUSH"))
	CmdRpush = Cmd(makeCmdBytes("RPUSH"))
	CmdRpushx = Cmd(makeCmdBytes("RPUSHX"))

	cmdPsubscribe = Cmd(makeCmdBytes("PSUBSCRIBE"))
	CmdPubsub = Cmd(makeCmdBytes("PUBSUB"))
	CmdPublish = Cmd(makeCmdBytes("PUBLISH"))
	cmdPunsubscribe = Cmd(makeCmdBytes("PUNSUBSCRIBE"))
	cmdSubscribe = Cmd(makeCmdBytes("SUBSCRIBE"))
	cmdUnsubscribe = Cmd(makeCmdBytes("UNSUBSCRIBE"))

	CmdEval = Cmd(makeCmdBytes("EVAL"))
	CmdEvalsha = Cmd(makeCmdBytes("EVALSHA"))
	CmdScript = Cmd(makeCmdBytes("SCRIPT"))

	CmdBgrewriteaof = Cmd(makeCmdBytes("BGREWRITEAOF"))
	CmdBgsave = Cmd(makeCmdBytes("BGSAVE"))
	CmdClient = Cmd(makeCmdBytes("CLIENT"))
	CmdCommand = Cmd(makeCmdBytes("COMMAND"))
	CmdConfig = Cmd(makeCmdBytes("CONFIG"))
	CmdDbsize = Cmd(makeCmdBytes("DBSIZE"))
	CmdDebug = Cmd(makeCmdBytes("DEBUG"))
	CmdFlushAll = Cmd(makeCmdBytes("FLUSHALL"))
	CmdFlushdb = Cmd(makeCmdBytes("FLUSHDB"))
	CmdInfo = Cmd(makeCmdBytes("INFO"))
	CmdLastsave = Cmd(makeCmdBytes("LASTSAVE"))
	cmdMonitor = Cmd(makeCmdBytes("MONITOR"))
	CmdRole = Cmd(makeCmdBytes("ROLE"))
	CmdSave = Cmd(makeCmdBytes("SAVE"))
	CmdShutdown = Cmd(makeCmdBytes("SHUTDOWN"))
	CmdSlaveof = Cmd(makeCmdBytes("SLAVEOF"))
	CmdSlowlog = Cmd(makeCmdBytes("SLOWLOG"))
	CmdSync = Cmd(makeCmdBytes("SYNC"))
	CmdTime = Cmd(makeCmdBytes("TIME"))

	CmdSadd = Cmd(makeCmdBytes("SADD"))
	CmdScard = Cmd(makeCmdBytes("SCARD"))
	CmdSdiff = Cmd(makeCmdBytes("SDIFF"))
	CmdSdiffstore = Cmd(makeCmdBytes("SDIFFSTORE"))
	CmdSinter = Cmd(makeCmdBytes("SINTER"))
	CmdSinterstore = Cmd(makeCmdBytes("SINTERSTORE"))
	CmdSismember = Cmd(makeCmdBytes("SISMEMBER"))
	CmdSmembers = Cmd(makeCmdBytes("SMEMBERS"))
	CmdSmove = Cmd(makeCmdBytes("SMOVE"))
	CmdSpop = Cmd(makeCmdBytes("SPOP"))
	CmdSrandmember = Cmd(makeCmdBytes("SRANDMEMBER"))
	CmdSrem = Cmd(makeCmdBytes("SREM"))
	CmdSunion = Cmd(makeCmdBytes("SUNION"))
	CmdSunionstore = Cmd(makeCmdBytes("SUNIONSTORE"))
	CmdSscan = Cmd(makeCmdBytes("SSCAN"))

	CmdZadd = Cmd(makeCmdBytes("ZADD"))
	CmdZcard = Cmd(makeCmdBytes("ZCARD"))
	CmdZcount = Cmd(makeCmdBytes("ZCOUNT"))
	CmdZincrby = Cmd(makeCmdBytes("ZINCRBY"))
	CmdZinterstore = Cmd(makeCmdBytes("ZINTERSTORE"))
	CmdZlexcount = Cmd(makeCmdBytes("ZLEXCOUNT"))
	CmdZrange = Cmd(makeCmdBytes("ZRANGE"))
	CmdZrangebylex = Cmd(makeCmdBytes("ZRANGEBYLEX"))
	CmdZrevrangebylex = Cmd(makeCmdBytes("ZREVRANGEBYLEX"))
	CmdZrangebyscore = Cmd(makeCmdBytes("ZRANGEBYSCORE"))
	CmdZrank = Cmd(makeCmdBytes("ZRANK"))
	CmdZrem = Cmd(makeCmdBytes("ZREM"))
	CmdZremrangebylex = Cmd(makeCmdBytes("ZREMRANGEBYLEX"))
	CmdZremrangebyrank = Cmd(makeCmdBytes("ZREMRANGEBYRANK"))
	CmdZremrangebyscore = Cmd(makeCmdBytes("ZREMRANGEBYSCORE"))
	CmdZrevrange = Cmd(makeCmdBytes("ZREVRANGE"))
	CmdZrevrangebyscore = Cmd(makeCmdBytes("ZREVRANGEBYSCORE"))
	CmdZrevrank = Cmd(makeCmdBytes("ZREVRANK"))
	CmdZscore = Cmd(makeCmdBytes("ZSCORE"))
	CmdZunionstore = Cmd(makeCmdBytes("ZUNIONSTORE"))
	CmdZscan = Cmd(makeCmdBytes("ZSCAN"))

	CmdAppend = Cmd(makeCmdBytes("APPEND"))
	CmdBitcount = Cmd(makeCmdBytes("BITCOUNT"))
	CmdBitop = Cmd(makeCmdBytes("BITOP"))
	CmdBitpos = Cmd(makeCmdBytes("BITPOS"))
	CmdDecr = Cmd(makeCmdBytes("DECR"))
	CmdDecrby = Cmd(makeCmdBytes("DECRBY"))
	CmdGet = Cmd(makeCmdBytes("GET"))
	CmdGetbit = Cmd(makeCmdBytes("GETBIT"))
	CmdGetrange = Cmd(makeCmdBytes("GETRANGE"))
	CmdGetset = Cmd(makeCmdBytes("GETSET"))
	CmdIncr = Cmd(makeCmdBytes("INCR"))
	CmdIncrby = Cmd(makeCmdBytes("INCRBY"))
	CmdIncybyfloat = Cmd(makeCmdBytes("INCRBYFLOAT"))
	CmdMget = Cmd(makeCmdBytes("MGET"))
	CmdMset = Cmd(makeCmdBytes("MSET"))
	CmdMsetnx = Cmd(makeCmdBytes("MSETNX"))
	CmdPsetex = Cmd(makeCmdBytes("PSETEX"))
	CmdSet = Cmd(makeCmdBytes("SET"))
	CmdSetbit = Cmd(makeCmdBytes("SETBIT"))
	CmdSetex = Cmd(makeCmdBytes("SETEX"))
	CmdSetnx = Cmd(makeCmdBytes("SETNX"))
	CmdSetrange = Cmd(makeCmdBytes("SETRANGE"))
	CmdStrlen = Cmd(makeCmdBytes("STRLEN"))

	cmdDiscard = Cmd(makeCmdBytes("DISCARD"))
	CmdExec = Cmd(makeCmdBytes("EXEC"))
	CmdMulti = Cmd(makeCmdBytes("MULTI"))
	cmdUnwatch = Cmd(makeCmdBytes("UNWATCH"))
	cmdWatch = Cmd(makeCmdBytes("WATCH"))
}
