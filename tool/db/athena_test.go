package db

import "testing"

func TestInitAthena(t *testing.T) {
	db, err := InitAthena(
		WithAccessID(""),
		WithAccessKey(""),
		WithDatabase(""),
		WithRegion(""),
		WithLogPath(""),
	)
	if err != nil {
		t.Error(err)
	}
	_ = db
}

func TestStorageAthenaDatabase(t *testing.T) {
	err := StorageAthenaDatabase("",
		WithDatabase(""),
		WithRegion(""),
		WithAccessKey(""),
		WithAccessID(""),
		WithLogPath(""))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLoadAthenaDatabase(t *testing.T) {
	//只获取database实例,如果获取不到,返回空
	db, err := LoadAthenaDatabase("")
	if err != nil {
		t.Error(err)
		return
	}
	//获取database实例,如果实例为空,通过配置初始化并存储至内存
	db, err = LoadAthenaDatabase("",
		WithDatabase(""),
		WithRegion(""),
		WithAccessKey(""),
		WithAccessID(""),
		WithLogPath(""))
	if err != nil {
		t.Error(err)
		return
	}
	_ = db
}

func TestDelAthenaDatabase(t *testing.T) {
	//在内存中移除 app_key 对应的database实例,并将实例关闭
	DelAthenaDatabase("")
}
