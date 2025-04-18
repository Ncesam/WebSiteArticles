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
		USER     string `env:"POSTGRES_USER" env-required:"true" env-description:"Пользователь базы"`
		PASSWORD string `env:"POSTGRES_PASSWORD" env-required:"true" env-description:"Пароль от базы"`
		DATABASE string `env:"POSTGRES_DATABASE" env-required:"true" env-description:"Имя базы данных"`
	}

	AUTH struct {
		ACCESS struct {
			KEY       string `env:"AUTH_ACCESS_KEY" env-required:"true" env-description:"Секрет для access-токенов"`
			ALGORITHM string `env:"AUTH_ACCESS_ALGO" env-default:"HS256" env-description:"Алгоритм шифрования access-токена"`
		}
		REFRESH struct {
			KEY       string `env:"AUTH_REFRESH_KEY" env-required:"true" env-description:"Секрет для refresh-токенов"`
			ALGORITHM string `env:"AUTH_REFRESH_ALGO" env-default:"HS256" env-description:"Алгоритм шифрования refresh-токена"`
		}
	}
}
