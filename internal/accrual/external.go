package accrual

import (
	"encoding/json"
	"fmt"
	"io"
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
	defer resp.Body.Close()

	if err != nil {
		logger.Errorf("Got error from accrual system call: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("HTTP status from accrual system call is not 200: %d", resp.StatusCode)
		var result error
		switch resp.StatusCode {
		case http.StatusNoContent:
			result = ErrNoOrder
		case http.StatusTooManyRequests:
			retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
			if err != nil {
				logger.Errorf("Could not convert Retry-After header to int: %+v", err)
				return nil, err
			}
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

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Could not get body: %v", err)
		return nil, err
	}

	if err := json.Unmarshal(buf, &result); err != nil {
		logger.Errorf("Could not parse JSON from '%s': %+v", string(buf), err)
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
