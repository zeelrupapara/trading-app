package utils

import (
	"github.com/zeelrupapara/trading-api/models"
	"github.com/zeelrupapara/trading-api/structs"
)

// swagger:parameters RequestCreateUser
type RequestCreateUser struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqRegisterUser
	}
}

// swagger:response ResponseCreateUser
type ResponseCreateUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestGetUser
type RequestGetUser struct {
	// in:path
	UserId string `json:"userId"`
}

// swagger:response ResponseGetUser
type ResponseGetUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestAuthnUser
type RequestAuthnUser struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqLoginUser
	}
}

// swagger:response ResponseAuthnUser
type ResponseAuthnUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestMarketData
type RequestMarketData struct {
	// in:parameters
	// required: true
	Symbol string `json:"symbol"`
}

type MarketData struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

// swagger:response ResponseMarketData
type ResponseMarketData struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			MarketData
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestPlaceOrder
type RequestPlaceOrder struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqPlaceOrder
	}
}

// swagger:response ResponsePlaceOrder
type ResponsePlaceOrder struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			structs.ResPlaceOrder
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestTradeHistory
type RequestTradeHistory struct {
	// in:parameters
	Page string `json:"page"`
}

// swagger:response ResponseTradeHistory
type ResponseTradeHistory struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   []struct {
			structs.ResPlaceOrder
		} `json:"data"`
	} `json:"body"`
}

////////////////////
// --- GENERIC ---//
////////////////////

// Response is okay
// swagger:response GenericResOk
type ResOK struct {
	// in:body
	Body struct {
		// enum:success
		Status string `json:"status"`
	}
}

// Fail due to user invalid input
// swagger:response GenericResFailBadRequest
type ResFailBadRequest struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to user invalid input
// swagger:response ResForbiddenRequest
type ResForbiddenRequest struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Server understand request but refuse to authorize it
// swagger:response GenericResFailConflict
type ResFailConflict struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to server understand request but unable to process
// swagger:response GenericResFailUnprocessableEntity
type ResFailUnprocessableEntity struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to resource not exists
// swagger:response GenericResFailNotFound
type ResFailNotFound struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Unexpected error occurred
// swagger:response GenericResError
type ResError struct {
	// in: body
	Body struct {
		// enum: error
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	} `json:"body"`
}
