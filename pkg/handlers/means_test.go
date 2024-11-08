package handlers

import (
	"slices"
	"testing"
)

// TestInsert tests the Insert function of the BST.
func TestInsert(t *testing.T) {
	bst := &BST{}

	bst.Insert(10, 100)
	bst.Insert(5, 50)
	bst.Insert(20, 200)
	bst.Insert(15, 150)

	// Verify the structure
	if bst.root.qts != 10 {
		t.Errorf("Expected root qts to be 10, got %d", bst.root.qts)
	}
	if bst.root.left.qts != 5 {
		t.Errorf("Expected left child qts to be 5, got %d", bst.root.left.qts)
	}
	if bst.root.right.qts != 20 {
		t.Errorf("Expected right child qts to be 20, got %d", bst.root.right.qts)
	}
	if bst.root.right.left.qts != 15 {
		t.Errorf("Expected right-left child qts to be 15, got %d", bst.root.right.left.qts)
	}
}

// TestSearch tests the Search function of the BST.
func TestSearch(t *testing.T) {
	bst := &BST{}

	// Insert nodes
	bst.Insert(10, 100)
	bst.Insert(5, 50)
	bst.Insert(20, 200)
	bst.Insert(15, 150)
	bst.Insert(25, 250)

	// Search within a range that includes multiple nodes
	result := bst.Search(5, 20)
	expected := []int32{100, 50, 200, 150}

	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Search within a range that includes only one node
	result = bst.Search(10, 10)
	expected = []int32{100}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Search outside the range of inserted nodes
	result = bst.Search(30, 40)
	expected = []int32{}
	if !slices.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
