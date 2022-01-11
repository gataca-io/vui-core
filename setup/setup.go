package setup

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gataca-io/vui-core/admin"
	adminC "github.com/gataca-io/vui-core/admin/controller"
	adminD "github.com/gataca-io/vui-core/admin/dao"
	adminS "github.com/gataca-io/vui-core/admin/service"
	"github.com/gataca-io/vui-core/log"
	"github.com/gataca-io/vui-core/models"
	"github.com/gataca-io/vui-core/security"
	"github.com/gataca-io/vui-core/service"
	"github.com/gataca-io/vui-core/service/impl"
	"github.com/gataca-io/vui-core/tools"
	vcApis "github.com/gataca-io/vui-core/vc-apis/controller"
	"github.com/spf13/viper"

	"github.com/AlexanderGrom/go-event"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	basePath = "/"
)

type Services struct {
	DefaultVMethod             models.DefaultVMethod
	ConfigRepo                 admin.TokenTenantDao
	ConfigReader               service.ConfigReader
	Tokenizer                  service.Tokenizer
	CryptoService              service.CryptoService
	EidasBridge                service.IEidasBridge
	DidService                 service.DidService
	GovernanceService          service.GovernanceService
	SSIService                 service.SSIService
	NotificationService        service.NotificationService
	GatacaLoginService         admin.GatacaLoginService
	ConnectService             admin.ConnectService
	CatalogService             admin.CatalogService
	CatalogPresentationService admin.CatalogPresentationService
	SetupService               admin.SetupService
	EOSService                 service.EOSService
	ETHService                 service.EthereumService
}

// InitServer Create a new Echo Server
func InitServer() *echo.Echo {
	e := echo.New()
	if len(viper.GetStringMap("app.traces")) > 0 {
		tH := tools.GetOrElse("app.traces.traceIdHeader", echo.HeaderXRequestID)
		sH := tools.GetOrElse("app.traces.spanIdHeader", log.HeaderXSpanId)
		e.Pre(log.RequestIDWithHeaders(tH, sH))
	} else {
		e.Pre(log.RequestID())
	}
	// Middleware
	e.Use(log.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowMethods:  []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS, echo.HEAD},
		ExposeHeaders: []string{"Token", "Token_type", "X-Connect-Id", "X-Connect-Jwt"},
	}))

	log.Debug("Installed & configured echo middlewares")

	if viper.GetBool("app.metrics.enabled") {
		log.Debug("Enabling metrics...")
		go SetupMetrics(viper.GetString("app.metrics.address"))
		log.Info("Metrics enabled")
	}
	return e
}

// InitDb Create a New DB connection
func InitDb(dbHost string, dbPort string, dbUser string, dbPass string, dbName string, dbSchema string, idleConns string, maxOpenConns string, connTTL string, sslMode string) (*gorm.DB, error) {
	var dataSourceName string
	if dbSchema != "" {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, sslMode, dbSchema)
	} else {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, sslMode)
	}
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	idle, err := strconv.Atoi(idleConns)
	if err != nil {
		return nil, fmt.Errorf("invalid postgresql idleConns param: %v", idleConns)
	}
	maxOpen, err := strconv.Atoi(maxOpenConns)
	if err != nil {
		return nil, fmt.Errorf("invalid postgresql maxOpenConns param: %v", maxOpenConns)
	}
	ttl, err := strconv.Atoi(connTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid postgresql connTTL param: %v", connTTL)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error configuring DB")
	}

	sqlDB.SetMaxIdleConns(idle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(ttl) * time.Second)

	// Check connection
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitConfigRepo Initialize the config repository
func InitConfigRepo(s *Services, dbMap *gorm.DB, dbConfig bool) error {
	if s.ConfigRepo == nil {
		if dbConfig {
			configRepo := adminD.NewTokenTenantDao(dbMap)
			s.ConfigRepo = configRepo
			config := configRepo.GetDefaultConfig()
			if config == nil {
				log.Error("Config not initialized", config)
				return models.ErrNotInitialized
			}
		} else {
			s.ConfigRepo = adminD.NewFileConfig()
		}
	}
	return nil
}

// InitConfigReader Initialize a config reader
func InitConfigReader(s *Services) {
	if s.ConfigReader == nil {
		s.ConfigReader = service.NewConfigReader(s.ConfigRepo)
		s.DefaultVMethod = s.ConfigRepo.GetDefaultVMethod
	}
}

// InitCrypto Create the services for crypto
func InitCrypto(s *Services, externalUrls []string) {
	if s.ConfigReader != nil {
		if s.Tokenizer == nil {
			s.Tokenizer = impl.NewTokenizer(s.ConfigReader, externalUrls)
		}
		if s.CryptoService == nil {
			s.CryptoService = impl.NewCryptoService(s.ConfigReader)
		}
	}
}

