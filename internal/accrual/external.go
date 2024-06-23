package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

type ExternalAccrualConfig struct {
	Address string
	Timeout time.Duration
	Retries int
}

type ExternalAccrual struct {
	cfg    ExternalAccrualConfig
	Client *http.Client
}

func (a ExternalAccrual) CalculateAmount(order storage.Order) (*CalculationResult, error) {
	logger := utils.Log

	url := fmt.Sprintf(`%s/api/orders/%s`, a.cfg.Address, order.OrderNumber)
	logger.Infof("About to call external accrual system at '%s'", url)
	resp, err := a.Client.Get(url)

	if err != nil {
		logger.Infof("Got error from accrual system call: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Infof("HTTP status from accrual system call is not 200: %d", resp.StatusCode)
		var result error
		switch resp.StatusCode {
		case http.StatusNoContent:
			result = ErrNoOrder
		case http.StatusTooManyRequests:
			retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
			if retryAfter == 0 {
				retryAfter = 10
			}
			result = ErrRateLimit(retryAfter)
		default:
			result = ErrInternalError
		}
		return nil, result
	}

	var result CalculationResult

	var buf bytes.Buffer
	defer resp.Body.Close()

	if _, err := buf.ReadFrom(resp.Body); err != nil {
		logger.Infof("Could not get body: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		logger.Infof("Could not parse JSON from '%s': %+v", buf.String(), err)
		return nil, err
	}

	utils.Log.WithField("response", result).Infof("Response from external accrual system")

	return &result, nil
}

func NewExternalAccrual(cfg ExternalAccrualConfig) Accrual {
	return ExternalAccrual{
		cfg: cfg,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
