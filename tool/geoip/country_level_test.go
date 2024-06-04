package geoip

import (
	"fmt"
	"testing"
)

func TestCountryByCode(t *testing.T) {
	fmt.Println(CountryByCode("USSR"))
}
