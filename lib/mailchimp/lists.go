package mailchimp

import (
	"fmt"
	"github.com/FabLabBerlin/localmachines/models/settings"
)

func LocationSubscribe(locId int64, email string) (err error) {
	locSettings, err := settings.GetAllAt(locId)
	if err != nil {
		return fmt.Errorf("get settings: %v", err)
	}

	apiKey := locSettings.GetString(locId, settings.MAILCHIMP_API_KEY)
	if apiKey == nil {
		return fmt.Errorf("no api key set")
	}

	listId := locSettings.GetString(locId, settings.MAILCHIMP_LIST_ID)
	if listId == nil {
		return fmt.Errorf("no list id set")
	}

	s := Subscription{
		ApiKey:      *apiKey,
		Id:          *listId,
		SendWelcome: true,
	}

	s.Email.Email = email

	if err := s.Submit(); err != nil {
		return fmt.Errorf("submit: %v", err)
	}

	return
}
