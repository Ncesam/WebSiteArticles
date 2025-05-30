package config

type Config struct {
	VERSION string `env:"VERSION" env-default:"0.1" env-description:"Версия приложения"`

	APP struct {
		PORT  int    `env:"APP_PORT" env-default:"8080" env-description:"Порт на котором запускается приложение"`
		HOST  string `env:"APP_HOST" env-default:"localhost" env-description:"Хост"`
		DEBUG bool   `env:"APP_DEBUG" env-description:"Режим отладки"`
	}

	POSTGRES struct {
		PORT     int    `env:"POSTGRES_PORT" env-default:"5432" env-description:"Порт базы данных"`
		HOST     string `env:"POSTGRES_HOST" env-default:"postgres" env-description:"Хост базы"`
		USERNAME string `env:"POSTGRES_USER" env-required:"true" env-description:"Пользователь базы"`
		PASSWORD string `env:"POSTGRES_PASSWORD" env-required:"true" env-description:"Пароль от базы"`
		DATABASE string `env:"POSTGRES_DATABASE" env-required:"true" env-description:"Имя базы данных"`
	}
	MONGO struct {
		PORT     int    `env:"MONGO_PORT" env-default:"27017" env-description:"Порт MongoDB"`
		HOST     string `env:"MONGO_HOST" env-default:"mongo" env-description:"Хост MongoDB"`
		USERNAME string `env:"MONGO_USER" env-description:"Пользователь MongoDB (опционально)"`
		PASSWORD string `env:"MONGO_PASSWORD" env-description:"Пароль MongoDB (опционально)"`
		DATABASE string `env:"MONGO_DATABASE" env-default:"admin" env-description:"Имя базы данных MongoDB"`
	}

	AUTH struct {
		ACCESS struct {
			DURATION  int    `env:"AUTH_ACCESS_DURATION" env-required:"true" env-description:"Время жизни токена"`
			KEY       string `env:"AUTH_ACCESS_KEY" env-required:"true" env-description:"Секрет для access-токенов"`
			ALGORITHM string `env:"AUTH_ACCESS_ALGO" env-default:"HS256" env-description:"Алгоритм шифрования access-токена"`
		}
		REFRESH struct {
			DURATION  int    `env:"AUTH_REFRESH_DURATION" env-required:"true" env-description:"Время жизни токена"`
			KEY       string `env:"AUTH_REFRESH_KEY" env-required:"true" env-description:"Секрет для refresh-токенов"`
			ALGORITHM string `env:"AUTH_REFRESH_ALGO" env-default:"HS256" env-description:"Алгоритм шифрования refresh-токена"`
		}
	}
	AUTH_SERVICE struct {
		ADDRESS string `env:"AUTH_SERVICE_ADDRESS" env-required:"true" env-description:"Адрес микросервиса"`
	}
	CONFIG_SERVICE struct {
		ADDRESS string `env:"CONFIG_SERVICE_ADDRESS" env-required:"true" env-description:"Адрес микросервиса"`
	}
	QUEUE_SERVICE struct {
		ADDRESS string `env:"QUEUE_SERVICE_ADDRESS" env-required:"true" env-description:"Адрес микросервиса"`
	}
	REQUEST_SERVICE struct {
		ADDRESS string `env:"REQUEST_SERVICE_ADDRESS" env-required:"true" env-description:"Адрес микросервиса"`
	}
	OPEN_AI_KEY string `env:"OPEN_AI_KEY" env-required:"false" env-description:"Ключ API Open AI"`
	CLID_DTF    string `env:"CLID_DTF" env-required:"false" env-description:"CLID для яндекс маркета DTF"`
	CLID_VC     string `env:"CLID_VC" env-required:"false" env-description:"CLID для яндекс маркета VC"`
}
