package internal

import (
	env "github.com/Netflix/go-env"
)

type Environment struct {
	Env  string `env:"ENV,required=true"`
	User struct {
		Email    string `env:"USER_EMAIL,required=true"`
		Name     string `env:"USER_NAME,required=true"`
		Password string `env:"USER_PASS,required=true"`
	}

	Post struct {
		Comment      string `env:"POST_COMMENT,required=true"`
		CommentTimes int    `env:"POST_COMMENT_TIMES,required=true"`
		URL          string `env:"POST_URL,required=true"`
	}
}

var Env Environment

func init() {
	_, err := env.UnmarshalFromEnviron(&Env)
	if err != nil {
		panic(err)
	}
}
