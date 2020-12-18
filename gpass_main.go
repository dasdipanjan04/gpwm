package main

import (
	mkm "gpwm/masterkeymanager"
)

func main() {
	db := mkm.CreateMasterKeyTable()
	mkm.InsertMasterKeyDataToDB(db, "TestFN_1", "TestLN_1", "1test@test.test", "tesertewtertettest", "pazsword", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_2", "TestLN_2", "2test@test.test", "ertwertertertertert", "word", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_3", "TestLN_3", "3test@test.test", "testtertertertertest", "pward", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_4", "TestLN_4", "4test@test.test", "ertertertert", "passkey", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_5", "TestLN_5", "5test@test.test", "qrewtewrtwertwert", "p@$$stuff", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_6", "TestLN_6", "6test@test.test", "testfdvdfbbvbtest", "PPaasswwoorrdd", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_7", "TestLN_7", "7test@test.test", "xcvbxcvbxxvcbcvbxcvb", "paaassssssdddd", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_8", "TestLN_8", "8test@test.test", "aerqwer2rdadxv", "pswd", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_9", "TestLN_9", "9test@test.test", "eufztrqwiefuqgkjahfvkj", "cred", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_10", "TestLN_10", "10test@test.test", "@@jhfgequfgfvhkcvjh", "entials", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_11", "TestLN_11", "11test@test.test", "dhfgaskdfhgcvmc{{2397813", "holydfgfgf", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_12", "TestLN_12", "12test@test.test", "kruzgriugv23235467ffdhg0==d", "hello", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_13", "TestLN_13", "13test@test.test", "kurzefrage1kurzefrage2achso??", "golang", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_14", "TestLN_14", "14test@test.test", "dvasdvavdsdv", "alltest", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_15", "TestLN_15", "15test@test.test", "fdfwrwer345fgrqtdfsda", "testpass", true)
	mkm.InsertMasterKeyDataToDB(db, "TestFN_16", "TestLN_16", "16test@test.test", "adsfqwert4t114355", "trespass", true)
}
