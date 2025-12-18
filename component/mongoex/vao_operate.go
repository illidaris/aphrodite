package mongoex

// mongo 运算符

//go:generate stringer -type=Operate -output vao_operate_string.go
type Operate int32

const (
	OperateNil Operate = iota + 1
)

// 比较操作符
const (
	OperateEq  Operate = 100100 + iota // $eq
	OperateNe                          // $ne
	OperateGt                          // $gt
	OperateGte                         // $gte
	OperateLt                          // $lt
	OperateLte                         // $lte
	OperateIn                          // $in
	OperateNin                         // $nin
)

// 逻辑操作符
const (
	OperateAnd Operate = 100200 + iota // $and
	OperateOr                          // $or
	OperateNot                         // $not
	OperateNor                         // $nor
)

// 元素操作符
const (
	OperateExists Operate = 100300 + iota // &$exists
	OperateType                           // $type
)

// 数组操作符
const (
	OperateAll       Operate = 100400 + iota // $all
	OperateElemMatch                         // $elemMatch
	OperateSize                              // $size
)

// 其他操作符
const (
	OperateRegex     Operate = 100500 + iota // $regex
	OperateText                              // $text
	OperateWhere                             // $where
	OperateNear                              // $near
	OperateGeoWithin                         // $geoWithin
)

// 聚合操作符
const (
	OperateSum      Operate = 200100 + iota // $sum
	OperateAvg                              // $avg
	OperateMin                              // $min
	OperateMax                              // $max
	OperatePush                             // $push
	OperateAddToSet                         // $addToSet
	OperateFirst                            // $first
	OperateLast                             // $last
)

// 常用聚合管道操作符
const (
	OperateProject Operate = 300100 + iota // $project
	OperateMatch                           // $match
	OperateLimit                           // $limit
	OperateSkip                            // $skip
	OperateUnwind                          // $unwind
	OperateGroup                           // $group
	OperateSort                            // $sort
	OperateGeoNear                         // $geoNear
)

var OperateCodeMap = map[Operate]string{
	OperateEq:        "$eq",
	OperateNe:        "$ne",
	OperateGt:        "$gt",
	OperateGte:       "$gte",
	OperateLt:        "$lt",
	OperateLte:       "$lte",
	OperateIn:        "$in",
	OperateNin:       "$nin",
	OperateAnd:       "$and",
	OperateOr:        "$or",
	OperateNot:       "$not",
	OperateNor:       "$nor",
	OperateExists:    "&$exists",
	OperateType:      "$type",
	OperateAll:       "$all",
	OperateElemMatch: "$elemMatch",
	OperateSize:      "$size",
	OperateRegex:     "$regex",
	OperateText:      "$text",
	OperateWhere:     "$where",
	OperateNear:      "$near",
	OperateGeoWithin: "$geoWithin",
	OperateSum:       "$sum",
	OperateAvg:       "$avg",
	OperateMin:       "$min",
	OperateMax:       "$max",
	OperatePush:      "$push",
	OperateAddToSet:  "$addToSet",
	OperateFirst:     "$first",
	OperateLast:      "$last",
	OperateProject:   "$project",
	OperateMatch:     "$match",
	OperateLimit:     "$limit",
	OperateSkip:      "$skip",
	OperateUnwind:    "$unwind",
	OperateGroup:     "$group",
	OperateSort:      "$sort",
	OperateGeoNear:   "$geoNear",
}
