package regotest

a := "top"
allow {
    some c
    b := "nested"
    input.data[c] == a
}
