module txtresearch

go 1.17

require github.com/syndtr/goleveldb v1.0.0

require golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect

require (
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/kardianos/service v1.2.2
	github.com/thinkeridea/go-extend v1.3.2
)
//以下是引用本地topsdk包
require "topsdk" v0.0.0
replace "topsdk" => "../topsdk"

