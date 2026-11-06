package tgbotapi

// Виджеты (HTMLBubbles) — Bot API 10.4, живут на beta-эндпоинте
// (https://api.telegram.org/beta). Контракт методов и семантика состояния:
// https://github.com/TelegramMessenger/widget-sandbox/blob/main/docs/tgw-format.md
//
// Код виджета (HTML) задаётся при отправке и не меняется; меняется только
// состояние: chat_data — на это сообщение, public_data — на все его копии.
// Состояние пишет только бот, виджет его читает и получает уведомление
// chatDataChanged.

// Widget описывает отправленный виджет в составе сообщения.
type Widget struct {
	// ID есть только у shared-виджета (отправленного с is_shared или
	// allow_queries форвардов): по нему правят общее состояние.
	ID string `json:"id,omitempty"`

	IsCopy                   bool `json:"is_copy,omitempty"`
	CanSendQueries           bool `json:"can_send_queries,omitempty"`
	CanSendBackgroundQueries bool `json:"can_send_background_queries,omitempty"`
}

// SendWidgetConfig отправляет новый виджет (метод sendWidget).
type SendWidgetConfig struct {
	ChatID          int64
	ChannelUsername string
	MessageThreadID int

	// HTML — код виджета: до 32 КБ и до 4096 байт в zlib -9.
	HTML string

	// Width и Height — желаемый размер, 0 = по умолчанию (600×600 максимум).
	Width  int
	Height int

	// AllowQueries разрешает виджету задавать вопросы боту (иначе запросы
	// отклоняются как WIDGET_QUERY_NOT_ALLOWED).
	AllowQueries bool
	// IsShared делает виджет объектом сервера: public_data общий для всех копий.
	IsShared bool

	Caption string

	// ChatData — состояние этого сообщения (до 1024 байт). nil = не передавать.
	ChatData *string
	// PublicData — общее состояние; учитывается при IsShared.
	PublicData *string

	ReplyParameters *ReplyParameters
}

// NewSendWidget создаёт конфигурацию отправки виджета.
func NewSendWidget(chatID int64, html string) SendWidgetConfig {
	return SendWidgetConfig{ChatID: chatID, HTML: html}
}

func (config SendWidgetConfig) method() string {
	return "sendWidget"
}

func (config SendWidgetConfig) params() (Params, error) {
	params := make(Params)

	if err := params.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername); err != nil {
		return params, err
	}
	if err := params.AddInterface("html", config.HTML); err != nil {
		return params, err
	}
	params.AddNonZero("message_thread_id", config.MessageThreadID)
	params.AddNonZero("width", config.Width)
	params.AddNonZero("height", config.Height)
	params.AddBool("allow_queries", config.AllowQueries)
	params.AddBool("is_shared", config.IsShared)
	params.AddNonEmpty("caption", config.Caption)
	if config.ChatData != nil {
		params["chat_data"] = *config.ChatData
	}
	if config.PublicData != nil {
		params["public_data"] = *config.PublicData
	}
	if config.ReplyParameters != nil {
		if err := params.AddInterface("reply_parameters", config.ReplyParameters); err != nil {
			return params, err
		}
	}

	return params, nil
}

// EditMessageWidgetDataConfig правит состояние виджета в сообщении
// (метод editMessageWidgetData). Поля-указатели: nil = «не менять»,
// пустая строка = «очистить» — как в Bot API.
type EditMessageWidgetDataConfig struct {
	ChatID          int64
	ChannelUsername string
	MessageID       int

	// InlineMessageID — для сообщения через бота (инлайн-режим).
	InlineMessageID string

	ChatData   *string
	PublicData *string
}

// NewEditMessageWidgetData меняет состояние чата у виджета в сообщении.
func NewEditMessageWidgetData(chatID int64, messageID int, chatData string) EditMessageWidgetDataConfig {
	return EditMessageWidgetDataConfig{
		ChatID:    chatID,
		MessageID: messageID,
		ChatData:  &chatData,
	}
}

func (config EditMessageWidgetDataConfig) method() string {
	return "editMessageWidgetData"
}

func (config EditMessageWidgetDataConfig) params() (Params, error) {
	params := make(Params)

	if config.InlineMessageID != "" {
		params["inline_message_id"] = config.InlineMessageID
	} else {
		if err := params.AddFirstValid("chat_id", config.ChatID, config.ChannelUsername); err != nil {
			return params, err
		}
		params.AddNonZero("message_id", config.MessageID)
	}
	if config.ChatData != nil {
		params["chat_data"] = *config.ChatData
	}
	if config.PublicData != nil {
		params["public_data"] = *config.PublicData
	}

	return params, nil
}

// EditSharedWidgetDataConfig правит общее состояние shared-виджета по его id
// (метод editSharedWidgetData) — меняет его у всех сообщений сразу.
type EditSharedWidgetDataConfig struct {
	WidgetID   string
	PublicData *string
	// CacheTime — период опроса общего состояния, секунды (4..3600).
	CacheTime int
}

func (config EditSharedWidgetDataConfig) method() string {
	return "editSharedWidgetData"
}

func (config EditSharedWidgetDataConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("widget_id", config.WidgetID)
	if config.PublicData != nil {
		params["public_data"] = *config.PublicData
	}
	params.AddNonZero("cache_time", config.CacheTime)

	return params, nil
}

// AnswerWidgetQueryConfig отвечает на запрос виджета (метод answerWidgetQuery).
// Неотвеченный запрос заканчивается в виджете как TIMEOUT.
type AnswerWidgetQueryConfig struct {
	WidgetQueryID string
	// Data — результат для виджета, до 1024 байт.
	Data string
}

func (config AnswerWidgetQueryConfig) method() string {
	return "answerWidgetQuery"
}

func (config AnswerWidgetQueryConfig) params() (Params, error) {
	params := make(Params)

	params.AddNonEmpty("widget_query_id", config.WidgetQueryID)
	params.AddNonEmpty("data", config.Data)

	return params, nil
}
