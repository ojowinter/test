// Helper function to be added in the code transformed.

package main

// Adds public names to the map named in "pkg".
func _export(pkg map[interface{}]interface{}, exported []interface{}) {
	for _, v := range exported {
		pkg[v] = v
	}
}
