package main

import (
	"sort"
	"strings"
)

type commalist []string

func (cl *commalist) Set(value string) error {
	*cl = strings.Split(value, ",")
	return nil
}

func (cl commalist) String() string { return strings.Join(cl, ",") }

type commaset map[string]struct{}

func (tgl *commaset) Set(value string) error {
	values := strings.Split(value, ",")
	if *tgl == nil {
		*tgl = make(map[string]struct{}, len(values))
	}
	for _, v := range values {
		(*tgl)[v] = struct{}{}
	}
	return nil
}

func (tgl commaset) String() string {
	arr := make([]string, 0, len(tgl))
	for k := range tgl {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return strings.Join(arr, ",")
}
