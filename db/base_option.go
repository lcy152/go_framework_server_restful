package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SortOption struct {
	Key    OptionKey
	Ascend bool
}

type option struct {
	Match     map[OptionKey]interface{}
	EQ        map[OptionKey]interface{}
	NEQ       map[OptionKey]interface{}
	GTE       map[OptionKey]interface{}
	GT        map[OptionKey]interface{}
	LTE       map[OptionKey]interface{}
	LT        map[OptionKey]interface{}
	IN        map[OptionKey]interface{}
	NIN       map[OptionKey]interface{}
	OR        []map[OptionKey]interface{}
	ElemMatch map[OptionKey]map[OptionKey]interface{}
	TimeKey   OptionKey
	TimeStart int64
	TimeEnd   int64
	PageSize  int
	PageIndex int
	Sort      []SortOption
}

func NewOptions() *option {
	return &option{
		Match:     map[OptionKey]interface{}{},
		EQ:        map[OptionKey]interface{}{},
		NEQ:       map[OptionKey]interface{}{},
		GTE:       map[OptionKey]interface{}{},
		GT:        map[OptionKey]interface{}{},
		LTE:       map[OptionKey]interface{}{},
		LT:        map[OptionKey]interface{}{},
		IN:        map[OptionKey]interface{}{},
		NIN:       map[OptionKey]interface{}{},
		OR:        []map[OptionKey]interface{}{},
		ElemMatch: map[OptionKey]map[OptionKey]interface{}{},
		TimeKey:   OptInvalid,
		TimeStart: 0,
		TimeEnd:   0,
		PageSize:  0,
		PageIndex: 0,
		Sort:      []SortOption{},
	}
}

func (opt *option) toQueryBsonM(need map[OptionKey]string, aggregate bool) []bson.M {
	var query []bson.M

	var eqQuery []bson.M
	var matchQuery []bson.M
	var orQuery []bson.M
	for key, val := range opt.EQ {
		if needKey, ok1 := need[key]; ok1 {
			switch val.(type) {
			case []string:
				value := val.([]string)
				if aggregate {
					eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$in": value}}})
				} else {
					eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$in": value}})
				}
			default:
				if needKey, ok1 := need[key]; ok1 {
					if aggregate {
						eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: val}})
					} else {
						eqQuery = append(eqQuery, bson.M{needKey: val})
					}
				}
			}
		}
	}
	for key, value := range opt.NEQ {
		if needKey, ok1 := need[key]; ok1 {
			if aggregate {
				eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$ne": value}}})
			} else {
				eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$ne": value}})
			}
		}
	}
	for key, value := range opt.LTE {
		if needKey, ok1 := need[key]; ok1 {
			if aggregate {
				eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$lte": value}}})
			} else {
				eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$lte": value}})
			}
		}
	}
	for key, value := range opt.LT {
		if needKey, ok1 := need[key]; ok1 {
			if aggregate {
				eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$lt": value}}})
			} else {
				eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$lt": value}})
			}
		}
	}
	for key, value := range opt.GTE {
		if needKey, ok1 := need[key]; ok1 {
			if aggregate {
				eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$gte": value}}})
			} else {
				eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$gte": value}})
			}
		}
	}
	for key, value := range opt.GT {
		if needKey, ok1 := need[key]; ok1 {
			if aggregate {
				eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$gt": value}}})
			} else {
				eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$gt": value}})
			}
		}
	}
	for key, val := range opt.IN {
		if needKey, ok1 := need[key]; ok1 {
			switch val.(type) {
			case []string:
				value := val.([]string)
				if aggregate {
					eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$in": value}}})
				} else {
					eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$in": value}})
				}
			}
		}
	}
	for key, val := range opt.NIN {
		if needKey, ok1 := need[key]; ok1 {
			switch val.(type) {
			case []string:
				value := val.([]string)
				if aggregate {
					eqQuery = append(eqQuery, bson.M{"$match": bson.M{needKey: bson.M{"$nin": value}}})
				} else {
					eqQuery = append(eqQuery, bson.M{needKey: bson.M{"$nin": value}})
				}
			}
		}
	}
	for key, val := range opt.Match {
		if needKey, ok1 := need[key]; ok1 {
			seteqQuery := func(value interface{}) {
				if aggregate {
				} else {
					matchQuery = append(matchQuery, bson.M{needKey: bson.M{"$regex": value, "$options": "$i"}})
				}
			}
			switch val.(type) {
			case string:
				value := val.(string)
				value = `\Q` + value + `\E`
				seteqQuery(value)
			}
		}
	}
	for _, valList := range opt.OR {
		var andRegex []bson.M
		for key, val := range valList {
			if needKey, ok1 := need[key]; ok1 {
				switch val.(type) {
				case []string:
					value := val.([]string)
					if aggregate {
						andRegex = append(andRegex, bson.M{"$match": bson.M{needKey: bson.M{"$in": value}}})
					} else {
						andRegex = append(andRegex, bson.M{needKey: bson.M{"$in": value}})
					}
				default:
					if needKey, ok1 := need[key]; ok1 {
						if aggregate {
							andRegex = append(andRegex, bson.M{"$match": bson.M{needKey: val}})
						} else {
							andRegex = append(andRegex, bson.M{needKey: val})
						}
					}
				}
			}
		}
		orQuery = append(orQuery, bson.M{"$and": andRegex})
	}
	var elemMatch []bson.M
	for key1, valMap := range opt.ElemMatch {
		if needKey1, ok1 := need[key1]; ok1 {
			var tempQuery []bson.M
			for key, value := range valMap {
				if needKey2, ok2 := need[key]; ok2 {
					if aggregate {
						tempQuery = append(tempQuery, bson.M{"$match": bson.M{needKey2: value}})
					} else {
						tempQuery = append(tempQuery, bson.M{needKey2: value})
					}
				}
			}
			elemMatch = append(elemMatch, bson.M{needKey1: bson.M{"$elemMatch": bson.M{"$and": tempQuery}}})
		}
	}
	if len(eqQuery) != 0 {
		query = append(query, eqQuery...)
	}
	if len(orQuery) != 0 {
		query = append(query, bson.M{"$or": orQuery})
	}
	if len(matchQuery) != 0 {
		query = append(query, bson.M{"$or": matchQuery})
	}
	if len(elemMatch) != 0 {
		query = append(query, elemMatch...)
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
	for _, val := range opt.Sort {
		if needKey, ok := need[val.Key]; ok {
			ascend := 1
			if !val.Ascend {
				ascend = -1
			}
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
	for _, val := range opt.Sort {
		if needKey, ok := need[val.Key]; ok {
			ascend := 1
			if !val.Ascend {
				ascend = -1
			}
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
