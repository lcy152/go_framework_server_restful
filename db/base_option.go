package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type option struct {
	Search    map[OptionKey]interface{}
	Regex     map[OptionKey]bool
	TimeKey   OptionKey
	TimeStart int64
	TimeEnd   int64
	PageSize  int64
	PageIndex int64
	Ascend    bool
	Sort      []OptionKey
}

func NewOptions() *option {
	return &option{
		Search:    map[OptionKey]interface{}{},
		Regex:     map[OptionKey]bool{},
		TimeKey:   OptInvalid,
		TimeStart: 0,
		TimeEnd:   0,
		PageSize:  0,
		PageIndex: 0,
		Ascend:    true,
		Sort:      []OptionKey{},
	}
}

func (opt *option) toQueryBsonM(need map[OptionKey]string, aggregate bool) []bson.M {
	var query []bson.M

	var nonregex []bson.M
	var regex []bson.M
	for key, val := range opt.Search {
		if needKey, ok1 := need[key]; ok1 {
			setNonregex := func(value interface{}) {
				if _, ok2 := opt.Regex[key]; ok2 {
					if !aggregate {
						regex = append(regex, bson.M{needKey: bson.M{"$regex": value, "$options": "$i"}})
					}
				} else {
					if aggregate {
						nonregex = append(nonregex, bson.M{"$match": bson.M{needKey: value}})
					} else {
						nonregex = append(nonregex, bson.M{needKey: value})
					}
				}
			}
			switch val.(type) {
			case int:
				value := val.(int)
				setNonregex(value)
			case int8:
				value := val.(int8)
				setNonregex(value)
			case int16:
				value := val.(int16)
				setNonregex(value)
			case int32:
				value := val.(int32)
				setNonregex(value)
			case int64:
				value := val.(int64)
				setNonregex(value)
			case bool:
				value := val.(bool)
				setNonregex(value)
			case primitive.DateTime:
				value := val.(primitive.DateTime)
				setNonregex(value)
			case string:
				value := val.(string)
				if _, ok2 := opt.Regex[key]; ok2 {
					value = `\Q` + value + `\E`
				}
				setNonregex(value)
			case []string:
				value := val.([]string)
				if aggregate {
					nonregex = append(nonregex, bson.M{"$match": bson.M{needKey: bson.M{"$in": value}}})
				} else {
					nonregex = append(nonregex, bson.M{needKey: bson.M{"$in": value}})
				}
			}
		}
	}
	if len(nonregex) != 0 {
		query = append(query, nonregex...)
	}
	if len(regex) != 0 && !aggregate {
		query = append(query, bson.M{"$or": regex})
	}

	if opt.TimeKey != OptInvalid && opt.TimeEnd != 0 {
		if needKey, ok := need[opt.TimeKey]; ok {
			if aggregate {
				query = append(query, bson.M{"$match": bson.M{needKey: bson.D{{"$gte", primitive.DateTime(opt.TimeStart)}, {"$lte", primitive.DateTime(opt.TimeEnd)}}}})
			} else {
				query = append(query, bson.M{needKey: bson.D{{"$gte", primitive.DateTime(opt.TimeStart)}, {"$lte", primitive.DateTime(opt.TimeEnd)}}})
			}
		}
	}

	return query
}

func (opt *option) toPagingBsonM(need map[OptionKey]string) []bson.M {
	var paging []bson.M

	var sort []bson.M
	ascend := 1
	if !opt.Ascend {
		ascend = -1
	}
	for _, val := range opt.Sort {
		if needKey, ok := need[val]; ok {
			sort = append(sort, bson.M{"$sort": bson.M{needKey: ascend}})
		}
	}
	if len(sort) != 0 {
		paging = append(paging, sort...)
	}

	if opt.PageSize != 0 {
		limit := int64(opt.PageSize)
		skip := int64(opt.PageSize * opt.PageIndex)

		paging = append(paging, bson.M{"$skip": skip})
		paging = append(paging, bson.M{"$limit": limit})
	}

	return paging
}

func (opt *option) toPagingFindOptions(need map[OptionKey]string) options.FindOptions {
	findOption := options.FindOptions{}

	var sort bson.D
	ascend := 1
	if !opt.Ascend {
		ascend = -1
	}
	for _, val := range opt.Sort {
		if needKey, ok := need[val]; ok {
			sort = append(sort, bson.E{needKey, ascend})
		}
	}
	if len(sort) != 0 {
		findOption.Sort = sort
	}

	if opt.PageSize != 0 {
		limit := int64(opt.PageSize)
		skip := int64(opt.PageSize * opt.PageIndex)

		findOption.Skip = &skip
		findOption.Limit = &limit
	}

	return findOption
}

func (opt *option) toAggregate(need map[OptionKey]string) []bson.M {
	query := opt.toQueryBsonM(need, true)

	paging := opt.toPagingBsonM(need)
	query = append(query, paging...)

	return query
}

func (opt *option) toFind(need map[OptionKey]string) (bson.M, options.FindOptions) {
	query := opt.toQueryBsonM(need, false)

	findOptions := opt.toPagingFindOptions(need)
	q := bson.M{}
	if query != nil {
		q = bson.M{"$and": query}
	}
	return q, findOptions
}
