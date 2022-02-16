package models

type HTTP404 struct {
	error
}

type HTTP500 struct {
	error
}

type PSQLStorage struct {
	error
}

type WrongStorageType struct {
	error
}

type MissingStorageType struct {
	error
}

type RunError struct {
	error
}
