function runOnce  {
  { /usr/bin/time $2 ; } 2> /tmp/o 1> /dev/null
  printf "$1 = "
  cat /tmp/o | awk -v N=1 '{print $N"s"}'
}

function run {
  echo ""
  runOnce "$1" "$2"
  runOnce "$1" "$2"
  runOnce "$1" "$2"
}

run "Go" "go run main.go 40" 
