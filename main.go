package main

import (
	"github.com/zc310/utils"
	"math/rand"
	"net/http"
	"strconv"

	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Code struct {
		// ID 单号
		ID utils.WideString `json:"id"` // string
		// Count 注数
		Count int `json:"count"` // 0
		// Lot  玩法
		Lot utils.WideString `json:"lot"` // 竞彩足球
		// Note 备注
		Note utils.WideString `json:"note"` // string
		// Time 方案时间
		Time time.Time `json:"time"` // 2019-04-25T00:21:25.324Z
		// Issue 期号
		Issue string `json:"issue"` // 2019001
		// Code 投注号码
		Code string `json:"-"`
	}
)

var (
	codes = map[string]*Code{}
)

const KEY = "special-key"

//----------
// Handlers
//----------

// getCode 根据时间查询号码列表
func getCode(c echo.Context) error {
	var list []*Code
	for _, v := range codes {
		list = append(list, v)
	}
	return c.JSON(http.StatusOK, list)
}

// getCodeByID 根据单号查询号码内容
func getCodeByID(c echo.Context) error {
	co, ok := codes[c.Param("id")]
	if !ok {
		return c.String(http.StatusNotFound, "")
	}
	return c.String(http.StatusOK, co.Code)
}

func GetApiKey(c echo.Context) string {
	key := c.Request().Header.Get("api_key")
	if len(key) == 0 {
		key = c.QueryParam("api_key")
	}
	return key
}

// apiKeyCheck 检查api_key 是否符合要求

func apiKeyCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if GetApiKey(c) == KEY {
			return next(c)
		}
		return c.NoContent(http.StatusUnauthorized)
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/api/swagger.yaml", "api/swagger.yaml")

	g := e.Group("/code/v1", apiKeyCheck)
	// Routes
	g.GET("/code", getCode)
	g.GET("/code/:id", getCodeByID)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func init() {
	newCode(99)
}
func newCode(n int) {
	var id int
	var cid string
	for i := 0; i < n; i++ {
		cid = bson.NewObjectId().Hex()
		id = rand.Intn(len(LotName))
		co := &Code{ID: utils.WideString(cid),
			Count: rand.Intn(100),
			Lot:   LotName[id],
			Issue: strconv.Itoa(2019001 + rand.Intn(99)),
			Note:  utils.WideString("号码说明"),
			Time:  time.Now(),
			Code:  LotCode[id],
		}

		codes[cid] = co
	}

}

var (
	LotName = []utils.WideString{"竞彩足球", "竞彩篮球", "单场", "胜负彩", "任选九", "半全场", "进球彩", "排列3", "排列5", "大乐透", "七星彩"}
	LotCode = []string{JcZq, JcLq, Dc, SFC, R9, BQC, JQC, PL3, PL5, DLT, QXC}

	JcZq = `SPF|171125017=0/1/3,171124025=0,171126014=3|3*1|99
SPF|171120003=3|1*1|99
RSP|171120003=3|1*1|99
BQC|171121002=33/30/11/03/00,171121003=31/13/10/01|2*1|91
CBF|171121002=10/21/31/40/42/51/90/11/33/01/12/13/04/24/15/09,171121003=20/30/32/41/50/52/00/22/99/02/03/23/14/05/25|2*1|99
JQS|171121004=0/2/4/6,171121005=1/3/5/7|2*1|99
ZHH|171122002-RSP=3/1/0,171123001-BQC=33/31,171123002-SPF=3/1/0,171123003-BQC=33/31/30/13/11/10/03/01/00|4*1|11`

	JcLq = `SF|171125017=0/3,171124025=0,171126014=3|3*1|99
RSF|171120003=3|1*1|99
SFC|171121002=03/01/11/04/05,171121003=11/12/01/02|2*1|91
DXF|171121002=1/0,171121003=0/1|2*1|99
LHH|171122002-SF=3/0,171123001-SFC=03/11,171123002-RSF=3/1/0|3*1|11`

	Dc = `SPF|49=3/0,50=3/0,51=3/0,52=3/0,53=3/0,54=3/0,55=3/0,56=3/0,57=3/0|9*1|1
CBF|20=10/20/21/30/00/11/22/01/02/12/03,21=10/20/21/30/00/11/22/01/02/12/03,22=10/20/21/30/00/11/22/01/02/12/03|3*1|1
JQS|44=0/1/2/3/4/5,45=0/1/2/3/4/5,46=0/1/2/3/4/5|3*1|1
BQC|53=33/31/30/13/11/10/00,54=33/30/13/11/10/00,55=33/30/13/00|3*1|1
SXP|26=0/1/2/3,27=0/1/2/3,28=0/1/2/3,29=0/1/2/3|4*1|1`

	SFC = `31103101133301
33310003110331
03131100313103
33110031310303
11130013001333
0,0,0,0,0,0,0,3,1,0,0,0,0,013`
	R9 = `31#0#101#33##1
##31#00#11033#
031311###13##3
3#110#3#31#3#1
11##00#30#13#0
0,30,1,#,#,#,#,0,0,0,#,3,3,3`
	JQC = `01232101
3,2,1,0,1,2,3,213
(0)(1)(2)(3)(03)(1)(2)(3)`

	BQC = `000000111113
333333000001
000100300111
0,0,0,0,0,3,1,0,0,0,0,013`

	DLT = `01 02 03 04 05 06 17 + 01 02
01 02 03 13 23 + 01 12
02 06 13 18 20 + 04 11
03 09 23 24 26 + 10 12
04 10 13 14 16 + 08 11
05 22 23 24 26 + 07 09`

	QXC = `0123456
1234567
2345678
3456789
4567898
56,0,1,9,9,9,9`

	PL3 = `012
123
234
345
12,9,9`
	PL5 = `012
123
234
345
12,9,9`
)