// InitSecurity Set up the enforcement of security
func InitSecurity(s *Services, e *echo.Echo, dbMap *gorm.DB, loadedPaths []interface{}) {
	repo := security.AccessRepository{DbMap: dbMap}
	enforcer := &security.Enforcer{
		ConfigPaths: loadedPaths,
		Tokenizer:   s.Tokenizer,
		Repo:        &repo,
	}
	enforcer.Install(e)
}

// SetupEidasBridge Configure an Eidas Bridge using the config
func SetupEidasBridge(s *Services, configKeys []interface{}) {
	if s.EidasBridge == nil {
		if len(configKeys) > 0 {
			log.Debug("Configuring eiDas Bridge")
			qKeys := make(map[string]models.QECKey)
			for _, key := range configKeys {
				qKey := key.(map[string]interface{})
				password := qKey["password"].(string)
				if password == "" {
					password = viper.GetString("eidasBridge.password")
				}
				qecHex := ""
				qecBytes := qKey["qecBytes"]
				if qecBytes == nil {
					qecFile := qKey["qecFile"]
					if qecFile != nil {
						file := qecFile.(string)
						file = filepath.Join(basePath, filepath.Clean(file))
						/* #nosec G304 */
						bytes, err := ioutil.ReadFile(file)
						if err != nil {
							log.Error("Cannot read P12 of QEC file at ", file)
							return
						}
						qecHex = hex.EncodeToString(bytes)
					} else {
						log.Info("No EidasBridge certificate file provided.")
					}
				} else {
					qecHex = qecBytes.(string)
				}
				did := qKey["did"].(string)
				qec := models.QECKey{
					Host:     qKey["host"].(string),
					Password: password,
					QEC:      qecHex,
					DID:      did,
				}
				qKeys[did] = qec
			}
			eb := service.NewEidasBridge(qKeys)
			for _, qKey := range qKeys {
				err := eb.CreateKey(nil, &qKey)
				if err != nil {
					log.Error("Cannot register key in Eidas Bridge", err)
				}
			}
			s.EidasBridge = eb
		}
	}
}

// InitAdminCatalog Create a service to query the remote catalog and the local catalog
func InitAdminCatalog(s *Services, dbMap *gorm.DB, ev event.Dispatcher) {
	if s.CatalogService == nil {
		catalogDao := adminD.NewPostgreSQLCatalogRepository(dbMap)
		initCatalogs(catalogDao)
		servCat := adminS.NewCatalogService(catalogDao, ev)
		s.CatalogService = servCat
	}
}

// InitPlatformServices Init all services used to interact with the Gataca Platform
func InitPlatformServices(s *Services, ev event.Dispatcher, didRHost string, governanceHost string, accountsHost string, remoteCatalog bool) {
	if s.DidService == nil {
		loginService := service.NewLoginService(governanceHost, s.CryptoService)
		s.DidService = service.NewDidService(didRHost, loginService)
		s.GovernanceService = service.NewGovernanceService(governanceHost, loginService, s.DidService, s.CatalogService, remoteCatalog, ev)
		s.NotificationService = service.NewNotificationService(accountsHost, loginService)
	}
	if s.SSIService == nil {
		s.SSIService = impl.NewSSIService(s.CryptoService, s.EidasBridge, s.DidService, s.GovernanceService)
	}
}

func InitMonitoringApps(e *echo.Echo) {
	adminC.NewVersionHandler(e)
}

func InitTokenAndAdminApps(s *Services, dbMap *gorm.DB, e *echo.Echo, ev event.Dispatcher, ebsiURL, ebsiOnboardingIssuer string) {
	connectLoginConfigRepo := adminD.NewConnectLoginTenantDao(dbMap)

	if s.EOSService == nil {
		if s.ETHService == nil {
			s.ETHService = service.NewEthereumService(dbMap)
		}
		s.EOSService = service.NewEOSService(s.ConfigReader, s.ETHService, ebsiURL, ebsiOnboardingIssuer, s.GovernanceService, s.SSIService, s.DidService)
	}

	if s.ConnectService == nil {
		s.ConnectService = adminS.NewConnectService(connectLoginConfigRepo)
	}
	if s.GatacaLoginService == nil {
		userDao := adminD.NewUserDao(dbMap)
		s.GatacaLoginService = adminS.NewGatacaLoginService(userDao, s.ConnectService, s.Tokenizer, s.ConfigReader, s.CryptoService)
	}
	if s.SetupService == nil {
		s.SetupService = adminS.NewSetupService(s.ConfigRepo, connectLoginConfigRepo, s.EOSService, s.ETHService)
	}

	adminC.NewAdminHandler(e, s.GatacaLoginService, s.ConfigReader)
	adminC.NewSetupHandler(e, ev, s.SetupService)
}

