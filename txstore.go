package txstore

// Store is a transaction storage specified in the document linked below.
// https://www.notion.so/Transactional-Key-Value-Store-d72f26aa31e34eef9aa7442507215ce7
type Store interface {
	// Set store the value for key
	Set(key, value string)
	// Get return the current value for key
	Get(key string) string
	// Delete remove the entry for key
	Delete(key string)
	// Count return the number of keys that have the given value
	Count(value string) int
	// Begin start a new transaction
	Begin()
	// Commit complete the current transaction
	Commit()
	// Rollback revert to state prior to BEGIN call
	Rollback()
}
