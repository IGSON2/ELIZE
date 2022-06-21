package elizebch

import "elize/database"

type fakeDB struct {
	fakeLastBLockPoint func() []byte
	fakeFindBlock      func(string) []byte
}

func (f fakeDB) LastBlockPoint() []byte {
	database.InitDB()
	return database.RestoreChain()
}
func (f fakeDB) FindBlock(hash string) []byte {
	return database.FindOneBlock(hash)
}
func (f fakeDB) SaveBlock(hash string, data []byte) {}
func (f fakeDB) SaveBlockchain(data []byte)         {}
func (f fakeDB) DeleteBlocks()                      {}
