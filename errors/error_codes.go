package errors

// API error codes
var (
	ErrCodeRecordNotFound   = newErrorCode("RECORD_NOT_FOUND", "Registro não encontrado.")
	ErrCodeInvalidInputData = newErrorCode("INPUT_DATA_ERROR", "Erro nos dados enviados.")
	ErrIncorrectUnmarshal   = newErrorCode("OUTPUT_DATA_ERROR", "Erro no output da API")
	ErrCodeInternal         = newErrorCode("INTERNAL_SERVER_ERROR", "Ocorreu um error no sistema.")
	ErrEchoHTTPGenericError = newErrorCode("HTTP_GENERIC_ERROR", "Ocorreu um error no sistema.")
	ErrCodeUnauthorized     = newErrorCode("ACCESS_UNAUTHORIZED", "Acesso não autorizado.")
	ErrCodeForbidden        = newErrorCode("ACCESS_FORBIDDEN", "Escopo do token não é valido para esta solicitação. Verifique o token antes de tentar novamente.")

	ErrCodeInvalidAuth = newErrorCode("INVALID_AUTH", "Login incorreto. Verifique os dados informados.")

	ErrCodeInvalidProductID = newErrorCode("INVALID_PRODUCT_ID", "Id produto (productId) inválido.")
	ErrCodeInvalidCart      = newErrorCode("INVALID_CART", "Cesta de compras inválida.")
	ErrCodeEmptyCart        = newErrorCode("EMPTY_CART", "Cesta de compras vazia.")

	ErrCodeUserAlreadyExists = newErrorCode("USER_ALREADY_EXISTS", "Usuário (e-amil) já existe na base de dados")
)
