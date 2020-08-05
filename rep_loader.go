package main

// ReplicationLoader provides an ability 
// to load replication via InnerReplication
type ReplicationLoader struct {
	FileLoader
}

// GetReplicationFiles looks for files with *.rep pattern
func (loader *ReplicationLoader) GetReplicationFiles() []string {
	return loader.GetFiles("*.rep")
}