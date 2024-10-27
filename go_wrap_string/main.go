package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
  // 1 read the file
  fileBytes, err := os.ReadFile("./example.js")
  if err != nil {
    fmt.Printf("Error reading file: %s\n", err)
    os.Exit(1)
  }
  
  // 2 create the import pattern
  pattern := regexp.MustCompile(`import(\s)((.*)(\s+)from(\s+)(.*);|((.*\n){0,})\}(\s)from(\s)(.*);)`)

  // 3 parse the file data with the pattern
  adjusted := pattern.ReplaceAll(fileBytes,[]byte(`/*$0*/`))

  // 4 write back to file
  os.WriteFile("./example-import-commented.js", adjusted, 0755)


  // for _, found := range foundStatements {
  //   fmt.Println(found)
  // }
}
