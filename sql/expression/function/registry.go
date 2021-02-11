// Copyright 2020-2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package function

import (
	"math"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression/function/aggregation"
	"github.com/dolthub/go-mysql-server/sql/expression/function/aggregation/window"
)

// Defaults is the function map with all the default functions.
var Defaults = []sql.Function{
	// elt, find_in_set, insert, load_file, locate
	sql.Function1{Name: "abs", Fn: NewAbsVal},
	sql.Function1{Name: "acos", Fn: NewAcos},
	sql.Function1{Name: "array_length", Fn: NewArrayLength},
	sql.Function1{Name: "ascii", Fn: NewAscii},
	sql.Function1{Name: "asin", Fn: NewAsin},
	sql.Function1{Name: "atan", Fn: NewAtan},
	sql.Function1{Name: "avg", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewAvg(e) }},
	sql.Function1{Name: "bin", Fn: NewBin},
	sql.Function1{Name: "bit_length", Fn: NewBitlength},
	sql.Function1{Name: "ceil", Fn: NewCeil},
	sql.Function1{Name: "ceiling", Fn: NewCeil},
	sql.Function1{Name: "char_length", Fn: NewCharLength},
	sql.Function1{Name: "character_length", Fn: NewCharLength},
	sql.FunctionN{Name: "coalesce", Fn: NewCoalesce},
	sql.FunctionN{Name: "concat", Fn: NewConcat},
	sql.FunctionN{Name: "concat_ws", Fn: NewConcatWithSeparator},
	sql.NewFunction0("connection_id", NewConnectionID),
	sql.Function1{Name: "cos", Fn: NewCos},
	sql.Function1{Name: "cot", Fn: NewCot},
	sql.Function1{Name: "count", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewCount(e) }},
	sql.Function1{Name: "crc32", Fn: NewCrc32},
	sql.NewFunction0("curdate", NewCurrDate),
	sql.NewFunction0("current_date", NewCurrentDate),
	sql.NewFunction0("current_time", NewCurrentTime),
	sql.NewFunction0("current_timestamp", NewCurrTimestamp),
	sql.NewFunction0("current_user", NewCurrentUser),
	sql.NewFunction0("curtime", NewCurrTime),
	sql.Function1{Name: "date", Fn: NewDate},
	sql.FunctionN{Name: "date_add", Fn: NewDateAdd},
	sql.Function2{Name: "date_format", Fn: NewDateFormat},
	sql.FunctionN{Name: "date_sub", Fn: NewDateSub},
	sql.FunctionN{Name: "datetime", Fn: NewDatetime},
	sql.Function1{Name: "day", Fn: NewDay},
	sql.Function1{Name: "dayname", Fn: NewDayName},
	sql.Function1{Name: "dayofmonth", Fn: NewDay},
	sql.Function1{Name: "dayofweek", Fn: NewDayOfWeek},
	sql.Function1{Name: "dayofyear", Fn: NewDayOfYear},
	sql.Function1{Name: "degrees", Fn: NewDegrees},
	sql.Function1{Name: "explode", Fn: NewExplode},
	sql.Function1{Name: "first", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewFirst(e) }},
	sql.Function1{Name: "floor", Fn: NewFloor},
	sql.Function1{Name: "from_base64", Fn: NewFromBase64},
	sql.FunctionN{Name: "greatest", Fn: NewGreatest},
	sql.Function1{Name: "hex", Fn: NewHex},
	sql.Function1{Name: "hour", Fn: NewHour},
	sql.Function3{Name: "if", Fn: NewIf},
	sql.Function2{Name: "ifnull", Fn: NewIfNull},
	sql.Function2{Name: "instr", Fn: NewInstr},
	sql.Function1{Name: "is_binary", Fn: NewIsBinary},
	sql.FunctionN{Name: "json_extract", Fn: NewJSONExtract},
	sql.Function1{Name: "json_unquote", Fn: NewJSONUnquote},
	sql.Function1{Name: "last", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewLast(e) }},
	sql.Function1{Name: "lcase", Fn: NewLower},
	sql.FunctionN{Name: "least", Fn: NewLeast},
	sql.Function2{Name: "left", Fn: NewLeft},
	sql.Function1{Name: "length", Fn: NewLength},
	sql.Function1{Name: "ln", Fn: NewLogBaseFunc(float64(math.E))},
	sql.FunctionN{Name: "log", Fn: NewLog},
	sql.Function1{Name: "log10", Fn: NewLogBaseFunc(float64(10))},
	sql.Function1{Name: "log2", Fn: NewLogBaseFunc(float64(2))},
	sql.Function1{Name: "lower", Fn: NewLower},
	sql.FunctionN{Name: "lpad", Fn: NewPadFunc(lPadType)},
	sql.Function1{Name: "ltrim", Fn: NewTrimFunc(lTrimType)},
	sql.Function1{Name: "max", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewMax(e) }},
	sql.Function1{Name: "md5", Fn: NewMD5},
	sql.Function1{Name: "microsecond", Fn: NewMicrosecond},
	sql.FunctionN{Name: "mid", Fn: NewSubstring},
	sql.Function1{Name: "min", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewMin(e) }},
	sql.Function1{Name: "minute", Fn: NewMinute},
	sql.Function1{Name: "month", Fn: NewMonth},
	sql.Function1{Name: "monthname", Fn: NewMonthName},
	sql.FunctionN{Name: "now", Fn: NewNow},
	sql.Function2{Name: "nullif", Fn: NewNullIf},
	sql.Function2{Name: "pow", Fn: NewPower},
	sql.Function2{Name: "power", Fn: NewPower},
	sql.Function1{Name: "radians", Fn: NewRadians},
	sql.FunctionN{Name: "rand", Fn: NewRand},
	sql.FunctionN{Name: "regexp_matches", Fn: NewRegexpMatches},
	sql.Function2{Name: "repeat", Fn: NewRepeat},
	sql.Function3{Name: "replace", Fn: NewReplace},
	sql.Function1{Name: "reverse", Fn: NewReverse},
	sql.FunctionN{Name: "round", Fn: NewRound},
	sql.Function0{Name: "row_number", Fn: window.NewRowNumber},
	sql.FunctionN{Name: "rpad", Fn: NewPadFunc(rPadType)},
	sql.Function1{Name: "rtrim", Fn: NewTrimFunc(rTrimType)},
	sql.Function1{Name: "second", Fn: NewSecond},
	sql.Function1{Name: "sha", Fn: NewSHA1},
	sql.Function1{Name: "sha1", Fn: NewSHA1},
	sql.Function2{Name: "sha2", Fn: NewSHA2},
	sql.Function1{Name: "sign", Fn: NewSign},
	sql.Function1{Name: "sin", Fn: NewSin},
	sql.Function1{Name: "sleep", Fn: NewSleep},
	sql.Function1{Name: "soundex", Fn: NewSoundex},
	sql.Function2{Name: "split", Fn: NewSplit},
	sql.Function1{Name: "sqrt", Fn: NewSqrt},
	sql.FunctionN{Name: "substr", Fn: NewSubstring},
	sql.FunctionN{Name: "substring", Fn: NewSubstring},
	sql.Function3{Name: "substring_index", Fn: NewSubstringIndex},
	sql.Function1{Name: "sum", Fn: func(e sql.Expression) sql.Expression { return aggregation.NewSum(e) }},
	sql.Function1{Name: "tan", Fn: NewTan},
	sql.Function1{Name: "time_to_sec", Fn: NewTimeToSec},
	sql.FunctionN{Name: "timestamp", Fn: NewTimestamp},
	sql.Function1{Name: "to_base64", Fn: NewToBase64},
	sql.Function1{Name: "trim", Fn: NewTrimFunc(bTrimType)},
	sql.Function1{Name: "ucase", Fn: NewUpper},
	sql.Function1{Name: "unhex", Fn: NewUnhex},
	sql.FunctionN{Name: "unix_timestamp", Fn: NewUnixTimestamp},
	sql.FunctionN{Name: "utc_timestamp", Fn: NewUTCTimestamp},
	sql.Function2{Name: "timediff", Fn: NewTimeDiff},
	sql.Function1{Name: "upper", Fn: NewUpper},
	sql.NewFunction0("user", NewUser),
	sql.FunctionN{Name: "week", Fn: NewWeek},
	sql.Function1{Name: "values", Fn: NewValues},
	sql.Function1{Name: "weekday", Fn: NewWeekday},
	sql.Function1{Name: "weekofyear", Fn: NewWeekOfYear},
	sql.Function1{Name: "year", Fn: NewYear},
	sql.FunctionN{Name: "yearweek", Fn: NewYearWeek},
}

func GetLockingFuncs(ls *sql.LockSubsystem) []sql.Function {
	return []sql.Function{
		sql.Function2{Name: "get_lock", Fn: CreateNewGetLock(ls)},
		sql.Function1{Name: "is_free_lock", Fn: NewIsFreeLock(ls)},
		sql.Function1{Name: "is_used_lock", Fn: NewIsUsedLock(ls)},
		sql.NewFunction0("release_all_locks", NewReleaseAllLocks(ls)),
		sql.Function1{Name: "release_lock", Fn: NewReleaseLock(ls)},
	}
}
