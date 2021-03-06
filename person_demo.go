package main

import (
	"fmt"

	"github.com/fmagellan/ygotDemo/pkg/person"
	"github.com/openconfig/ygot/ygot"
)

// The following generation rule uses the generator binary to create the pkg/person package

//go:generate generator -path=yang -output_file=pkg/person/person.go -package_name=person -generate_fakeroot -fakeroot_name=device -compress_paths=true -shorten_enum_leaf_names -typedef_enum_with_defmod yang/person.yang

func main() {
	// Create a new device which is named according to the fake root specified above. To generate
	// the fakeroot then generate_fakeroot should be specified. This entity corresponds to the
	// root of the YANG schema tree. The fakeroot name is the CamelCase version of the name
	// supplied by the fakeroot_name argument.
	d := &person.Device{}

	// To render the device (which is currently empty) to JSON in RFC7951 format, then we
	// simply call the ygot.EmitJSON method with the relevant arguments.
	jsonConfig := &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	}
	jsonString, err := ygot.EmitJSON(d, jsonConfig)

	// If an error was returned (which occurs if the struct's contents could not be validated
	// or an error occurred with rendering to JSON), then this should be handled by the
	// calling code.
	if err != nil {
		panic(err)
	}
	fmt.Printf("Empty JSON: %v\n", jsonString)

	p := new(person.Person)
	p.Age = ygot.Uint32(33)
	p.Name = ygot.String("Magellan")

	d.Person = p

	// We can then validate the contents of the interface that we created.
	if err := d.Person.Validate(); err != nil {
		panic(fmt.Sprintf("Person validation failed: %v", err))
	}

	// We can also validate the device overall.
	if err := d.Validate(); err != nil {
		panic(fmt.Sprintf("Device validation failed: %v", err))
	}

	// EmitJSON from the ygot library directly does .Validate() and outputs JSON in
	// the specified format.
	jsonString, err = ygot.EmitJSON(d, jsonConfig)
	if err != nil {
		panic(fmt.Sprintf("JSON demo error: %v", err))
	}
	fmt.Println(jsonString)

	// The generated code includes an Unmarshal function, which can be used to load
	// a data tree such as the one that we just created.
	unmarshalledResult := &person.Device{}
	if err := person.Unmarshal([]byte(jsonString), unmarshalledResult); err != nil {
		panic(fmt.Sprintf("Can't unmarshal JSON: %v", err))
	}

	fmt.Println("---After unmarshalling---")
	fmt.Println("Name:", *unmarshalledResult.Person.Name)
	fmt.Println("Age:", *unmarshalledResult.Person.Age)
}
