package main

import (
	"flag"
	"fmt"
	"os"
)

type Celsius float64
type Fehrenheit float64
type Kelvin float64

func (c Celsius) ToFehrenhet() Fehrenheit {
  return Fehrenheit( c * (9.0/5.0) + 32.0 )
}

func (c Celsius) ToKelvin() Kelvin {
  return Kelvin(c + 273.15)
}

func (f Fehrenheit) ToCelsius() Celsius {
  return Celsius((f - 32.0) / (9.0/5.0))
}

func (f Fehrenheit) ToKelvin() Kelvin {
  return f.ToCelsius().ToKelvin()
}


func FToC(f Fehrenheit) Celsius {
  return Celsius((f - 32.0) / (9.0/5.0))
}

func CToF(c Celsius) Fehrenheit {
  return Fehrenheit( c * (9.0/5.0) + 32.0 )
}

func CToK(c Celsius) Kelvin {
  return Kelvin(c + 273.15)
}

func FToK(f Fehrenheit) Kelvin {
  return CToK(FToC(f))
}

func KToC(k Kelvin) Celsius {
  return Celsius(k - 273.15) 
}

func KToF(k Kelvin) Fehrenheit {
  return CToF(KToC(k)) 
}

func (c Celsius) String() string {
  return fmt.Sprintf("%g°C", c)
}

func (f Fehrenheit) String() string {
  return fmt.Sprintf("%g°F", f)
}

func (k Kelvin) String() string {
  return fmt.Sprintf("%gK", k)
}

type celsiusFlag struct {
  Celsius
}

type kelvinFlag struct {
  Kelvin
}

func (cf *celsiusFlag) Set(input string) error {
  var value float64
  var unit string
  fmt.Sscanf(input, "%f%s", &value, &unit)
  switch unit {
  case "C", "c", "°C", "°c":
    cf.Celsius = Celsius(value)
    return nil
  case "F", "f", "°F", "°f":
    cf.Celsius = FToC(Fehrenheit(value))
    return nil
  case "K", "k", "°K", "°k":
    cf.Celsius = KToC(Kelvin(value))
    return nil
  default:
    return fmt.Errorf("invalid format or unit")
  }
}

func (kf *kelvinFlag) Set(input string) error {
  var value float64
  var unit string
  fmt.Sscanf(input, "%f%s", &value, &unit)
  switch unit {
  case "C", "c", "°C", "°c":
    kf.Kelvin = CToK(Celsius(value))
    return nil
  case "F", "f", "°F", "°f":
    kf.Kelvin = FToK(Fehrenheit(value))
    return nil
  case "K", "k", "°K", "°k":
    kf.Kelvin = Kelvin(value)
    return nil
  default:
    return fmt.Errorf("invalid format or unit")
  }
}

func CelciusFlag(name string, defaultVal Celsius, usage string) *Celsius {
  f := celsiusFlag{defaultVal}
  flag.CommandLine.Var(&f, name, usage)
  return &f.Celsius
}

func KelvinFlag(name string, defaultVal Kelvin, usage string) *Kelvin {
  f := kelvinFlag{defaultVal}
  flag.CommandLine.Var(&f, name, usage)
  return &f.Kelvin
}

func main() {
  c := CelciusFlag("c", 18, "the temperature in celsius")
  k := KelvinFlag("k", 18, "the temperature in kelvin")
  flag.Parse()
  fmt.Println(*c)
  fmt.Println(*k)
}
