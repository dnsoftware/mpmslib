package configloader

import (
	"fmt"
	"strings"
)

const partsNumber = 6 // количество частей, на которые должен быть разбит идентификатор приложения

// AppID работа с идентификатором приложения вида "ALPH_POOLER_PROP_EU_BEGET_1"
type AppID struct {
	coin        string // Код монеты (ALPH, KASPA, NEXA и т.д.)
	appType     string // Тип приложения (POOLER - собственно пул, PAYER - платежи, REWARDER - начисление вознаграждений)
	rewardType  string // Метод начисления вознаграждения (PROP, PPLNS, SOLO и т.д.)
	region      string // Код региона размещения сервера (RU, EU, ASIA, USA и т.п.)
	provider    string // Код хостинг провайдера (OVH, BEGET, RUVDS и т.п.)
	instanceTag string // Метка экземпляра (номер или какой-то другой идентификатор чтобы отличить один экземпляр приложения от другого).
	separator   string // Символ разделитель ("_", "-")
	fullID      string // Полная строка-идентификатор
}

func NewAppID(code string, separator string) (AppID, error) {
	parts := strings.Split(code, separator)
	if len(parts) != partsNumber {
		return AppID{}, fmt.Errorf("AppID parsing error: bad parts count")
	}

	return AppID{
		coin:        parts[0],
		appType:     parts[1],
		rewardType:  parts[2],
		region:      parts[3],
		provider:    parts[4],
		instanceTag: parts[5],
		separator:   separator,
		fullID:      code,
	}, nil
}

// GetFullID получить полный идентификатор
func (a *AppID) GetFullID() string {
	return a.fullID
}

// GetCoinID получить идентификатор уровня монеты (например ALPH).
// Нужно, в частности, для получение общего конфига для сервисов работающих с монетой (пул, распределение наград, выплаты)
func (a *AppID) GetCoinID() string {
	return a.coin
}
