package client

import (
	"fmt"
	"git.sr.ht/~akilan1999/p2p-rendering-computation/server/docker"
	"testing"
)

// Testing out if a new group is getting created
func TestCreateGroup(t *testing.T) {
	group, err := CreateGroup()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	PrettyPrint(group)
}

// Testing if the group gets removed when a
// group ID is provided
func TestRemoveGroup(t *testing.T) {
	// Creates a new group
	group, err := CreateGroup()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Removes the new group
	// it created
	err = RemoveGroup(group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

// Testing if container information is added
// to the created group
func TestAddContainerToGroup(t *testing.T) {
	// Creates a new group
	group, err := CreateGroup()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Creating and adding the container to the
	// tracked list
	container1 ,err := docker.BuildRunContainer(0,"false","")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Testing the AddTrackContainer Function and adding the first container created
	err = AddTrackContainer(container1,"0.0.0.0")
	if err != nil {
		// Killing docker container created
		err = docker.StopAndRemoveContainer(container1.ID)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(err)
		t.Fail()
	}

	// Adds container information to the group
	Group, err := AddContainerToGroup(container1.ID,group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	PrettyPrint(Group)

	// Killing docker container created
	err = docker.StopAndRemoveContainer(container1.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Removing container 1 from the tracked list
	err = RemoveTrackedContainer(container1.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Removes the new group
	// it created
	err = RemoveGroup(group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}

// Testing if the container information is removed from the group
func TestGroup_RemoveContainerGroup(t *testing.T) {
	// Creates a new group
	group, err := CreateGroup()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Creating and adding the container to the
	// tracked list
	container1 ,err := docker.BuildRunContainer(0,"false","")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Testing the AddTrackContainer Function and adding the first container created
	err = AddTrackContainer(container1,"0.0.0.0")
	if err != nil {
		// Killing docker container created
		err = docker.StopAndRemoveContainer(container1.ID)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		fmt.Println(err)
		t.Fail()
	}

	// Adds container information to the group
	Group, err := AddContainerToGroup(container1.ID,group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("Container added")
	PrettyPrint(Group)

	// Removing docker container from the group
	Group, err = RemoveContainerGroup(container1.ID,group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("Container removed")
	PrettyPrint(Group)

	// Killing docker container created
	err = docker.StopAndRemoveContainer(container1.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Removing container 1 from the tracked list
	err = RemoveTrackedContainer(container1.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// Removes the new group
	// it created
	err = RemoveGroup(group.ID)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
