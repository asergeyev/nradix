// Copyright (C) 2015 Alex Sergeyev
// This project is licensed under the terms of the MIT license.
// Read LICENSE file for information for all notices and permissions.

package nradix

import "testing"

func TestTree(t *testing.T) {
	tr := NewTree(0)
	if tr == nil || tr.root == nil {
		t.Error("Did not create tree properly")
	}
	err := tr.AddCIDR("1.2.3.0/25", 1)
	if err != nil {
		t.Error(err)
	}

	// Matching defined cidr
	inf, err := tr.FindCIDR("1.2.3.1/25")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 1 {
		t.Errorf("Wrong value, expected 1, got %v", inf)
	}

	// Inside defined cidr
	inf, err = tr.FindCIDR("1.2.3.60/32")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 1 {
		t.Errorf("Wrong value, expected 1, got %v", inf)
	}
	inf, err = tr.FindCIDR("1.2.3.60")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 1 {
		t.Errorf("Wrong value, expected 1, got %v", inf)
	}

	// Outside defined cidr
	inf, err = tr.FindCIDR("1.2.3.160/32")
	if err != nil {
		t.Error(err)
	}
	if inf != nil {
		t.Errorf("Wrong value, expected nil, got %v", inf)
	}
	inf, err = tr.FindCIDR("1.2.3.160")
	if err != nil {
		t.Error(err)
	}
	if inf != nil {
		t.Errorf("Wrong value, expected nil, got %v", inf)
	}

	inf, err = tr.FindCIDR("1.2.3.128/25")
	if err != nil {
		t.Error(err)
	}
	if inf != nil {
		t.Errorf("Wrong value, expected nil, got %v", inf)
	}

	// Covering not defined
	inf, err = tr.FindCIDR("1.2.3.0/24")
	if err != nil {
		t.Error(err)
	}
	if inf != nil {
		t.Errorf("Wrong value, expected nil, got %v", inf)
	}

	// Covering defined
	err = tr.AddCIDR("1.2.3.0/24", 2)
	if err != nil {
		t.Error(err)
	}
	inf, err = tr.FindCIDR("1.2.3.0/24")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 2 {
		t.Errorf("Wrong value, expected 2, got %v", inf)
	}

	inf, err = tr.FindCIDR("1.2.3.160/32")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 2 {
		t.Errorf("Wrong value, expected 2, got %v", inf)
	}

	// Delete internal
	err = tr.DeleteCIDR("1.2.3.0/25")
	if err != nil {
		t.Error(err)
	}

	// Hit covering with old IP
	inf, err = tr.FindCIDR("1.2.3.0/32")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 2 {
		t.Errorf("Wrong value, expected 2, got %v", inf)
	}

}

func TestTree6(t *testing.T) {
	tr := NewTree(0)
	if tr == nil || tr.root == nil {
		t.Error("Did not create tree properly")
	}
	err := tr.AddCIDR("dead::0/16", 3)
	if err != nil {
		t.Error(err)
	}

	// Matching defined cidr
	inf, err := tr.FindCIDR("dead::beef")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 3 {
		t.Errorf("Wrong value, expected 3, got %v", inf)
	}

	// Outside
	inf, err = tr.FindCIDR("deed::beef/32")
	if err != nil {
		t.Error(err)
	}
	if inf != nil {
		t.Errorf("Wrong value, expected nil, got %v", inf)
	}

	// Subnet
	err = tr.AddCIDR("dead:beef::0/48", 4)
	if err != nil {
		t.Error(err)
	}

	// Match defined subnet
	inf, err = tr.FindCIDR("dead:beef::0a5c:0/64")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 4 {
		t.Errorf("Wrong value, expected 4, got %v", inf)
	}

	// Match outside defined subnet
	inf, err = tr.FindCIDR("dead:0::beef:0a5c:0/64")
	if err != nil {
		t.Error(err)
	}
	if inf.(int) != 3 {
		t.Errorf("Wrong value, expected 3, got %v", inf)
	}

}
