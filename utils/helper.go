package utils

import (
	"fmt"
	"strconv"
)

func SliceStringToUInt(data []string) []uint {
	ids := []uint{}

	for _, d := range data {
		id, err := strconv.Atoi(d)

		if err != nil {
			fmt.Println(err.Error())
		}

		ids = append(ids, uint(id))
	}
	return ids
}
