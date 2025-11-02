// internal/models/presets.go
package models

type PresetRecord map[string]interface{}

func GetPresetUsers() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "name": "Alice", "email": "alice@example.com", "age": 25, "isActive": true},
		{"id": "2", "name": "Bob", "email": "bob@example.com", "age": 30, "isActive": false},
	}
}

func GetPresetProducts() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "title": "Wireless Mouse", "price": 499, "inStock": true},
		{"id": "2", "title": "Mechanical Keyboard", "price": 1999, "inStock": false},
	}
}

func GetPresetComments() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "postId": "101", "author": "Alice", "content": "Nice post!"},
		{"id": "2", "postId": "102", "author": "Bob", "content": "I totally agree!"},
	}
}

func GetPresetCarts() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "userId": "1", "productId": "2", "quantity": 1},
		{"id": "2", "userId": "2", "productId": "1", "quantity": 2},
	}
}

func GetPresetPosts() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "title": "Welcome to MockNode", "content": "This is a sample post!"},
		{"id": "2", "title": "Building APIs", "content": "Learn how to mock APIs easily."},
	}
}

func GetPresetTodos() []PresetRecord {
	return []PresetRecord{
		{"id": "1", "task": "Finish project setup", "done": false},
		{"id": "2", "task": "Test API endpoints", "done": true},
	}
}
