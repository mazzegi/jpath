package jpath

import (
	"fmt"
	"reflect"
	"testing"
)

func assertNoErr(t *testing.T, err error) {
	if err == nil {
		return
	}
	t.Fatalf("expect no-err; got: %v", err)
}

func assertErr(t *testing.T, err error) {
	if err != nil {
		return
	}
	t.Fatalf("expect err; got none")
}

func assertEqual(t *testing.T, want any, have any) {
	if reflect.DeepEqual(want, have) {
		return
	}
	t.Fatalf("want %v; have %v", want, have)
}

type TT1 struct {
	F1 string
	F2 int
	F3 float64
	T2 TT2
}

type TT2 struct {
	G1 string
	G2 int
	G3 float64
	S1 []int
	S2 []TT3
}

type TT3 struct {
	H1   string
	H2   int
	M1   map[string]string
	TT2s []TT2
}

func TestQuery(t *testing.T) {
	v := TT1{
		F1: "a-string",
		F2: 42,
		F3: 0.42,
		T2: TT2{
			G1: "a-g-string",
			G2: 433,
			G3: 12.433,
			S1: []int{2, 3, 4},
			S2: []TT3{
				{
					H1: "h1-s",
					H2: 1,
				},
				{
					H1: "h2-s",
					H2: 2,
					M1: map[string]string{
						"cows": "are-flying",
						"cats": "are-swimming",
					},
					TT2s: []TT2{
						{
							S1: []int{34, 35, 36},
						},
					},
				},
			},
		},
	}

	r, err := Query(v, "F1")
	assertNoErr(t, err)
	assertEqual(t, "a-string", r)

	_, err = Query(v, "C1")
	assertErr(t, err)

	r, err = Query(v, "T2/G2")
	assertNoErr(t, err)
	assertEqual(t, 433, r)

	r, err = Query(v, "T2/G3")
	assertNoErr(t, err)
	assertEqual(t, 12.433, r)

	r, err = Query(v, "T2/S1/1")
	assertNoErr(t, err)
	assertEqual(t, 3, r)

	r, err = Query(v, "T2/S2/0/H1")
	assertNoErr(t, err)
	assertEqual(t, "h1-s", r)

	r, err = Query(v, "T2/S2/1/M1/cats")
	assertNoErr(t, err)
	assertEqual(t, "are-swimming", r)
}

func TestSet(t *testing.T) {
	v := TT1{
		F1: "a-string",
		F2: 42,
		F3: 0.42,
		T2: TT2{
			G1: "a-g-string",
			G2: 433,
			G3: 12.433,
			S1: []int{2, 3, 4},
			S2: []TT3{
				{
					H1: "h1-s",
					H2: 1,
				},
				{
					H1: "h2-s",
					H2: 2,
					M1: map[string]string{
						"cows": "are-flying",
						"cats": "are-swimming",
					},
					TT2s: []TT2{
						{
							S1: []int{34, 35, 36},
						},
					},
				},
			},
		},
	}

	err := Set(&v, "T2/S2/0/H2", 2)
	assertNoErr(t, err)
	assertEqual(t, 2, v.T2.S2[0].H2)

	err = Set(&v, "F1", "hola")
	assertNoErr(t, err)
	assertEqual(t, "hola", v.F1)

	err = Set(&v, "T2/S2/1/TT2s/0/S1/2", 42)
	assertNoErr(t, err)
	assertEqual(t, 42, v.T2.S2[1].TT2s[0].S1[2])

	// err = Set(&v, "T2/S2/1/M1/cats", "get milk")
	// assertNoErr(t, err)
	// assertEqual(t, "get milk", v.T2.S2[1].M1["cats"])
}

func TestMapValue(t *testing.T) {
	m := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	rv := reflect.ValueOf(m)
	mrv1 := rv.MapIndex(reflect.ValueOf("k1"))
	fmt.Printf("%s\n", mrv1.String())

	if !mrv1.CanAddr() {
		fmt.Printf("value cannot addr\n")
	}
	if !mrv1.CanSet() {
		fmt.Printf("value cannot be set\n")
	}
}
