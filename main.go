package main

import (
	"fmt"
	"os"
	"strings"

	"encoding/xml"
)

type Address struct {
	XMLName    xml.Name `xml:"address"`
	Street     string   `xml:"street"`
	City       string   `xml:"city"`
	PostalCode *string  `xml:"postal_code"`
}

type Phone struct {
	XMLName  xml.Name `xml:"phone"`
	Mobile   string   `xml:"mobile"`
	Cellular string   `xml:"landline"`
}

type Family struct {
	XMLName xml.Name `xml:"family"`
	Name    string   `xml:"name"`
	Born    string   `xml:"born"`
	Address *Address
	Phone   *Phone
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstname"`
	LastName  string   `xml:"lastname"`
	Address   *Address
	Phone     *Phone
	Family    []*Family
}

type People struct {
	XMLName xml.Name `xml:"people"`
	People  []*Person
}

type Args struct {
	InputPath  string
	OutputPath string
}

func parseArgs(_args []string) (args Args, err error) {
	for i := 0; i < len(_args); i++ {
		arg := _args[i]

		switch arg {
		case "-o", "--output":
			if i+1 >= len(_args) {
				err = fmt.Errorf("Missing value for %s", arg)
				return
			}
			args.OutputPath = _args[i+1]
			i++
		default:
			if len(arg) > 0 && arg[0] == '-' {
				err = fmt.Errorf("Unknown flag: %s", _args)
				return
			}
			if args.InputPath != "" {
				err = fmt.Errorf("Only one input file can be provided")
				return
			}
			args.InputPath = arg
		}
	}

	if args.InputPath == "" {
		err = fmt.Errorf("No input file path given")
		return
	}

	return
}

func main() {
	args, err := parseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\nUsage: %s <file> [-o <output>]\n", err, os.Args[0])
		os.Exit(1)
	}

	data, err := os.ReadFile(args.InputPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file %s: %v\n", args.InputPath, err)
		os.Exit(1)
	}

	output := People{}

	lines := strings.Split(string(data), "\n")

	var person *Person
	var family *Family

	for _, line := range lines {
		sections := strings.Split(line, "|")
		for i := range sections {
			sections[i] = strings.TrimSpace(sections[i])
		}
		t := sections[0]
		switch t {
		case "P":
			person = &Person{
				FirstName: sections[1],
				LastName:  sections[2],
			}
			output.People = append(output.People, person)
			family = nil
		case "T":
			{
				phone := &Phone{
					Mobile:   sections[1],
					Cellular: sections[2],
				}
				if family != nil {
					family.Phone = phone
				} else {
					person.Phone = phone
				}
			}
		case "A":
			{
				address := &Address{
					Street: sections[1],
					City:   sections[2],
				}
				if len(sections) == 4 {
					address.PostalCode = &sections[3]
				}
				if family != nil {
					family.Address = address
				} else {
					person.Address = address
				}
			}
		case "F":
			family = &Family{
				Name: sections[1],
				Born: sections[2],
			}
			person.Family = append(person.Family, family)
		}
	}
	data, err = xml.MarshalIndent(output, "", "	")

	if args.OutputPath == "" {
		fmt.Printf("%s\n", data)
	} else if err := os.WriteFile(args.OutputPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file %s: %v\n", args.OutputPath, err)
		os.Exit(1)
	} else {
		fmt.Printf("Wrote %s\n", args.OutputPath)
	}
}
