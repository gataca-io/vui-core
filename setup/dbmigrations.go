package setup

//Scripts must be retrocompatible. They will be executed on each start of the app.
var DBScripts []string = []string{
	//v1.1.0: mark tokenConfigs with default flag
	"UPDATE \"tokenConfigs\" set is_default = false where is_default is NULL",
	"CREATE UNIQUE index IF NOT EXISTS one_default on \"tokenConfigs\" (is_default) where is_default = true",
}
