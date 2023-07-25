// Backen > ResponseGenerator > CacheTimeTables
// This file provides a cache generation mechanism for fully static bootstrap nodes.

package responsegenerator

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/logging"
	"errors"
	// "time"
)

// type CacheTimeTable struct {
// 	ResponseUrl string
// 	StartsFrom  int64
// 	EndsAt      int64
// }

// type CacheTimeTable api.ResultCache

var cacheTimeBlocks = []api.Timestamp{6000000, 600000, 60000, 6000, 600}

//                            100000m  10000m  1000m  100m  10m
//                            1660h    166h    16.6h  1.66h
//                            69d      6.9d

/*
NewCacheTimeTable generates caches table based on the given beginning and the end. If multiple intervals are provided, it needs to start from the largest interval and go in order, and all the intervals need to be divisible with each other for the minimal disturbances to existing caches.
*/
func NewCacheTimeTable(beginning api.Timestamp, end api.Timestamp, intervals []api.Timestamp) []api.ResultCache {
	var timetable []api.ResultCache

	// return timetable
	/*
	   For each time block, scan through the range and insert the biggest block possible first, then go to smaller blocks and insert those.
	*/
	for _, tb := range intervals {
		for beginning+tb <= end {
			timetable = append(timetable, api.ResultCache{
				StartsFrom: beginning,
				EndsAt:     beginning + tb,
			})
			beginning = beginning + tb
		}
	}
	/*
	   In the end, if there's a gap between the end and provided end, add the difference to the last cache.
	*/
	if len(timetable) > 0 {
		timetable[len(timetable)-1].EndsAt = end
	}
	return timetable
}

func MakeConsolidatedTimeTable(tt *[]api.ResultCache, timeBlocks []api.Timestamp) []api.ResultCache {
	var ctt []api.ResultCache

	for _, timeBlock := range timeBlocks {
		for {
			cti, err := MakeConsolidatedCacheItem(tt, timeBlock)
			if err != nil {
				break
			}
			ctt = append(ctt, cti)
		}
	}
	return ctt
}

func MakeConsolidatedCacheItem(tt *[]api.ResultCache, timeBlock api.Timestamp) (api.ResultCache, error) {
	if len(*tt) == 0 {
		return api.ResultCache{}, errors.New("The source cache time table is empty.")
	}
	if delta := (*tt)[0].EndsAt - (*tt)[0].StartsFrom; delta >= timeBlock {
		cci := (*tt)[0]
		newtt := DeleteFromCacheTimeTable(0, *tt)
		*tt = newtt
		return cci, nil // All is fine, not touching the cache. Returning as is and removing the entity from the base.
	}
	// Our cache is smaller than the time block we want to build. Let's try to add to it.
	cci := api.ResultCache{
		StartsFrom: (*tt)[0].StartsFrom,
		EndsAt:     (*tt)[0].EndsAt,
	}
	for counter := 1; counter < len(*tt); counter++ {
		if delta := cci.EndsAt - cci.StartsFrom; delta < timeBlock {
			cci.EndsAt = (*tt)[counter].EndsAt
			continue
		}
		// Successfully created a consolidated cache. Delete the caches we've used from the base.
		for i := 0; i < counter; i++ {
			newtt := DeleteFromCacheTimeTable(0, *tt)
			// ^ Always deleting from zero because we're deleting from the beginning, which makes the item that would be [1] to [0]
			*tt = newtt
		}
		return cci, nil
	}
	// We couldn't complete the cache build to a size larger than the given time block. Bail.
	return api.ResultCache{}, errors.New("This time block size is too large.")
}

// DeleteFromCacheTimeTable deletes a timetable item from a timetable
func DeleteFromCacheTimeTable(i int, tt []api.ResultCache) []api.ResultCache {
	if i > len(tt)-1 {
		logging.Logf(1, "DeleteFromCacheTimeTable was requested to delete a cachetimetable item that was out of bounds. Requested deletion index: %v, TimeTable: %v\n", i, tt)
	}
	return append(tt[0:i], tt[i+1:len(tt)]...)
}
