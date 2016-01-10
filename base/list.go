package base

import (
	"fmt"
	"sort"
)

func getLocalList() ([]License, error) {
	content, err := readIndex()

	if err != nil {
		return nil, err
	}

	return jsonToList(content)
}

func getRemoteList() ([]License, error) {
	body, err := fetchIndex()

	if err != nil {
		return nil, err
	}

	return jsonToList(body)
}

// printList prints the provided list of licenses
// after sorting them. Side-effect: the underlying
// array for the slice is sorted.
func printList(licenses []License) {
	sort.Sort(ByLicenseKey(licenses))

	fmt.Print("Available licenses:\n\n")
	for _, l := range licenses {
		fmt.Printf("%s%-14s(%s)\n", indent, l.Key, l.Name)
	}
	fmt.Println()
}

// ListLocal reads the list of available local licenses
// and prints the list.
func ListLocal() error {
	licenses, err := getLocalList()

	if err != nil {
		return newErrReadFailed()
	}

	printList(licenses)
	return nil
}

// ListRemote fetches the list of remote licenses
// and prints the list.
func ListRemote() error {
	licenses, err := getRemoteList()

	if err != nil {
		return newErrFetchFailed()
	}

	printList(licenses)
	return nil
}
