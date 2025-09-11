// This constants should be synced with users.action table in your db

package configs

const WRITE = "write"
const READ = "read"
const UPDATE = "update"
const DELETE = "delete"
const READ_OWN = "readOwn"
const UPDATE_OWN = "updateOwn"
const DELETE_OWN = "deleteOwn"

// Function that checks if user can read anything
type CanReadRole func(int64) bool
