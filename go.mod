module github.com/fumiama/ReiBot

go 1.18

require (
	github.com/FloatTech/floatbox v0.0.0-20250513111443-adba80e84e80
	github.com/FloatTech/zbpctrl v1.7.0
	github.com/RomiChan/syncx v0.0.0-20240418144900-b7402ffdebc7
	github.com/fumiama/go-base16384 v1.7.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	github.com/wdvxdr1123/ZeroBot v1.8.0
)

require (
	github.com/FloatTech/sqlite v1.7.0 // indirect
	github.com/FloatTech/ttl v0.0.0-20240716161252-965925764562 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-registry v0.2.7 // indirect
	github.com/fumiama/go-simple-protobuf v0.2.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/fumiama/terasu v0.0.0-20241027183601-987ab91031ce // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20230315142452-642cacee5cc0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.61.0 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.33.1 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
