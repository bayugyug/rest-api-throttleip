
##TODO
env DEPNOLOCK=1 dep init -v && go test ./...

@users
    POST     /v1/api/user
    PUT      /v1/api/user
    GET      /v1/api/user/{id}
    DELETE   /v1/api/user/{id}
    POST     /v1/api/otp
    POST     /v1/api/login



@Icecream
    POST     /v1/api/icecream
    PUT      /v1/api/icecream
    GET      /v1/api/icecream/{id}
    DELETE   /v1/api/icecream/{id}

@ingredients
    POST     /v1/api/ingredient
    PUT      /v1/api/ingredient
    GET      /v1/api/ingredient/{id}
    DELETE   /v1/api/ingredient/{id}


@sourcing_values
    POST     /v1/api/sourcing
    PUT      /v1/api/sourcing
    GET      /v1/api/sourcing/{id}
    DELETE   /v1/api/sourcing/{id}

