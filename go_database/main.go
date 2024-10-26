package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const VERSION = "1.0.0"

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

type Driver struct {
	mu  sync.Mutex
	mus map[string]*sync.Mutex
	dir string
	log Logger
}

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir: dir,
		mus: make(map[string]*sync.Mutex),
		log: opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Create the database at '%s'\n", dir)
	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write(collection string, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	path := filepath.Join(dir, resource+".json")

	tempPath := path + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := os.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("Missing collection")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource")
	}

	record := filepath.Join(d.dir, collection, resource)
	if _, err := stat(record); err != nil {
		return err
	}

	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)

}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Missing collection")
	}

	colPath := filepath.Join(d.dir, collection)

	if _, err := stat(colPath); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(colPath)
	if err != nil {
		return nil, fmt.Errorf("Cannot read collection: %s", err)
	}
	var records []string
	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(colPath, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("Cannot read collection: %s", err)
		}
		records = append(records, string(b))
	}
	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
	if collection == "" {
		return fmt.Errorf("No Collection specified")
	}
	mu := d.getOrCreateMutex(collection)
	mu.Lock()
	defer mu.Unlock()

	path := filepath.Join(d.dir, collection, resource)

	switch fi, err := stat(path); {
	case fi == nil, err != nil:
		return fmt.Errorf("Cannot find the record: %s", err)
	case fi.Mode().IsDir():
		return os.RemoveAll(path)
	case fi.Mode().IsRegular():
		return os.RemoveAll(path + ".json")
	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, ok := d.mus[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mus[collection] = m
	}

	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

type User struct {
	Name    string  `json:"name"`
	Age     int32   `json:"age"`
	Contact string  `json:"contact"`
	Company string  `json:"compony"`
	Address Address `json:"address"`
}

type Address struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	employes := []User{
		{"John", 23, "Doe", "John Doe", Address{"John", "Doe", "John Doe"}},
	}

	for _, employee := range employes {
		db.Write("users", employee.Name, User{
			Name:    employee.Name,
			Age:     employee.Age,
			Contact: employee.Contact,
			Company: employee.Company,
			Address: employee.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error ReadAll", err)
		return
	}
	fmt.Println(records)
	allusers := []User{}
	for _, record := range records {
		var employeeFound User
		if err := json.Unmarshal([]byte(record), &employeeFound); err != nil {
			fmt.Println("Err Unmarshalling", err)
		}
		allusers = append(allusers, employeeFound)
	}
	fmt.Println(allusers)

	if err := db.Delete("users", "John"); err != nil {
		fmt.Println("Delete error", err)
		return
	}

	if err := db.Delete("users", ""); err != nil {
		fmt.Println("Delete all error", err)
		return
	}
}