func InitCatalogApps(s *Services, dbMap *gorm.DB, e *echo.Echo) {
	if s.CatalogPresentationService == nil {
		s.CatalogPresentationService = adminS.NewCatalogPresentationService(s.GovernanceService)
	}

	adminC.NewCatalogHandler(e, s.CatalogService, s.CatalogPresentationService)
	adminC.NewSchemaHandler(e, s.CatalogPresentationService)
}

func InitVCAPIS(s *Services, e *echo.Echo) {
	vcApis.NewVCAPIsHandler(e, s.SSIService)
}

// SetUp Make a generic setup receiving all requested params
func SetUp(e *echo.Echo, dbMap *gorm.DB, ev event.Dispatcher, dbConfig bool, eidasConfig []interface{}, externalUrls []string, loadedPaths []interface{}, backboneHost string, governanceHost string, accountsHost string, remoteCatalog bool, ebsiURL, ebsiOnboardingIssuer string) (*Services, error) {
	log.Debug("Init setup")
	s := &Services{}
	err := InitConfigRepo(s, dbMap, dbConfig)
	InitConfigReader(s)
	SetupEidasBridge(s, eidasConfig)
	InitCrypto(s, externalUrls) //
	InitSecurity(s, e, dbMap, loadedPaths)
	InitAdminCatalog(s, dbMap, ev)
	InitPlatformServices(s, ev, backboneHost, governanceHost, accountsHost, remoteCatalog)
	initApps(s, dbMap, e, ev, ebsiURL, ebsiOnboardingIssuer)
	return s, err
}

// CompleteSetup Create a full setup reading from config and creating all needed objects
func CompleteSetup() (*echo.Echo, *gorm.DB, event.Dispatcher, *Services, error) {
	e := InitServer()

	dbHost := viper.GetString("database.host")
	dbPort := tools.GetOrElse("database.port", "5432")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")
	dbSchema := viper.GetString("database.schema")
	idleConns := tools.GetOrElse("database.connection.idleCons", "10")
	maxOpenConns := tools.GetOrElse("database.connection.maxOpenConns", "100")
	connTTL := tools.GetOrElse("database.connection.connTTL", "300")
	sslMode := tools.GetOrElse("database.connection.sslmode", "require")
	db, err := InitDb(dbHost, dbPort, dbUser, dbPass, dbName, dbSchema, idleConns, maxOpenConns, connTTL, sslMode)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	ev := event.New()

	configKeysData := viper.Get("eidasbridge")
	var configKeys []interface{}
	if configKeysData != nil {
		configKeys = configKeysData.([]interface{})
	}
	loadedPaths := viper.Get("security.paths").([]interface{})
	externalUrls := viper.GetStringSlice("security.externalJwks")
	dids := viper.GetString("app.dids.host")
	accounts := viper.GetString("app.accounts.host")
	governance := viper.GetString("app.governance.host")
	ebsiURL := tools.GetMandatoryString("app.ebsi.url")
	ebsiOnboardingIssuer := tools.GetMandatoryString("app.ebsi.onboardingIssuer")
	s, err := SetUp(e, db, ev, true, configKeys, externalUrls, loadedPaths, dids, accounts, governance, false, ebsiURL, ebsiOnboardingIssuer)

	return e, db, ev, s, err
}

// ##########
// # PRIVATE
// ##########
func initCatalogs(catalogDao admin.CatalogDao) {
	t := time.Now()
	for _, key := range []string{adminS.Authority, adminS.CredentialGroup, adminS.CredentialType} {
		for _, suffix := range []string{"", "_sbx"} {
			_, err := catalogDao.GetByKey(key + suffix)
			if err != nil {
				catalog := models.Catalog{
					Key:       key + suffix,
					Value:     "[]",
					CreatedAt: &t,
					UpdatedAt: &t,
				}
				err := catalogDao.CreateKey(&catalog)
				if err != nil {
					log.Error("Error creating key", key+suffix, err)
				}
			}
		}
	}
}

// initApps Init common core application services
func initApps(s *Services, dbMap *gorm.DB, e *echo.Echo, ev event.Dispatcher, ebsiURL, ebsiOnboardingIssuer string) {
	InitMonitoringApps(e)
	InitTokenAndAdminApps(s, dbMap, e, ev, ebsiURL, ebsiOnboardingIssuer)
	InitCatalogApps(s, dbMap, e)
	InitVCAPIS(s, e)
}

// PerformDBMigrations This is not automatically invoked. It must be invoked from the main.go of the parent app.
func PerformDBMigrations(db *gorm.DB) {
	for _, script := range DBScripts {
		err := db.Exec(script).Error
		if err != nil {
			log.Fatalf("Cannot perform migration %s :\n%s", script, err)
		}
	}
}
